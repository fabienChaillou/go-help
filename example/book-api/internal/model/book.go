package model

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// BookStatus represents the availability state of a book.
type BookStatus string

const (
	StatusAvailable BookStatus = "available"
	StatusBorrowed  BookStatus = "borrowed"
	StatusReserved  BookStatus = "reserved"
	StatusLost      BookStatus = "lost"
)

// Book is the GORM model persisted in the `books` table.
// gorm.Model provides ID, CreatedAt, UpdatedAt and DeletedAt (soft-delete).
type Book struct {
	gorm.Model
	Title       string         `gorm:"not null;index"                json:"title"`
	Author      string         `gorm:"not null;index"                json:"author"`
	Description string         `                                     json:"description"`
	Year        int            `gorm:"index"                         json:"year"`
	Status      BookStatus     `gorm:"type:varchar(20);default:'available';index" json:"status"`
	Keywords    pq.StringArray `gorm:"type:text[]"                   json:"keywords"`
	ISBN        string         `gorm:"uniqueIndex"                   json:"isbn"`
	Publisher   string         `gorm:"index"                         json:"publisher"`
	Pages       int            `                                     json:"pages"`
}
