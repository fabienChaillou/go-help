package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/fabienChaillou/book-api/internal/model"
	"github.com/fabienChaillou/book-api/internal/repository"
	"gorm.io/gorm"
)

// ErrNotFound is returned when a book does not exist.
var ErrNotFound = errors.New("book not found")

// BookService defines the business operations on books.
// Depending on service abstractions on repository interfaces (DIP).
type BookService interface {
	Create(ctx context.Context, book *model.Book) error
	CreateBatch(ctx context.Context, books []*model.Book) error
	GetByID(ctx context.Context, id uint) (*model.Book, error)
	List(ctx context.Context, filter repository.BookFilter) ([]*model.Book, int64, error)
	Update(ctx context.Context, book *model.Book) error
	Delete(ctx context.Context, id uint) error
	DeleteBatch(ctx context.Context, ids []uint) error
}

type bookService struct {
	repo repository.BookRepository
}

// NewBookService creates a BookService using the provided repository.
func NewBookService(repo repository.BookRepository) BookService {
	return &bookService{repo: repo}
}

func (s *bookService) Create(ctx context.Context, book *model.Book) error {
	if err := validate(book); err != nil {
		return err
	}
	return s.repo.Create(ctx, book)
}

func (s *bookService) CreateBatch(ctx context.Context, books []*model.Book) error {
	for _, b := range books {
		if err := validate(b); err != nil {
			return fmt.Errorf("book %q: %w", b.Title, err)
		}
	}
	return s.repo.CreateBatch(ctx, books)
}

func (s *bookService) GetByID(ctx context.Context, id uint) (*model.Book, error) {
	book, err := s.repo.FindByID(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return book, err
}

func (s *bookService) List(ctx context.Context, filter repository.BookFilter) ([]*model.Book, int64, error) {
	return s.repo.List(ctx, filter)
}

func (s *bookService) Update(ctx context.Context, book *model.Book) error {
	if _, err := s.GetByID(ctx, book.ID); err != nil {
		return err
	}
	if err := validate(book); err != nil {
		return err
	}
	return s.repo.Update(ctx, book)
}

func (s *bookService) Delete(ctx context.Context, id uint) error {
	if _, err := s.GetByID(ctx, id); err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}

func (s *bookService) DeleteBatch(ctx context.Context, ids []uint) error {
	return s.repo.DeleteBatch(ctx, ids)
}

// validate enforces invariants before persistence.
func validate(b *model.Book) error {
	if b.Title == "" {
		return errors.New("title is required")
	}
	if b.Author == "" {
		return errors.New("author is required")
	}
	return nil
}
