package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/fabienChaillou/book-api/internal/api/dto"
	"github.com/fabienChaillou/book-api/internal/model"
	"github.com/fabienChaillou/book-api/internal/repository"
	"github.com/fabienChaillou/book-api/internal/service"
)

// BookHandler holds the Huma handlers for the /books resource.
type BookHandler struct {
	svc service.BookService
}

func NewBookHandler(svc service.BookService) *BookHandler {
	return &BookHandler{svc: svc}
}

// Register wires all book routes onto the given Huma API.
func (h *BookHandler) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "create-book",
		Method:      http.MethodPost,
		Path:        "/books",
		Summary:     "Create a book",
		Tags:        []string{"books"},
	}, h.CreateBook)

	huma.Register(api, huma.Operation{
		OperationID: "create-books-batch",
		Method:      http.MethodPost,
		Path:        "/books/batch",
		Summary:     "Create multiple books",
		Tags:        []string{"books"},
	}, h.CreateBooksBatch)

	huma.Register(api, huma.Operation{
		OperationID: "list-books",
		Method:      http.MethodGet,
		Path:        "/books",
		Summary:     "List books with optional filters",
		Tags:        []string{"books"},
	}, h.ListBooks)

	huma.Register(api, huma.Operation{
		OperationID: "get-book",
		Method:      http.MethodGet,
		Path:        "/books/{id}",
		Summary:     "Get a book by ID",
		Tags:        []string{"books"},
	}, h.GetBook)

	huma.Register(api, huma.Operation{
		OperationID: "update-book",
		Method:      http.MethodPut,
		Path:        "/books/{id}",
		Summary:     "Update a book",
		Tags:        []string{"books"},
	}, h.UpdateBook)

	huma.Register(api, huma.Operation{
		OperationID: "delete-book",
		Method:      http.MethodDelete,
		Path:        "/books/{id}",
		Summary:     "Delete a book",
		Tags:        []string{"books"},
	}, h.DeleteBook)

	huma.Register(api, huma.Operation{
		OperationID: "delete-books-batch",
		Method:      http.MethodDelete,
		Path:        "/books/batch",
		Summary:     "Delete multiple books",
		Tags:        []string{"books"},
	}, h.DeleteBooksBatch)
}

// ---- Input / Output types ----

type CreateBookInput struct {
	Body dto.CreateBookRequest
}
type CreateBookOutput struct {
	Body *dto.BookResponse
}

type CreateBatchInput struct {
	Body dto.CreateBatchRequest
}
type CreateBatchOutput struct {
	Body []*dto.BookResponse
}

type ListBooksInput struct {
	Title     string `query:"title"      doc:"Filter by title (partial, case-insensitive)"`
	Author    string `query:"author"     doc:"Filter by author (partial, case-insensitive)"`
	Publisher string `query:"publisher"  doc:"Filter by publisher (partial, case-insensitive)"`
	Year      int    `query:"year"       doc:"Filter by publication year"`
	Status    string `query:"status"     doc:"Filter by status (available|borrowed|reserved|lost)"`
	Keyword   string `query:"keyword"    doc:"Filter by exact keyword match"`
	Page      int    `query:"page"       doc:"Page number (default 1)"     default:"1"`
	PageSize  int    `query:"page_size"  doc:"Items per page (default 10)" default:"10"`
}

type ListBooksOutput struct {
	Body dto.ListBooksResponse
}

type GetBookInput struct {
	ID uint `path:"id"`
}
type GetBookOutput struct {
	Body *dto.BookResponse
}

type UpdateBookInput struct {
	ID   uint `path:"id"`
	Body dto.UpdateBookRequest
}
type UpdateBookOutput struct {
	Body *dto.BookResponse
}

type DeleteBookInput struct {
	ID uint `path:"id"`
}

type DeleteBatchInput struct {
	Body dto.DeleteBatchRequest
}

// ---- Handlers ----

func (h *BookHandler) CreateBook(ctx context.Context, input *CreateBookInput) (*CreateBookOutput, error) {
	book := dto.FromCreateRequest(input.Body)
	if err := h.svc.Create(ctx, book); err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}
	return &CreateBookOutput{Body: dto.ToBookResponse(book)}, nil
}

func (h *BookHandler) CreateBooksBatch(ctx context.Context, input *CreateBatchInput) (*CreateBatchOutput, error) {
	books := make([]*model.Book, 0, len(input.Body.Books))
	for _, req := range input.Body.Books {
		books = append(books, dto.FromCreateRequest(req))
	}

	if err := h.svc.CreateBatch(ctx, books); err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}

	resp := make([]*dto.BookResponse, 0, len(books))
	for _, b := range books {
		resp = append(resp, dto.ToBookResponse(b))
	}
	return &CreateBatchOutput{Body: resp}, nil
}

func (h *BookHandler) ListBooks(ctx context.Context, input *ListBooksInput) (*ListBooksOutput, error) {
	filter := repository.BookFilter{
		Title:     input.Title,
		Author:    input.Author,
		Publisher: input.Publisher,
		Year:      input.Year,
		Status:    input.Status,
		Keyword:   input.Keyword,
		Page:      input.Page,
		PageSize:  input.PageSize,
	}

	books, total, err := h.svc.List(ctx, filter)
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	items := make([]*dto.BookResponse, 0, len(books))
	for _, b := range books {
		items = append(items, dto.ToBookResponse(b))
	}

	return &ListBooksOutput{
		Body: dto.ListBooksResponse{
			Items:    items,
			Total:    total,
			Page:     input.Page,
			PageSize: input.PageSize,
		},
	}, nil
}

func (h *BookHandler) GetBook(ctx context.Context, input *GetBookInput) (*GetBookOutput, error) {
	book, err := h.svc.GetByID(ctx, input.ID)
	if err != nil {
		return nil, mapServiceError(err)
	}
	return &GetBookOutput{Body: dto.ToBookResponse(book)}, nil
}

func (h *BookHandler) UpdateBook(ctx context.Context, input *UpdateBookInput) (*UpdateBookOutput, error) {
	book, err := h.svc.GetByID(ctx, input.ID)
	if err != nil {
		return nil, mapServiceError(err)
	}

	dto.ApplyUpdateRequest(book, input.Body)

	if err := h.svc.Update(ctx, book); err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}
	return &UpdateBookOutput{Body: dto.ToBookResponse(book)}, nil
}

func (h *BookHandler) DeleteBook(ctx context.Context, input *DeleteBookInput) (*struct{}, error) {
	if err := h.svc.Delete(ctx, input.ID); err != nil {
		return nil, mapServiceError(err)
	}
	return nil, nil
}

func (h *BookHandler) DeleteBooksBatch(ctx context.Context, input *DeleteBatchInput) (*struct{}, error) {
	if err := h.svc.DeleteBatch(ctx, input.Body.IDs); err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}
	return nil, nil
}

func mapServiceError(err error) error {
	if errors.Is(err, service.ErrNotFound) {
		return huma.Error404NotFound("book not found")
	}
	return huma.Error500InternalServerError(err.Error())
}
