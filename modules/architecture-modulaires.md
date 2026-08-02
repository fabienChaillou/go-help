# Architecture modulaire en Go — proposition complète (exemple prêt à adapter)

Voici une architecture modulaire (inspirée de *hexagonal / ports & adapters* + *service/repository*), pensée pour une API web en Go (Gin + sqlx), facilement testable et extensible. Je fournis :

1. l'arborescence recommandée,
2. principes courts,
3. exemples de fichiers clés (main, router, controller/handler, service, repository, domain, config, db),
4. notes sur tests, mocks et déploiement.

---

# 1. Arborescence proposée

```
myapp/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── app/
│   │   ├── http/
│   │   │   ├── router.go
│   │   │   └── user_handler.go
│   │   ├── service/
│   │   │   └── user_service.go
│   │   └── repository/
│   │       └── user_repo_sql.go
│   ├── domain/
│   │   └── user.go
│   └── pkg/
│       ├── config/
│       │   └── config.go
│       ├── db/
│       │   └── db.go
│       └── logger/
│           └── logger.go
├── migrations/
├── scripts/
├── Dockerfile
├── Makefile
└── go.mod
```

---

# 2. Principes / responsabilités

* `cmd/server`: bootstrap de l'application (DI minimale, lecture config, ouverture DB, démarrage HTTP).
* `internal/pkg`: utilitaires réutilisables (config, DB, logger, middleware).
* `internal/domain`: entités métier (structures + règles éventuelles).
* `internal/app/service`: logique métier (use-cases). Ne connaît pas Gin ni SQL.
* `internal/app/repository`: implémentations concrètes des ports pour stockage (sqlx, pgx...).
* `internal/app/http`: adaptateurs HTTP (handlers, validation, conversion DTO <-> domain).
* Tests unitaires par couche ; mocks pour interfaces (ex : `mockgen` ou `mockery`).

---

# 3. Exemples de code

> Les exemples sont réduits pour la lisibilité — adapte selon ton style (erreurs, context, logger, validation, transactions...).

### `cmd/server/main.go`

```go
package main

import (
    "context"
    "log"
    "net/http"
    "time"

    "myapp/internal/pkg/config"
    "myapp/internal/pkg/db"
    "myapp/internal/pkg/logger"
    "myapp/internal/app/repository"
    "myapp/internal/app/service"
    "myapp/internal/app/http"
)

func main() {
    cfg := config.Load() // lit env / fichier
    logg := logger.New(cfg)

    sqlxDB, err := db.NewSQLX(cfg.DatabaseURL)
    if err != nil {
        log.Fatalf("db init: %v", err)
    }
    defer sqlxDB.Close()

    // Repos
    userRepo := repository.NewUserSQLRepo(sqlxDB)

    // Services (use-cases)
    userSvc := service.NewUserService(userRepo, logg)

    // HTTP layer
    router := http.NewRouter(userSvc, logg)

    srv := &http.Server{
        Addr:    cfg.HTTPAddr,
        Handler: router,
    }

    // Graceful shutdown minimal
    go func() {
        log.Printf("Listening on %s", cfg.HTTPAddr)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("server error: %v", err)
        }
    }()

    // Wait / shutdown signals omitted for brevity...
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    _ = srv.Shutdown(ctx)
}
```

---

### `internal/pkg/config/config.go`

```go
package config

import "os"

type Config struct {
    HTTPAddr    string
    DatabaseURL string
    LogLevel    string
}

func Load() *Config {
    // simple: lire depuis env ; remplacer par viper if needed
    return &Config{
        HTTPAddr:    getEnv("HTTP_ADDR", ":8080"),
        DatabaseURL: getEnv("DATABASE_URL", "postgres://user:pass@localhost:5432/mydb?sslmode=disable"),
        LogLevel:    getEnv("LOG_LEVEL", "info"),
    }
}

func getEnv(k, def string) string {
    if v := os.Getenv(k); v != "" {
        return v
    }
    return def
}
```

---

### `internal/pkg/db/db.go`

```go
package db

import (
    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"
)

func NewSQLX(dsn string) (*sqlx.DB, error) {
    db, err := sqlx.Connect("postgres", dsn)
    if err != nil {
        return nil, err
    }
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(5)
    return db, nil
}
```

---

### `internal/domain/user.go`

```go
package domain

import "time"

type User struct {
    ID        int64     `db:"id" json:"id"`
    Email     string    `db:"email" json:"email"`
    Name      string    `db:"name" json:"name"`
    CreatedAt time.Time `db:"created_at" json:"created_at"`
}
```

---

### `internal/app/repository/user_repo.go` (interface)

```go
package repository

import "myapp/internal/domain"

type UserRepository interface {
    Create(user *domain.User) error
    GetByID(id int64) (*domain.User, error)
    GetByEmail(email string) (*domain.User, error)
    List(limit, offset int) ([]*domain.User, error)
}
```

---

### `internal/app/repository/user_repo_sql.go` (impl sqlx)

