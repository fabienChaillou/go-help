Excellente question 👏 — bien architecturer un projet **Go** est crucial pour sa maintenabilité, sa testabilité et son évolutivité.

Voici une **référence complète et moderne** pour structurer un projet Go, adaptée à une application type **API REST, microservice ou outil CLI**.

---

## 🧱 1. Structure générale d’un projet Go

Une structure typique et propre pourrait ressembler à ceci :

```
myapp/
├── cmd/
│   └── myapp/
│       └── main.go            # Point d'entrée principal
│
├── internal/
│   ├── api/                   # Couche "controller" (handlers HTTP)
│   │   ├── v1/
│   │   │   ├── user_handler.go
│   │   │   ├── auth_handler.go
│   │   │   └── ...
│   │   └── middleware/
│   │       └── auth_middleware.go
│   │
│   ├── domain/                # Couche "business logic"
│   │   ├── user/
│   │   │   ├── model.go       # Entité métier
│   │   │   ├── service.go     # Logique métier
│   │   │   ├── repository.go  # Interface de persistance
│   │   │   └── mock/          # Mocks pour les tests
│   │   └── ...
│   │
│   ├── repository/            # Couche "data access"
│   │   ├── sqlx/
│   │   │   └── user_repo.go
│   │   ├── mongo/
│   │   │   └── ...
│   │   └── memory/
│   │       └── ...
│   │
│   ├── config/
│   │   ├── config.go
│   │   └── config.yaml
│   │
│   ├── app/
│   │   └── app.go             # Initialisation globale (router, db, etc.)
│   │
│   └── utils/
│       ├── logger.go
│       └── errors.go
│
├── pkg/                       # Code réutilisable (public)
│   └── validation/
│       └── validator.go
│
├── scripts/                   # Outils CI/CD, migration, etc.
│   └── migrate.sh
│
├── migrations/                # SQL migrations
│   ├── 001_init.sql
│   └── 002_add_user.sql
│
├── go.mod
├── go.sum
└── README.md
```

---

## 🧩 2. Découpage logique

| Couche                   | Rôle                                 | Exemple de fichiers                  |
| ------------------------ | ------------------------------------ | ------------------------------------ |
| **cmd/**                 | Point d’entrée (exécutable)          | `main.go`                            |
| **internal/api/**        | Routes HTTP, validation, réponses    | `user_handler.go`, `auth_handler.go` |
| **internal/domain/**     | Logique métier, entités, services    | `user/service.go`, `user/model.go`   |
| **internal/repository/** | Accès à la base de données, DAO      | `user_repo.go`                       |
| **internal/config/**     | Gestion de la configuration          | `config.go`, `.env`, `config.yaml`   |
| **internal/utils/**      | Outils transversaux                  | `logger.go`, `error.go`              |
| **pkg/**                 | Modules réutilisables (hors domaine) | `validation/validator.go`            |

---

## ⚙️ 3. Exemple minimal d’architecture API

### `cmd/myapp/main.go`

```go
package main

import (
    "myapp/internal/app"
    "log"
)

func main() {
    a := app.New()
    if err := a.Run(); err != nil {
        log.Fatal(err)
    }
}
```

### `internal/app/app.go`

```go
package app

import (
    "github.com/gin-gonic/gin"
    "myapp/internal/api/v1"
    "myapp/internal/config"
    "myapp/internal/repository/sqlx"
)

type App struct {
    router *gin.Engine
}

func New() *App {
    cfg := config.Load()
    repo := sqlx.NewRepository(cfg.DatabaseURL)
    userHandler := v1.NewUserHandler(repo)

    r := gin.Default()
    v1.RegisterRoutes(r, userHandler)

    return &App{router: r}
}

func (a *App) Run() error {
    return a.router.Run(":8080")
}
```

### `internal/domain/user/service.go`

```go
package user

type Service struct {
    repo Repository
}

func NewService(r Repository) *Service {
    return &Service{repo: r}
}

func (s *Service) CreateUser(u *User) error {
    // logique métier : validation, hashing, etc.
    return s.repo.Create(u)
}
```

### `internal/repository/sqlx/user_repo.go`

```go
package sqlx

import (
    "github.com/jmoiron/sqlx"
    "myapp/internal/domain/user"
)

type UserRepository struct {
    db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) Create(u *user.User) error {
    _, err := r.db.NamedExec(`INSERT INTO users (name, email) VALUES (:name, :email)`, u)
    return err
}
```

---

## 🧪 4. Organisation des tests

```
internal/
├── domain/
│   ├── user/
│   │   ├── service_test.go
│   │   └── mock/
│   │       └── repository_mock.go
```

Tests unitaires :

```go
func TestCreateUser(t *testing.T) {
    repo := mock.NewRepository()
    svc := user.NewService(repo)

    err := svc.CreateUser(&user.User{Name: "Alice"})
    require.NoError(t, err)
}
```

---

## 🧰 5. Bonnes pratiques

✅ **Respecter la règle "internal"**
→ Empêche l’import de code interne depuis d’autres modules.

✅ **Interfaces dans le domaine**, implémentations dans `repository/`
→ Inversion de dépendances → tests plus faciles.

✅ **Configuration via environnement + fichier YAML**
→ Exemple avec `spf13/viper`.

✅ **Logs structurés (zap, zerolog)**
✅ **Erreurs contextualisées (pkg/errors)**
✅ **Tests + mocks (go test + mockgen)**
✅ **CI/CD (GitHub Actions ou Makefile)**

---

## ⚙️ 6. Exemple Makefile

```makefile
run:
	go run ./cmd/myapp

test:
	go test ./... -cover

build:
	go build -o bin/myapp ./cmd/myapp

lint:
	golangci-lint run

migrate:
	migrate -path migrations -database "$(DB_URL)" up
```

---

