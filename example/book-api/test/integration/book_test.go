package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/fabienChaillou/book-api/internal/api/dto"
	apihandler "github.com/fabienChaillou/book-api/internal/api/handler"
	"github.com/fabienChaillou/book-api/internal/model"
	"github.com/fabienChaillou/book-api/internal/repository"
	"github.com/fabienChaillou/book-api/internal/service"
	"github.com/fabienChaillou/book-api/test/fixture"
)

// BookSuite groups all book integration tests against a real PostgreSQL container.
type BookSuite struct {
	suite.Suite
	container *tcpostgres.PostgresContainer
	db        *gorm.DB
	server    *httptest.Server
}

func TestBookSuite(t *testing.T) {
	suite.Run(t, new(BookSuite))
}

// SetupSuite starts the Postgres container and migrates the schema once per suite.
func (s *BookSuite) SetupSuite() {
	ctx := context.Background()

	container, err := tcpostgres.Run(ctx,
		"postgres:16-alpine",
		tcpostgres.WithDatabase("testdb"),
		tcpostgres.WithUsername("testuser"),
		tcpostgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2),
		),
	)
	require.NoError(s.T(), err)
	s.container = container

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(s.T(), err)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.NoError(s.T(), err)
	require.NoError(s.T(), db.AutoMigrate(&model.Book{}))
	s.db = db

	s.server = s.newTestServer()
}

func (s *BookSuite) TearDownSuite() {
	s.server.Close()
	if s.container != nil {
		_ = s.container.Terminate(context.Background())
	}
}

// TearDownTest wipes the books table between each test for isolation.
func (s *BookSuite) TearDownTest() {
	s.db.Exec("DELETE FROM books")
}

func (s *BookSuite) newTestServer() *httptest.Server {
	repo := repository.NewBookRepository(s.db)
	svc := service.NewBookService(repo)
	h := apihandler.NewBookHandler(svc)

	r := chi.NewRouter()
	api := humachi.New(r, huma.DefaultConfig("Book API Test", "1.0.0"))
	h.Register(api)

	return httptest.NewServer(r)
}

// ---- Helpers ----

func (s *BookSuite) url(path string) string {
	return s.server.URL + path
}

func (s *BookSuite) mustSeedBooks() []*model.Book {
	books := fixture.Books()
	for _, b := range books {
		require.NoError(s.T(), s.db.Create(b).Error)
	}
	return books
}

func (s *BookSuite) doJSON(method, url string, body interface{}) *http.Response {
	var buf bytes.Buffer
	if body != nil {
		require.NoError(s.T(), json.NewEncoder(&buf).Encode(body))
	}
	req, err := http.NewRequest(method, url, &buf)
	require.NoError(s.T(), err)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	require.NoError(s.T(), err)
	return resp
}

func (s *BookSuite) decodeBody(resp *http.Response, v interface{}) {
	s.T().Helper()
	defer resp.Body.Close()
	require.NoError(s.T(), json.NewDecoder(resp.Body).Decode(v))
}

// ---- Tests ----

func (s *BookSuite) TestCreateBook_Success() {
	payload := dto.CreateBookRequest{
		Title:     "Test Book",
		Author:    "Test Author",
		Year:      2024,
		Status:    "available",
		Keywords:  []string{"test", "go"},
		ISBN:      "978-0000000001",
		Publisher: "Test Publisher",
		Pages:     100,
	}

	resp := s.doJSON(http.MethodPost, s.url("/books"), payload)
	assert.Equal(s.T(), http.StatusOK, resp.StatusCode)

	var result dto.BookResponse
	s.decodeBody(resp, &result)

	assert.Equal(s.T(), "Test Book", result.Title)
	assert.Equal(s.T(), "Test Author", result.Author)
	assert.NotZero(s.T(), result.ID)
	assert.Equal(s.T(), []string{"test", "go"}, result.Keywords)
}

func (s *BookSuite) TestCreateBook_MissingTitle() {
	payload := dto.CreateBookRequest{Author: "Someone"}
	resp := s.doJSON(http.MethodPost, s.url("/books"), payload)
	assert.Equal(s.T(), http.StatusUnprocessableEntity, resp.StatusCode)
	resp.Body.Close()
}

func (s *BookSuite) TestCreateBatch_Success() {
	payload := dto.CreateBatchRequest{
		Books: []dto.CreateBookRequest{
			{Title: "Batch A", Author: "Author A", ISBN: "978-0000000010"},
			{Title: "Batch B", Author: "Author B", ISBN: "978-0000000011"},
		},
	}

	resp := s.doJSON(http.MethodPost, s.url("/books/batch"), payload)
	assert.Equal(s.T(), http.StatusOK, resp.StatusCode)

	var results []dto.BookResponse
	s.decodeBody(resp, &results)
	assert.Len(s.T(), results, 2)
}