```go
package repository

import (
    "github.com/jmoiron/sqlx"
    "myapp/internal/domain"
)

type userSQLRepo struct {
    db *sqlx.DB
}

func NewUserSQLRepo(db *sqlx.DB) UserRepository {
    return &userSQLRepo{db: db}
}

func (r *userSQLRepo) Create(u *domain.User) error {
    query := `INSERT INTO users (email, name, created_at) VALUES ($1, $2, $3) RETURNING id`
    return r.db.QueryRowx(query, u.Email, u.Name, u.CreatedAt).Scan(&u.ID)
}

func (r *userSQLRepo) GetByID(id int64) (*domain.User, error) {
    var u domain.User
    err := r.db.Get(&u, "SELECT id, email, name, created_at FROM users WHERE id=$1", id)
    if err != nil {
        return nil, err
    }
    return &u, nil
}

func (r *userSQLRepo) GetByEmail(email string) (*domain.User, error) {
    var u domain.User
    err := r.db.Get(&u, "SELECT id, email, name, created_at FROM users WHERE email=$1", email)
    if err != nil {
        return nil, err
    }
    return &u, nil
}

func (r *userSQLRepo) List(limit, offset int) ([]*domain.User, error) {
    users := []*domain.User{}
    err := r.db.Select(&users, "SELECT id, email, name, created_at FROM users ORDER BY id LIMIT $1 OFFSET $2", limit, offset)
    if err != nil {
        return nil, err
    }
    return users, nil
}
```

---

### `internal/app/service/user_service.go`

```go
package service

import (
    "errors"
    "time"

    "myapp/internal/app/repository"
    "myapp/internal/domain"
)

var ErrUserExists = errors.New("user already exists")

type UserService interface {
    CreateUser(email, name string) (*domain.User, error)
    GetUser(id int64) (*domain.User, error)
}

type userService struct {
    repo repository.UserRepository
    log  Logger // abstraction simple pour logging
}

func NewUserService(r repository.UserRepository, l Logger) UserService {
    return &userService{repo: r, log: l}
}

func (s *userService) CreateUser(email, name string) (*domain.User, error) {
    if ex, _ := s.repo.GetByEmail(email); ex != nil {
        return nil, ErrUserExists
    }
    u := &domain.User{
        Email:     email,
        Name:      name,
        CreatedAt: time.Now().UTC(),
    }
    if err := s.repo.Create(u); err != nil {
        return nil, err
    }
    s.log.Infof("user created id=%d", u.ID)
    return u, nil
}

func (s *userService) GetUser(id int64) (*domain.User, error) {
    return s.repo.GetByID(id)
}
```

*(`Logger` est une interface simple que tu implémente via `internal/pkg/logger` — facilite le test.)*

---

### `internal/app/http/router.go`

```go
package http

import (
    "github.com/gin-gonic/gin"
    "myapp/internal/app/service"
    "myapp/internal/pkg/logger"
)

func NewRouter(uSvc service.UserService, l logger.Logger) *gin.Engine {
    r := gin.New()
    r.Use(gin.Recovery(), gin.Logger())

    uh := NewUserHandler(uSvc, l)

    v1 := r.Group("/v1")
    {
        users := v1.Group("/users")
        users.POST("", uh.Create)
        users.GET("/:id", uh.GetByID)
        users.GET("", uh.List)
    }
    return r
}
```

---

### `internal/app/http/user_handler.go`

```go
package http

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "myapp/internal/app/service"
    "myapp/internal/pkg/logger"
)

type UserHandler struct {
    svc service.UserService
    log logger.Logger
}

func NewUserHandler(s service.UserService, l logger.Logger) *UserHandler {
    return &UserHandler{svc: s, log: l}
}

type createUserReq struct {
    Email string `json:"email" binding:"required,email"`
    Name  string `json:"name" binding:"required"`
}

func (h *UserHandler) Create(c *gin.Context) {
    var req createUserReq
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    u, err := h.svc.CreateUser(req.Email, req.Name)
    if err != nil {
        if err == service.ErrUserExists {
            c.JSON(http.StatusConflict, gin.H{"error": "user exists"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, u)
}

func (h *UserHandler) GetByID(c *gin.Context) {
    idS := c.Param("id")
    id, err := strconv.ParseInt(idS, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }
    u, err := h.svc.GetUser(id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
        return
    }
    c.JSON(http.StatusOK, u)
}

func (h *UserHandler) List(c *gin.Context) {
    // ex. parse query params; omitted for brièveté
    c.JSON(http.StatusOK, gin.H{"users": []string{}})
}
```

---

# 4. Tests & Mocks

* Crée une interface `UserRepository` (déjà fait) ; génère un mock pour les tests de service :

  * `//go:generate mockgen -destination=../mocks/mock_userrepo.go -package=mocks myapp/internal/app/repository UserRepository`
* Tests unitaires :

  * Tests du service : injecter mock repository.
  * Tests du handler : utiliser `httptest` + router Gin.
* Tests d'intégration : containerized Postgres (testcontainers) ou base dédiée ; exécuter migrations avant les tests.

---

# 5. Points d'amélioration / options

