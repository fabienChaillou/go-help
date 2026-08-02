package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/fabienChaillou/book-api/internal/api/handler"
	"github.com/fabienChaillou/book-api/internal/config"
	"github.com/fabienChaillou/book-api/internal/database"
	"github.com/fabienChaillou/book-api/internal/repository"
	"github.com/fabienChaillou/book-api/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg := config.Load()

	db, err := database.New(cfg.DB)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	// Dependency injection (DIP: wire concrete types through interfaces)
	bookRepo := repository.NewBookRepository(db)
	bookSvc := service.NewBookService(bookRepo)
	bookHandler := handler.NewBookHandler(bookSvc)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	apiCfg := huma.DefaultConfig("Book API", "1.0.0")
	apiCfg.Info.Description = "CRUD API for managing books"
	api := humachi.New(r, apiCfg)

	bookHandler.Register(api)

	// Graceful context (extend as needed with signal handling)
	_ = context.Background()

	addr := fmt.Sprintf(":%s", cfg.HTTP.Port)
	log.Printf("starting server on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
