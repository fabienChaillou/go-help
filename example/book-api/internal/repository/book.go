package repository

import (
	"context"

	"github.com/fabienChaillou/book-api/internal/model"
	"gorm.io/gorm"
)

// BookFilter carries all optional query parameters for listing books.
type BookFilter struct {
	Title     string
	Author    string
	Publisher string
	Year      int
	Status    string
	Keyword   string // matches any element in the keywords array
	Page      int
	PageSize  int
}

func (f BookFilter) offset() int {
	if f.Page <= 1 {
		return 0
	}
	return (f.Page - 1) * f.pageSize()
}

func (f BookFilter) pageSize() int {
	if f.PageSize <= 0 {
		return 10
	}
	return f.PageSize
}

// BookRepository extends the generic repository with book-specific queries.
type BookRepository interface {
	Repository[model.Book]
	List(ctx context.Context, filter BookFilter) ([]*model.Book, int64, error)
}

type bookRepository struct {
	*GormRepository[model.Book]
}

// NewBookRepository creates a BookRepository backed by GORM.
func NewBookRepository(db *gorm.DB) BookRepository {
	return &bookRepository{
		GormRepository: NewGormRepository[model.Book](db),
	}
}

// List returns a paginated, filtered list of books and the total count.
func (r *bookRepository) List(ctx context.Context, f BookFilter) ([]*model.Book, int64, error) {
	q := r.DB.WithContext(ctx).Model(&model.Book{})

	if f.Title != "" {
		q = q.Where("title ILIKE ?", "%"+f.Title+"%")
	}
	if f.Author != "" {
		q = q.Where("author ILIKE ?", "%"+f.Author+"%")
	}
	if f.Publisher != "" {
		q = q.Where("publisher ILIKE ?", "%"+f.Publisher+"%")
	}
	if f.Year > 0 {
		q = q.Where("year = ?", f.Year)
	}
	if f.Status != "" {
		q = q.Where("status = ?", f.Status)
	}
	if f.Keyword != "" {
		// PostgreSQL: check if the keyword exists in the text[] column
		q = q.Where("? = ANY(keywords)", f.Keyword)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var books []*model.Book
	if err := q.Offset(f.offset()).Limit(f.pageSize()).Find(&books).Error; err != nil {
		return nil, 0, err
	}

	return books, total, nil
}