* Auth (JWT) : middleware séparé `internal/app/http/middleware/auth.go`.
* Validation : utiliser `github.com/go-playground/validator/v10`.
* Transactions : exposition d'une interface `TxManager` / `WithTx` dans repository.
* DI avancé : `wire` ou `fx` pour projets très gros.
* Observabilité : expose metrics Prometheus, tracing OpenTelemetry.
* Feature flags / config remote : `consul` / `vault` / `envconfig` selon besoin.

---

# 6. Exemple de `Dockerfile` simple

```dockerfile
FROM golang:1.21-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/myapp ./cmd/server

FROM gcr.io/distroless/static
COPY --from=build /bin/myapp /bin/myapp
EXPOSE 8080
ENTRYPOINT ["/bin/myapp"]
```

---

# 7. Checklist de production (rapide)

* Config sécurisée (secrets, env, vault)
* Migrations (eg. golang-migrate)
* Logging structuré, niveaux
* Graceful shutdown & health/readiness endpoints
* Rate limiting/ CORS/Helmet
* CI: tests unitaires, lint (gofmt/golangci-lint), build image, scan vuln

---

Si tu veux, je peux maintenant :

* te générer un **squelette de projet complet** (fichiers `main.go`, handlers, service, repo) prêt à cloner, ou
* t'afficher **exemples de tests unitaires** pour la couche service + mocks, ou
* adapter l'architecture exactement à **postgres + sqlx + squirrel** (vu ton historique).

Dis-moi ce que tu veux que je produise tout de suite (code complet, tests, ou CI config) — je génère les fichiers ici directement.

Squelette complet

Excellent 👍
Voici donc un **squelette complet d’application Go modulaire**, avec architecture claire (repository / service / handler), incluant **exemples de tests unitaires** et **mocks générés** via `mockgen`.

C’est une base prête à cloner pour une API REST en **Go + Gin + sqlx + PostgreSQL**.

---

## 🧩 1. Arborescence du projet

```
myapp/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── app/
│   │   ├── http/
│   │   │   ├── router.go
│   │   │   └── user_handler.go
│   │   ├── repository/
│   │   │   ├── user_repo.go
│   │   │   └── user_repo_sql.go
│   │   ├── service/
│   │   │   └── user_service.go
│   │   └── mocks/
│   │       └── mock_userrepo.go
│   ├── domain/
│   │   └── user.go
│   └── pkg/
│       ├── config/
│       │   └── config.go
│       ├── db/
│       │   └── db.go
│       └── logger/
│           └── logger.go
├── tests/
│   ├── service_user_test.go
│   └── handler_user_test.go
├── go.mod
├── go.sum
└── Makefile
```

---

## ⚙️ 2. Fichiers principaux

### `cmd/server/main.go`

```go
package main

import (
    "log"

    "myapp/internal/app/http"
    "myapp/internal/app/repository"
    "myapp/internal/app/service"
    "myapp/internal/pkg/config"
    "myapp/internal/pkg/db"
    "myapp/internal/pkg/logger"
)

func main() {
    cfg := config.Load()
    logg := logger.New(cfg.LogLevel)

    sqlxDB, err := db.NewSQLX(cfg.DatabaseURL)
    if err != nil {
        log.Fatalf("DB init failed: %v", err)
    }
    defer sqlxDB.Close()

    userRepo := repository.NewUserSQLRepo(sqlxDB)
    userSvc := service.NewUserService(userRepo, logg)
    router := http.NewRouter(userSvc, logg)

    log.Printf("Server running on %s", cfg.HTTPAddr)
    if err := router.Run(cfg.HTTPAddr); err != nil {
        log.Fatalf("Server error: %v", err)
    }
}
```

---

### `internal/pkg/config/config.go`

```go
package config

import "os"

type Config struct {
    HTTPAddr    string
    DatabaseURL string
    LogLevel    string
}

func Load() *Config {
    return &Config{
        HTTPAddr:    getEnv("HTTP_ADDR", ":8080"),
        DatabaseURL: getEnv("DATABASE_URL", "postgres://user:pass@localhost:5432/myapp?sslmode=disable"),
        LogLevel:    getEnv("LOG_LEVEL", "info"),
    }
}

func getEnv(key, def string) string {
    if val := os.Getenv(key); val != "" {
        return val
    }
    return def
}
```

---

### `internal/pkg/db/db.go`

```go
package db

import (
    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"
)

func NewSQLX(dsn string) (*sqlx.DB, error) {
    return sqlx.Connect("postgres", dsn)
}
```

---

### `internal/pkg/logger/logger.go`

```go
package logger

import "log"

type Logger interface {
    Infof(format string, args ...any)
    Errorf(format string, args ...any)
}

type stdLogger struct{}

func New(level string) Logger {
    return &stdLogger{}
}

func (l *stdLogger) Infof(format string, args ...any)  { log.Printf("[INFO] "+format, args...) }
func (l *stdLogger) Errorf(format string, args ...any) { log.Printf("[ERROR] "+format, args...) }
```

---

### `internal/domain/user.go`

```go
package domain

import "time"

type User struct {
    ID        int64     `db:"id" json:"id"`
    Email     string    `db:"email" json:"email"`
    Name      string    `db:"name" json:"name"`
    CreatedAt time.Time `db:"created_at" json:"created_at"`
}
```

---