func (s *BookSuite) TestGetBook_Success() {
	books := s.mustSeedBooks()
	target := books[0]

	resp, err := http.Get(s.url(fmt.Sprintf("/books/%d", target.ID)))
	require.NoError(s.T(), err)
	assert.Equal(s.T(), http.StatusOK, resp.StatusCode)

	var result dto.BookResponse
	s.decodeBody(resp, &result)
	assert.Equal(s.T(), target.Title, result.Title)
}

func (s *BookSuite) TestGetBook_NotFound() {
	resp, err := http.Get(s.url("/books/9999999"))
	require.NoError(s.T(), err)
	defer resp.Body.Close()
	assert.Equal(s.T(), http.StatusNotFound, resp.StatusCode)
}

func (s *BookSuite) TestListBooks_NoFilter() {
	s.mustSeedBooks()

	resp, err := http.Get(s.url("/books"))
	require.NoError(s.T(), err)
	assert.Equal(s.T(), http.StatusOK, resp.StatusCode)

	var result dto.ListBooksResponse
	s.decodeBody(resp, &result)

	assert.Equal(s.T(), int64(5), result.Total)
	assert.Len(s.T(), result.Items, 5)
}

func (s *BookSuite) TestListBooks_FilterByAuthor() {
	s.mustSeedBooks()

	resp, err := http.Get(s.url("/books?author=Donovan"))
	require.NoError(s.T(), err)

	var result dto.ListBooksResponse
	s.decodeBody(resp, &result)

	assert.Equal(s.T(), int64(1), result.Total)
	assert.Equal(s.T(), "The Go Programming Language", result.Items[0].Title)
}

func (s *BookSuite) TestListBooks_FilterByStatus() {
	s.mustSeedBooks()

	resp, err := http.Get(s.url("/books?status=borrowed"))
	require.NoError(s.T(), err)

	var result dto.ListBooksResponse
	s.decodeBody(resp, &result)

	assert.Equal(s.T(), int64(1), result.Total)
	assert.Equal(s.T(), "Domain-Driven Design", result.Items[0].Title)
}

func (s *BookSuite) TestListBooks_FilterByKeyword() {
	s.mustSeedBooks()

	resp, err := http.Get(s.url("/books?keyword=ddd"))
	require.NoError(s.T(), err)

	var result dto.ListBooksResponse
	s.decodeBody(resp, &result)

	assert.Equal(s.T(), int64(1), result.Total)
	assert.Equal(s.T(), "Domain-Driven Design", result.Items[0].Title)
}

func (s *BookSuite) TestListBooks_FilterByPublisher() {
	s.mustSeedBooks()

	resp, err := http.Get(s.url("/books?publisher=Addison"))
	require.NoError(s.T(), err)

	var result dto.ListBooksResponse
	s.decodeBody(resp, &result)

	// 3 books from Addison-Wesley
	assert.Equal(s.T(), int64(3), result.Total)
}

func (s *BookSuite) TestListBooks_FilterByYear() {
	s.mustSeedBooks()

	resp, err := http.Get(s.url("/books?year=2015"))
	require.NoError(s.T(), err)

	var result dto.ListBooksResponse
	s.decodeBody(resp, &result)

	assert.Equal(s.T(), int64(1), result.Total)
}

func (s *BookSuite) TestListBooks_Pagination() {
	s.mustSeedBooks()

	resp, err := http.Get(s.url("/books?page=1&page_size=2"))
	require.NoError(s.T(), err)

	var result dto.ListBooksResponse
	s.decodeBody(resp, &result)

	assert.Equal(s.T(), int64(5), result.Total)
	assert.Len(s.T(), result.Items, 2)
}

func (s *BookSuite) TestUpdateBook_Success() {
	books := s.mustSeedBooks()
	target := books[0]

	update := dto.UpdateBookRequest{Title: "Updated Title", Status: "borrowed"}
	resp := s.doJSON(http.MethodPut, s.url(fmt.Sprintf("/books/%d", target.ID)), update)
	assert.Equal(s.T(), http.StatusOK, resp.StatusCode)

	var result dto.BookResponse
	s.decodeBody(resp, &result)

	assert.Equal(s.T(), "Updated Title", result.Title)
	assert.Equal(s.T(), "borrowed", result.Status)
	// Author unchanged
	assert.Equal(s.T(), target.Author, result.Author)
}

