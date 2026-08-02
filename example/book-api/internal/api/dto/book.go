package dto

import (
	"time"

	"github.com/fabienChaillou/book-api/internal/model"
	"github.com/lib/pq"
)

// ---- Request DTOs ----

// CreateBookRequest is the payload for POST /books.
type CreateBookRequest struct {
	Title       string   `json:"title"       required:"true"`
	Author      string   `json:"author"      required:"true"`
	Description string   `json:"description"`
	Year        int      `json:"year"`
	Status      string   `json:"status"`
	Keywords    []string `json:"keywords"`
	ISBN        string   `json:"isbn"`
	Publisher   string   `json:"publisher"`
	Pages       int      `json:"pages"`
}

// CreateBatchRequest is the payload for POST /books/batch.
type CreateBatchRequest struct {
	Books []CreateBookRequest `json:"books" required:"true"`
}

// UpdateBookRequest is the payload for PUT /books/{id}.
type UpdateBookRequest struct {
	Title       string   `json:"title"`
	Author      string   `json:"author"`
	Description string   `json:"description"`
	Year        int      `json:"year"`
	Status      string   `json:"status"`
	Keywords    []string `json:"keywords"`
	ISBN        string   `json:"isbn"`
	Publisher   string   `json:"publisher"`
	Pages       int      `json:"pages"`
}

// DeleteBatchRequest is the payload for DELETE /books/batch.
type DeleteBatchRequest struct {
	IDs []uint `json:"ids" required:"true"`
}

// ---- Response DTOs ----

// BookResponse is returned for a single book.
type BookResponse struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	Description string    `json:"description"`
	Year        int       `json:"year"`
	Status      string    `json:"status"`
	Keywords    []string  `json:"keywords"`
	ISBN        string    `json:"isbn"`
	Publisher   string    `json:"publisher"`
	Pages       int       `json:"pages"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ListBooksResponse wraps a paginated list.
type ListBooksResponse struct {
	Items    []*BookResponse `json:"items"`
	Total    int64           `json:"total"`
	Page     int             `json:"page"`
	PageSize int             `json:"page_size"`
}

// ---- Mappers ----

func ToBookResponse(b *model.Book) *BookResponse {
	kw := []string(b.Keywords)
	if kw == nil {
		kw = []string{}
	}
	return &BookResponse{
		ID:          b.ID,
		Title:       b.Title,
		Author:      b.Author,
		Description: b.Description,
		Year:        b.Year,
		Status:      string(b.Status),
		Keywords:    kw,
		ISBN:        b.ISBN,
		Publisher:   b.Publisher,
		Pages:       b.Pages,
		CreatedAt:   b.CreatedAt,
		UpdatedAt:   b.UpdatedAt,
	}
}

func FromCreateRequest(req CreateBookRequest) *model.Book {
	status := model.BookStatus(req.Status)
	if status == "" {
		status = model.StatusAvailable
	}
	return &model.Book{
		Title:       req.Title,
		Author:      req.Author,
		Description: req.Description,
		Year:        req.Year,
		Status:      status,
		Keywords:    pq.StringArray(req.Keywords),
		ISBN:        req.ISBN,
		Publisher:   req.Publisher,
		Pages:       req.Pages,
	}
}

func ApplyUpdateRequest(b *model.Book, req UpdateBookRequest) {
	if req.Title != "" {
		b.Title = req.Title
	}
	if req.Author != "" {
		b.Author = req.Author
	}
	if req.Description != "" {
		b.Description = req.Description
	}
	if req.Year != 0 {
		b.Year = req.Year
	}
	if req.Status != "" {
		b.Status = model.BookStatus(req.Status)
	}
	if req.Keywords != nil {
		b.Keywords = pq.StringArray(req.Keywords)
	}
	if req.ISBN != "" {
		b.ISBN = req.ISBN
	}
	if req.Publisher != "" {
		b.Publisher = req.Publisher
	}
	if req.Pages != 0 {
		b.Pages = req.Pages
	}
}