### `internal/app/repository/user_repo.go`

```go
package repository

import "myapp/internal/domain"

type UserRepository interface {
    Create(u *domain.User) error
    GetByID(id int64) (*domain.User, error)
    GetByEmail(email string) (*domain.User, error)
}
```

---

### `internal/app/repository/user_repo_sql.go`

```go
package repository

import (
    "github.com/jmoiron/sqlx"
    "myapp/internal/domain"
)

type userSQLRepo struct {
    db *sqlx.DB
}

func NewUserSQLRepo(db *sqlx.DB) UserRepository {
    return &userSQLRepo{db: db}
}

func (r *userSQLRepo) Create(u *domain.User) error {
    return r.db.QueryRowx(`INSERT INTO users (email, name, created_at)
                           VALUES ($1, $2, $3) RETURNING id`,
        u.Email, u.Name, u.CreatedAt).Scan(&u.ID)
}

func (r *userSQLRepo) GetByID(id int64) (*domain.User, error) {
    var u domain.User
    if err := r.db.Get(&u, `SELECT * FROM users WHERE id=$1`, id); err != nil {
        return nil, err
    }
    return &u, nil
}

func (r *userSQLRepo) GetByEmail(email string) (*domain.User, error) {
    var u domain.User
    if err := r.db.Get(&u, `SELECT * FROM users WHERE email=$1`, email); err != nil {
        return nil, err
    }
    return &u, nil
}
```

---

### `internal/app/service/user_service.go`

```go
package service

import (
    "errors"
    "time"

    "myapp/internal/app/repository"
    "myapp/internal/domain"
    "myapp/internal/pkg/logger"
)

var ErrUserExists = errors.New("user already exists")

type UserService interface {
    CreateUser(email, name string) (*domain.User, error)
    GetUser(id int64) (*domain.User, error)
}

type userService struct {
    repo repository.UserRepository
    log  logger.Logger
}

func NewUserService(r repository.UserRepository, l logger.Logger) UserService {
    return &userService{repo: r, log: l}
}

func (s *userService) CreateUser(email, name string) (*domain.User, error) {
    if ex, _ := s.repo.GetByEmail(email); ex != nil {
        return nil, ErrUserExists
    }
    u := &domain.User{Email: email, Name: name, CreatedAt: time.Now().UTC()}
    if err := s.repo.Create(u); err != nil {
        return nil, err
    }
    s.log.Infof("created user %d", u.ID)
    return u, nil
}

func (s *userService) GetUser(id int64) (*domain.User, error) {
    return s.repo.GetByID(id)
}
```

---

### `internal/app/http/router.go`

```go
package http

import (
    "github.com/gin-gonic/gin"
    "myapp/internal/app/service"
    "myapp/internal/pkg/logger"
)

func NewRouter(uSvc service.UserService, log logger.Logger) *gin.Engine {
    r := gin.Default()
    uh := NewUserHandler(uSvc, log)

    v1 := r.Group("/v1")
    {
        users := v1.Group("/users")
        users.POST("", uh.Create)
        users.GET("/:id", uh.GetByID)
    }
    return r
}
```

---

### `internal/app/http/user_handler.go`

```go
package http

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "myapp/internal/app/service"
    "myapp/internal/pkg/logger"
)

type UserHandler struct {
    svc service.UserService
    log logger.Logger
}

func NewUserHandler(s service.UserService, l logger.Logger) *UserHandler {
    return &UserHandler{svc: s, log: l}
}

type createUserReq struct {
    Email string `json:"email" binding:"required,email"`
    Name  string `json:"name" binding:"required"`
}

func (h *UserHandler) Create(c *gin.Context) {
    var req createUserReq
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    u, err := h.svc.CreateUser(req.Email, req.Name)
    if err != nil {
        if err == service.ErrUserExists {
            c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, u)
}

func (h *UserHandler) GetByID(c *gin.Context) {
    id, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }

    u, err := h.svc.GetUser(id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
        return
    }
    c.JSON(http.StatusOK, u)
}
```

---

## 🧪 3. Mocks & Tests

### `internal/app/mocks/mock_userrepo.go`

Généré avec :

```bash
go install github.com/golang/mock/mockgen@latest
mockgen -source=internal/app/repository/user_repo.go \
        -destination=internal/app/mocks/mock_userrepo.go \
        -package=mocks
```

---

### `tests/service_user_test.go`

