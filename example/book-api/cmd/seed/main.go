package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/fabienChaillou/book-api/internal/config"
	"github.com/fabienChaillou/book-api/internal/database"
	"github.com/fabienChaillou/book-api/internal/model"
	"github.com/lib/pq"
)

func main() {
	cfg := config.Load()

	db, err := database.New(cfg.DB)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	books := make([]*model.Book, 0, 1000)
	for i := 1; i <= 1000; i++ {
		books = append(books, randomBook(rng, i))
	}

	result := db.CreateInBatches(books, 100)
	if result.Error != nil {
		log.Fatalf("seed failed: %v", result.Error)
	}

	fmt.Printf("✓ %d books inserted\n", result.RowsAffected)
}

// ---- Data pools ----

var titles = []string{
	"The Art of %s", "Introduction to %s", "Mastering %s",
	"Advanced %s", "Learning %s", "Practical %s",
	"%s in Depth", "%s: A Complete Guide", "Modern %s",
	"Effective %s", "Understanding %s", "The %s Handbook",
	"%s Patterns", "Building with %s", "Exploring %s",
}

var topics = []string{
	"Go", "Rust", "Python", "Kubernetes", "Microservices",
	"Distributed Systems", "Machine Learning", "Clean Architecture",
	"Domain-Driven Design", "Algorithms", "Data Structures",
	"Concurrency", "API Design", "DevOps", "Cloud Computing",
	"PostgreSQL", "Redis", "Kafka", "GraphQL", "gRPC",
	"Functional Programming", "System Design", "Networking",
	"Cryptography", "Compilers", "Operating Systems",
}

var authors = []string{
	"Alice Martin", "Bob Chen", "Claire Dubois", "David Kim",
	"Emma Wilson", "François Leclerc", "Grace Okafor", "Henri Bernard",
	"Isabelle Roy", "James Nakamura", "Karen Patel", "Luca Romano",
	"Marie Curie Jr.", "Nathan Brooks", "Olivia Johnson", "Pierre Dupont",
	"Quinn Murphy", "Rachel Zhang", "Samuel Osei", "Tina Andersen",
}

var publishers = []string{
	"Addison-Wesley", "O'Reilly", "Prentice Hall", "Manning",
	"Packt", "No Starch Press", "Apress", "Springer", "MIT Press",
}

var statuses = []model.BookStatus{
	model.StatusAvailable, model.StatusAvailable, model.StatusAvailable,
	model.StatusBorrowed,
	model.StatusReserved,
	model.StatusLost,
}

var keywordPool = []string{
	"programming", "systems", "architecture", "patterns", "algorithms",
	"databases", "networking", "security", "testing", "devops",
	"cloud", "containers", "microservices", "api", "performance",
	"concurrency", "functional", "oop", "solid", "clean-code",
	"distributed", "streaming", "caching", "monitoring", "ci-cd",
}

// ---- Generator ----

func randomBook(rng *rand.Rand, i int) *model.Book {
	topic := topics[rng.Intn(len(topics))]
	titleFmt := titles[rng.Intn(len(titles))]

	kw := pickN(rng, keywordPool, 2+rng.Intn(4)) // 2–5 keywords

	return &model.Book{
		Title:       fmt.Sprintf(titleFmt, topic),
		Author:      authors[rng.Intn(len(authors))],
		Description: fmt.Sprintf("Book #%04d — a deep dive into %s.", i, topic),
		Year:        1990 + rng.Intn(36), // 1990–2025
		Status:      statuses[rng.Intn(len(statuses))],
		Keywords:    pq.StringArray(kw),
		ISBN:        fmt.Sprintf("978-%010d", i),
		Publisher:   publishers[rng.Intn(len(publishers))],
		Pages:       100 + rng.Intn(801), // 100–900
	}
}

// pickN returns n distinct elements picked at random from pool.
func pickN(rng *rand.Rand, pool []string, n int) []string {
	if n > len(pool) {
		n = len(pool)
	}
	perm := rng.Perm(len(pool))
	out := make([]string, n)
	for i := range out {
		out[i] = pool[perm[i]]
	}
	return out
}