func (s *BookSuite) TestUpdateBook_NotFound() {
	update := dto.UpdateBookRequest{Title: "X"}
	resp := s.doJSON(http.MethodPut, s.url("/books/9999999"), update)
	defer resp.Body.Close()
	assert.Equal(s.T(), http.StatusNotFound, resp.StatusCode)
}

func (s *BookSuite) TestDeleteBook_Success() {
	books := s.mustSeedBooks()
	target := books[0]

	req, _ := http.NewRequest(http.MethodDelete, s.url(fmt.Sprintf("/books/%d", target.ID)), nil)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(s.T(), err)
	resp.Body.Close()
	// Huma returns 204 No Content for empty responses
	assert.Equal(s.T(), http.StatusNoContent, resp.StatusCode)

	// Verify soft-delete: record not found via API
	getResp, _ := http.Get(s.url(fmt.Sprintf("/books/%d", target.ID)))
	defer getResp.Body.Close()
	assert.Equal(s.T(), http.StatusNotFound, getResp.StatusCode)
}

func (s *BookSuite) TestDeleteBatch_Success() {
	books := s.mustSeedBooks()
	ids := []uint{books[0].ID, books[1].ID}

	payload := dto.DeleteBatchRequest{IDs: ids}
	var buf bytes.Buffer
	require.NoError(s.T(), json.NewEncoder(&buf).Encode(payload))

	req, err := http.NewRequest(http.MethodDelete, s.url("/books/batch"), &buf)
	require.NoError(s.T(), err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(s.T(), err)
	resp.Body.Close()
	assert.Equal(s.T(), http.StatusNoContent, resp.StatusCode)

	// 3 books should remain
	listResp, _ := http.Get(s.url("/books"))
	var result dto.ListBooksResponse
	s.decodeBody(listResp, &result)
	assert.Equal(s.T(), int64(3), result.Total)
}

// ---- Unit-level service tests (no HTTP, no container needed) ----

func TestBookService_Validate(t *testing.T) {
	// These tests use an in-memory stub instead of the real DB.
	t.Run("missing title returns error", func(t *testing.T) {
		svc := service.NewBookService(&stubRepo{})
		err := svc.Create(context.Background(), &model.Book{Author: "A"})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "title")
	})

	t.Run("missing author returns error", func(t *testing.T) {
		svc := service.NewBookService(&stubRepo{})
		err := svc.Create(context.Background(), &model.Book{Title: "T"})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "author")
	})

	t.Run("valid book is accepted", func(t *testing.T) {
		svc := service.NewBookService(&stubRepo{})
		err := svc.Create(context.Background(), &model.Book{Title: "T", Author: "A"})
		require.NoError(t, err)
	})
}

// stubRepo is a minimal in-memory BookRepository for unit tests.
type stubRepo struct {
	books []*model.Book
}

func (r *stubRepo) Create(_ context.Context, b *model.Book) error {
	b.ID = uint(len(r.books) + 1)
	r.books = append(r.books, b)
	return nil
}
func (r *stubRepo) CreateBatch(_ context.Context, bs []*model.Book) error {
	for _, b := range bs {
		_ = r.Create(context.Background(), b)
	}
	return nil
}
func (r *stubRepo) FindByID(_ context.Context, id uint) (*model.Book, error) {
	for _, b := range r.books {
		if b.ID == id {
			return b, nil
		}
	}
	return nil, fmt.Errorf("record not found")
}
func (r *stubRepo) Update(_ context.Context, b *model.Book) error {
	for i, existing := range r.books {
		if existing.ID == b.ID {
			r.books[i] = b
			return nil
		}
	}
	return fmt.Errorf("record not found")
}
func (r *stubRepo) Delete(_ context.Context, id uint) error {
	for i, b := range r.books {
		if b.ID == id {
			r.books = append(r.books[:i], r.books[i+1:]...)
			return nil
		}
	}
	return nil
}
func (r *stubRepo) DeleteBatch(_ context.Context, ids []uint) error {
	idSet := make(map[uint]bool, len(ids))
	for _, id := range ids {
		idSet[id] = true
	}
	filtered := r.books[:0]
	for _, b := range r.books {
		if !idSet[b.ID] {
			filtered = append(filtered, b)
		}
	}
	r.books = filtered
	return nil
}
func (r *stubRepo) List(_ context.Context, _ repository.BookFilter) ([]*model.Book, int64, error) {
	return r.books, int64(len(r.books)), nil
}

// ensure unused import (io) is used somewhere