```go
package tests

import (
    "errors"
    "testing"
    "time"

    "github.com/golang/mock/gomock"
    "myapp/internal/app/mocks"
    "myapp/internal/app/service"
    "myapp/internal/domain"
    "myapp/internal/pkg/logger"
)

type dummyLogger struct{}

func (d *dummyLogger) Infof(format string, args ...any)  {}
func (d *dummyLogger) Errorf(format string, args ...any) {}

func TestCreateUser_Success(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := mocks.NewMockUserRepository(ctrl)
    log := &dummyLogger{}

    mockRepo.EXPECT().GetByEmail("a@b.com").Return(nil, errors.New("not found"))
    mockRepo.EXPECT().Create(gomock.Any()).DoAndReturn(func(u *domain.User) error {
        u.ID = 1
        return nil
    })

    svc := service.NewUserService(mockRepo, log)
    u, err := svc.CreateUser("a@b.com", "Alice")
    if err != nil {
        t.Fatalf("unexpected err: %v", err)
    }
    if u.ID != 1 {
        t.Errorf("expected ID 1, got %d", u.ID)
    }
}

func TestCreateUser_AlreadyExists(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := mocks.NewMockUserRepository(ctrl)
    log := &dummyLogger{}

    mockRepo.EXPECT().GetByEmail("a@b.com").Return(&domain.User{ID: 1}, nil)

    svc := service.NewUserService(mockRepo, log)
    _, err := svc.CreateUser("a@b.com", "Alice")
    if err != service.ErrUserExists {
        t.Errorf("expected ErrUserExists, got %v", err)
    }
}
```

---

### `tests/handler_user_test.go`

```go
package tests

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "time"

    "github.com/gin-gonic/gin"
    "myapp/internal/app/http"
    "myapp/internal/app/service"
    "myapp/internal/domain"
)

type fakeUserService struct{}

func (f *fakeUserService) CreateUser(email, name string) (*domain.User, error) {
    return &domain.User{ID: 99, Email: email, Name: name, CreatedAt: time.Now()}, nil
}
func (f *fakeUserService) GetUser(id int64) (*domain.User, error) {
    return &domain.User{ID: id, Email: "t@t.com", Name: "Test", CreatedAt: time.Now()}, nil
}

func TestUserHandler_Create(t *testing.T) {
    gin.SetMode(gin.TestMode)
    svc := &fakeUserService{}
    log := &dummyLogger{}
    router := gin.New()
    uh := http.NewUserHandler(svc, log)
    router.POST("/users", uh.Create)

    body := map[string]string{"email": "test@test.com", "name": "Test"}
    b, _ := json.Marshal(body)
    req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(b))
    req.Header.Set("Content-Type", "application/json")

    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    if w.Code != http.StatusCreated {
        t.Fatalf("expected 201, got %d", w.Code)
    }
}
```

---

## 🧱 4. Makefile

```makefile
test:
	go test ./... -cover

mocks:
	mockgen -source=internal/app/repository/user_repo.go \
		-destination=internal/app/mocks/mock_userrepo.go \
		-package=mocks
```

---

## ✅ 5. Lancer et tester

```bash
go mod tidy
make mocks
make test
go run ./cmd/server
```

Serveur dispo sur `http://localhost:8080/v1/users`.

---

Souhaites-tu que je t’ajoute :

* ✅ une **version Docker + docker-compose (Postgres)** pour développement local,
  ou
* ✅ une **version avec tests d’intégration** (exécution réelle sur DB avec `testcontainers-go`) ?

Parfait — je t’ai préparé **les deux** :

1. une configuration **Docker + docker-compose** pour dev local (Postgres + build de l’app),
2. des **tests d’intégration** en Go utilisant **testcontainers-go** (démarrage d’un container Postgres, application des migrations et exécution des tests).

Je fournis tous les fichiers et commandes nécessaires. Tu peux copier-coller directement.

---

# A. Docker / docker-compose (dev)

### `Dockerfile` (reprise, build statique)

```dockerfile
# build stage
FROM golang:1.21-alpine AS builder
WORKDIR /src
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /bin/myapp ./cmd/server

# runtime
FROM gcr.io/distroless/static
COPY --from=builder /bin/myapp /bin/myapp
EXPOSE 8080
ENTRYPOINT ["/bin/myapp"]
```

---

### `docker-compose.yml`

```yaml
version: "3.8"
services:
  db:
    image: postgres:15
    environment:
      POSTGRES_USER: myapp
      POSTGRES_PASSWORD: myapppw
      POSTGRES_DB: myapp_test
    ports:
      - "5432:5432"
    volumes:
      - db_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U myapp"]
      interval: 5s
      timeout: 5s
      retries: 5

  app:
    build: .
    environment:
      DATABASE_URL: postgres://myapp:myapppw@db:5432/myapp_test?sslmode=disable
      HTTP_ADDR: ":8080"
    ports:
      - "8080:8080"
    depends_on:
      db:
        condition: service_healthy

volumes:
  db_data:
```

---

### `.env` (optionnel)

```
HTTP_ADDR=:8080
DATABASE_URL=postgres://myapp:myapppw@localhost:5432/myapp_test?sslmode=disable
LOG_LEVEL=debug
```

---

### migration SQL (ex. `migrations/001_create_users.sql`)

```sql
CREATE TABLE IF NOT EXISTS users (
  id BIGSERIAL PRIMARY KEY,
  email TEXT NOT NULL UNIQUE,
  name TEXT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
```

---

### Commandes utiles

* Démarrer la stack : `docker-compose up --build`
* Appliquer migrations (local) : utilise `psql` ou un outil `migrate` ; ex :

  ```bash
  docker-compose exec db psql -U myapp -d myapp_test -c "CREATE TABLE IF NOT EXISTS users (id bigserial primary key, email text unique not null, name text not null, created_at timestamptz default now());"
  ```
* Accéder à l’API : `http://localhost:8080/v1/users`

