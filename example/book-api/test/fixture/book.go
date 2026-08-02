package fixture

import (
	"github.com/fabienChaillou/book-api/internal/model"
	"github.com/lib/pq"
)

// Books returns a set of deterministic book fixtures for testing.
func Books() []*model.Book {
	return []*model.Book{
		{
			Title:       "The Go Programming Language",
			Author:      "Alan Donovan",
			Description: "Comprehensive guide to Go",
			Year:        2015,
			Status:      model.StatusAvailable,
			Keywords:    pq.StringArray{"go", "programming", "systems"},
			ISBN:        "978-0134190440",
			Publisher:   "Addison-Wesley",
			Pages:       380,
		},
		{
			Title:       "Clean Code",
			Author:      "Robert C. Martin",
			Description: "A handbook of agile software craftsmanship",
			Year:        2008,
			Status:      model.StatusAvailable,
			Keywords:    pq.StringArray{"clean-code", "refactoring", "solid"},
			ISBN:        "978-0132350884",
			Publisher:   "Prentice Hall",
			Pages:       431,
		},
		{
			Title:       "Domain-Driven Design",
			Author:      "Eric Evans",
			Description: "Tackling complexity in the heart of software",
			Year:        2003,
			Status:      model.StatusBorrowed,
			Keywords:    pq.StringArray{"ddd", "architecture", "domain"},
			ISBN:        "978-0321125217",
			Publisher:   "Addison-Wesley",
			Pages:       560,
		},
		{
			Title:       "Designing Data-Intensive Applications",
			Author:      "Martin Kleppmann",
			Description: "The big ideas behind reliable, scalable systems",
			Year:        2017,
			Status:      model.StatusAvailable,
			Keywords:    pq.StringArray{"distributed", "databases", "streaming"},
			ISBN:        "978-1449373320",
			Publisher:   "O'Reilly",
			Pages:       611,
		},
		{
			Title:       "The Pragmatic Programmer",
			Author:      "David Thomas",
			Description: "Your journey to mastery",
			Year:        2019,
			Status:      model.StatusReserved,
			Keywords:    pq.StringArray{"pragmatic", "career", "craftsmanship"},
			ISBN:        "978-0135957059",
			Publisher:   "Addison-Wesley",
			Pages:       352,
		},
	}
}

// SingleBook returns one fixture book (useful for single-entity tests).
func SingleBook() *model.Book {
	return Books()[0]
}