---

# B. Tests d’intégration (testcontainers-go)

Ajoute ces dépendances :

```bash
go get github.com/testcontainers/testcontainers-go
go get github.com/stretchr/testify
go get github.com/jmoiron/sqlx
go get github.com/lib/pq
```

---

### `tests/integration_test.go`

```go
package tests

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/jmoiron/sqlx"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/stretchr/testify/require"
)

func setupPostgresContainer(ctx context.Context, t *testing.T) (tc.Container, string, func()) {
	req := tc.ContainerRequest{
		Image:        "postgres:15",
		Env: map[string]string{
			"POSTGRES_USER":     "itest",
			"POSTGRES_PASSWORD": "itestpw",
			"POSTGRES_DB":       "itestdb",
		},
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp").WithStartupTimeout(30 * time.Second),
	}
	postgresC, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	host, err := postgresC.Host(ctx)
	require.NoError(t, err)
	p, err := postgresC.MappedPort(ctx, "5432")
	require.NoError(t, err)

	dsn := fmt.Sprintf("postgres://itest:itestpw@%s:%s/itestdb?sslmode=disable", host, p.Port())
	// cleanup function
	cleanup := func() {
		_ = postgresC.Terminate(ctx)
	}

	return postgresC, dsn, cleanup
}

func runMigrations(t *testing.T, db *sqlx.DB) {
	// simple ad-hoc migration for tests
	_, err := db.Exec(`
CREATE TABLE IF NOT EXISTS users (
  id BIGSERIAL PRIMARY KEY,
  email TEXT NOT NULL UNIQUE,
  name TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
`)
	require.NoError(t, err)
}

func TestIntegration_CreateAndGetUser(t *testing.T) {
	ctx := context.Background()
	_, dsn, cleanup := setupPostgresContainer(ctx, t)
	defer cleanup()

	// wait a little for postgres to be fully ready
	time.Sleep(500 * time.Millisecond)

	db, err := sqlx.Connect("postgres", dsn)
	require.NoError(t, err)
	defer db.Close()

	runMigrations(t, db)

	// Insert a user (imitate repository.Create)
	now := time.Now().UTC()
	var id int64
	err = db.QueryRowx(
		`INSERT INTO users (email, name, created_at) VALUES ($1, $2, $3) RETURNING id`,
		"integration@local", "Integration", now,
	).Scan(&id)
	require.NoError(t, err)
	require.True(t, id > 0)

	// Get user back (imitate repository.GetByID)
	var got struct {
		ID        int64     `db:"id"`
		Email     string    `db:"email"`
		Name      string    `db:"name"`
		CreatedAt time.Time `db:"created_at"`
	}
	err = db.Get(&got, "SELECT id, email, name, created_at FROM users WHERE id=$1", id)
	require.NoError(t, err)
	require.Equal(t, "integration@local", got.Email)
	require.Equal(t, "Integration", got.Name)
}
```

Explications :

* Le test démarre un container Postgres via `testcontainers-go`.
* Il applique une migration minimale (création de la table).
* Il exécute un insert et une lecture pour vérifier le fonctionnement réel du repository/db.

---

# C. Makefile (mis à jour)

```makefile
test:
	go test ./... -cover

mocks:
	mockgen -source=internal/app/repository/user_repo.go \
		-destination=internal/app/mocks/mock_userrepo.go \
		-package=mocks

docker-up:
	docker-compose up --build

docker-down:
	docker-compose down -v

integration-test:
	go test -v ./tests -run TestIntegration -count=1
```

---

# D. Conseils pratiques / notes

* Sur CI (GitHub Actions), tu peux lancer `docker-compose` ou utiliser `services: postgres` pour les tests d’intégration. Pour testcontainers, CI doit autoriser Docker (GitHub Actions runners le supportent).
* Si tu préfères exécuter migrations SQL automatiquement à partir d’un dossier `migrations/`, remplace `runMigrations` par l’appel à `migrate` (github.com/golang-migrate/migrate/v4).
* Dans les tests d’intégration, vérifie la suppression du container via la fonction `cleanup` (déjà incluse).
* Pour accélérer le développement local, utilise une base Postgres locale (via docker-compose) et active un volume pour préserver l’état.

---

Si tu veux je peux maintenant (choisis une action ; je la fais tout de suite) :

* générer **le repo complet** (tous fichiers) dans un seul archive `.zip` prêt à downloader,
* ou **donner le fichier GitHub Actions** pour exécuter tests unitaires + integration tests via `docker-compose`,
* ou **ajouter `golang-migrate`** et montrer comment intégrer ses migrations dans `cmd/server/main.go`.

Que souhaites-tu que je fournisse en plus tout de suite ?

Excellent 👌 — ajoutons **`golang-migrate`** proprement à ton architecture Go modulaire.
Tu auras :
✅ migrations versionnées dans `/migrations`,
✅ intégration automatique dans `cmd/server/main.go`,
✅ un `Makefile` et commandes CLI pratiques,
✅ option d’exécution via Docker ou à la main.

---

## 🧩 1. Installation & dépendances

### Ajoute dans `go.mod`

```bash
go get -u github.com/golang-migrate/migrate/v4
go get -u github.com/golang-migrate/migrate/v4/database/postgres
go get -u github.com/golang-migrate/migrate/v4/source/file
```

---

## 📂 2. Arborescence modifiée

```
myapp/
├── cmd/server/main.go
├── internal/
│   ├── ...
├── migrations/
│   ├── 001_create_users.up.sql
│   └── 001_create_users.down.sql
├── Makefile
├── docker-compose.yml
└── go.mod
```

---

## 🧾 3. Exemple de migrations SQL

### `migrations/001_create_users.up.sql`

```sql
CREATE TABLE IF NOT EXISTS users (
  id BIGSERIAL PRIMARY KEY,
  email TEXT NOT NULL UNIQUE,
  name TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

### `migrations/001_create_users.down.sql`

```sql
DROP TABLE IF EXISTS users;
```

---

## ⚙️ 4. Ajout d’un module `internal/pkg/migrate/migrate.go`

```go
package migrate

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(databaseURL, migrationsPath string) error {
	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsPath),
		databaseURL,
	)
	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}

	defer func() {
		srcErr, dbErr := m.Close()
		if srcErr != nil || dbErr != nil {
			log.Printf("warning: closing migrate: src=%v db=%v", srcErr, dbErr)
		}
	}()

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration failed: %w", err)
	}

	log.Println("✅ Migrations applied successfully")
	return nil
}
```

---

## 🚀 5. Intégration dans `cmd/server/main.go`

On appelle la fonction `RunMigrations` avant de démarrer l’API :

```go
package main

import (
    "log"

    "myapp/internal/app/http"
    "myapp/internal/app/repository"
    "myapp/internal/app/service"
    "myapp/internal/pkg/config"
    "myapp/internal/pkg/db"
    "myapp/internal/pkg/logger"
    "myapp/internal/pkg/migrate"
)

func main() {
    cfg := config.Load()
    logg := logger.New(cfg.LogLevel)

    // 1️⃣ Connexion DB
    sqlxDB, err := db.NewSQLX(cfg.DatabaseURL)
    if err != nil {
        log.Fatalf("DB init failed: %v", err)
    }
    defer sqlxDB.Close()

    // 2️⃣ Lancer migrations avant de servir l’API
    if err := migrate.RunMigrations(cfg.DatabaseURL, "migrations"); err != nil {
        log.Fatalf("Migration failed: %v", err)
    }

    // 3️⃣ DI
    userRepo := repository.NewUserSQLRepo(sqlxDB)
    userSvc := service.NewUserService(userRepo, logg)
    router := http.NewRouter(userSvc, logg)

    log.Printf("Server running on %s", cfg.HTTPAddr)
    if err := router.Run(cfg.HTTPAddr); err != nil {
        log.Fatalf("Server error: %v", err)
    }
}
```

---

## 🧰 6. Mise à jour du `Makefile`

```makefile
test:
	go test ./... -cover

mocks:
	mockgen -source=internal/app/repository/user_repo.go \
		-destination=internal/app/mocks/mock_userrepo.go \
		-package=mocks

docker-up:
	docker-compose up --build -d

docker-down:
	docker-compose down -v

# Run migrations locally using golang-migrate CLI (optional)
migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down

integration-test:
	go test -v ./tests -run TestIntegration -count=1
```

⚠️ Pour ces cibles `migrate-up` et `migrate-down`, installe `migrate` CLI :

```bash
brew install golang-migrate
# ou
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

---

## 🐳 7. Docker-compose (rappel)

Le container `app` exécutera les migrations automatiquement au démarrage (via le code Go dans `main.go`), donc pas besoin d’un container séparé pour `migrate`.

---

## ✅ 8. Exécution complète

```bash
# local dev
go run ./cmd/server

# ou avec docker
docker-compose up --build

# ou appliquer migrations manuellement
make migrate-up
```

---

## 🧪 9. Exemple de logs lors du démarrage

```
2025/10/16 10:32:41 ✅ Migrations applied successfully
2025/10/16 10:32:41 Server running on :8080
```

---

## ✨ Bonus — migrations versionnées auto

Pour créer facilement de nouvelles migrations :

```bash
migrate create -ext sql -dir migrations -seq add_phone_to_user
```

Cela créera :

```
migrations/
├── 001_create_users.up.sql
├── 001_create_users.down.sql
├── 002_add_phone_to_user.up.sql
└── 002_add_phone_to_user.down.sql
```

---

Souhaites-tu que je t’ajoute **le code du workflow GitHub Actions CI/CD** pour :

* lancer les tests unitaires,
* démarrer Postgres (service docker),
* exécuter `golang-migrate up`,
* puis lancer les tests d’intégration ?

Parfait 💪 — voici une **configuration GitLab CI/CD complète**
pour ton projet Go modulaire avec PostgreSQL, `golang-migrate`, et tests unitaires + d’intégration.

Elle gère :
✅ build + lint + test unitaire
✅ démarrage d’un conteneur PostgreSQL
✅ exécution des migrations via `golang-migrate`
✅ tests d’intégration réels (avec `testcontainers-go` ou DB réelle)

---

# 🧩 `.gitlab-ci.yml` complet

```yaml
stages:
  - lint
  - test
  - build

variables:
  GO_VERSION: "1.21"
  POSTGRES_DB: myapp_test
  POSTGRES_USER: myapp
  POSTGRES_PASSWORD: myapppw
  DATABASE_URL: postgres://myapp:myapppw@postgres:5432/myapp_test?sslmode=disable
  GOPROXY: https://proxy.golang.org
  CGO_ENABLED: "0"

services:
  - name: postgres:15
    alias: postgres

before_script:
  - echo "Using Go version $GO_VERSION"
  - apt-get update -qq && apt-get install -y curl unzip
  - go version
  - go mod download

# --- Stage 1: Linter ---
lint:
  stage: lint
  image: golang:${GO_VERSION}-alpine
  script:
    - go install golang.org/x/lint/golint@latest
    - golint ./...
    - go vet ./...
  allow_failure: true
  only:
    - merge_requests
    - main

# --- Stage 2: Unit tests ---
unit_tests:
  stage: test
  image: golang:${GO_VERSION}-alpine
  script:
    - go test ./internal/... -v -coverprofile=coverage.out
    - go tool cover -func=coverage.out | grep total
  artifacts:
    when: always
    reports:
      junit: junit.xml
    paths:
      - coverage.out
  only:
    - merge_requests
    - main

# --- Stage 3: Integration tests ---
integration_tests:
  stage: test
  image: golang:${GO_VERSION}
  services:
    - name: postgres:15
      alias: postgres
  variables:
    DATABASE_URL: postgres://myapp:myapppw@postgres:5432/myapp_test?sslmode=disable
  before_script:
    - apt-get update -qq && apt-get install -y wget
    - wget -q https://github.com/golang-migrate/migrate/releases/download/v4.17.1/migrate.linux-amd64.tar.gz
    - tar -xzf migrate.linux-amd64.tar.gz -C /usr/local/bin
    - chmod +x /usr/local/bin/migrate
    - go mod download
    - until pg_isready -h postgres -U $POSTGRES_USER; do echo "Waiting for postgres..."; sleep 1; done
    # appliquer migrations
    - migrate -path migrations -database "$DATABASE_URL" up
  script:
    - go test ./tests -v -count=1
  only:
    - merge_requests
    - main

# --- Stage 4: Build ---
build:
  stage: build
  image: golang:${GO_VERSION}
  script:
    - go build -o bin/myapp ./cmd/server
  artifacts:
    paths:
      - bin/myapp
  only:
    - main
```

---

# 📦 Explications par étape

| Stage                 | Description                                                                                                                               |
| --------------------- | ----------------------------------------------------------------------------------------------------------------------------------------- |
| **lint**              | Vérifie la qualité du code avec `golint` et `go vet`.                                                                                     |
| **unit_tests**        | Exécute les tests unitaires avec couverture (`go test ./internal/...`).                                                                   |
| **integration_tests** | Lance PostgreSQL dans un service GitLab, applique les migrations avec `golang-migrate`, puis exécute les tests d’intégration (`./tests`). |
| **build**             | Compile ton binaire Go (production-ready).                                                                                                |

---

# 🧪 Si tu utilises `testcontainers-go`

Tu n’as pas besoin de `golang-migrate` dans le pipeline CI, car les tests créent leur propre container Postgres.
Mais il faut **autoriser Docker-in-Docker** dans ton runner GitLab :

```yaml
integration_tests:
  stage: test
  image: docker:24.0.6
  services:
    - docker:dind
  variables:
    DOCKER_DRIVER: overlay2
  script:
    - apk add --no-cache go bash git
    - go mod download
    - go test ./tests -v -count=1
```

👉 Cette version exécute réellement `testcontainers-go` (les containers Postgres sont gérés directement dans les tests).

---

# ⚙️ ️Exemple de structure complète

```
.
├── .gitlab-ci.yml
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── app/
│   ├── pkg/
│   └── ...
├── migrations/
│   ├── 001_create_users.up.sql
│   └── 001_create_users.down.sql
├── tests/
│   └── integration_test.go
├── go.mod
└── go.sum
```

---

# ✅ Exemple de logs GitLab CI

```
$ migrate -path migrations -database "$DATABASE_URL" up
✅ Migrations applied successfully
$ go test ./tests -v
=== RUN   TestIntegration_CreateAndGetUser
--- PASS: TestIntegration_CreateAndGetUser (1.23s)
PASS
ok      myapp/tests  1.9s
```

---

# 🚀 Pour aller plus loin

* Ajoute un **job “deploy”** après le build si tu veux déployer via Docker ou Kubernetes.
* Tu peux générer automatiquement la couverture HTML :

  ```yaml
  artifacts:
    paths:
      - coverage.html
  ```

  en ajoutant :

  ```bash
  go tool cover -html=coverage.out -o coverage.html
  ```

---

Souhaites-tu que je t’ajoute une **étape de déploiement** (par exemple Docker push → registry.gitlab.com, ou déploiement sur un serveur SSH / Kubernetes) dans ce même `.gitlab-ci.yml` ?
