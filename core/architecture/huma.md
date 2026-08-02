Très bonne idée : utiliser **Huma** change un peu l’architecture parce qu’il offre déjà des couches d’abstraction (OpenAPI, validation, etc.). Voici comment je structurerais un projet Go avec Huma, avec les avantages/inconvénients, et un exemple concret de squelette.

---

## 🔍 Rappel rapide : Huma

Premières choses à connaître :

* Huma est un micro-framework REST/RPC pour Go, basé sur **OpenAPI 3.1** et **JSON Schema**. ([Huma][1])
* Il est agnostique au router : tu peux utiliser Chi, Gin, Fiber, ou même `net/http` standard via des « adapters ». ([Huma][2])
* Il gère automatiquement la validation des entrées, des schémas, les erreurs standardisées, la génération de documentation. ([Go Forum][3])

Donc on peut se concentrer sur la logique métier, les modèles, etc. avec moins de boilerplate.

---

## 🧱 Proposition d’architecture avec Huma

Voici une structure que je recommande quand on utilise Huma pour une application type API REST / microservice :

```
myapp/
├── cmd/
│   └── myapp/
│       └── main.go
│
├── internal/
│   ├── api/
│   │   ├── v1/
│   │   │   ├── handlers.go        # Handlers / opérations Huma pour v1
│   │   │   └── models.go          # V1 structs pour requêtes/réponses
│   │   └── middleware/            # Middlewares spécifiques (auth, logging, etc.)
│   │
│   ├── domain/
│   │   ├── user/                  
│   │   │   ├── service.go         # Logique métier
│   │   │   ├── model.go           # Entité métier
│   │   │   ├── repository.go      # Interfaces
│   │   └── ...                    # Autres domaines
│   │
│   ├── repository/                
│   │   ├── sql/                   # Implémentation SQL (Postgres, etc.)
│   │   └── mem/                   # Optionnel : stockage en mémoire pour tests
│   │
│   ├── config/                    
│   │   └── config.go              # Lecture des configs, env, etc.
│   │
│   ├── utils/                     
│   │   ├── logger.go              
│   │   ├── errors.go              
│   │   └── helpers.go             
│   │
│   ├── app/                       
│   │   └── app.go                 # Initialisation de Huma + router + base setup
│   │
├── docs/                          # Documentation, schémas custom, etc.
│   └── openapi/                   
│
├── migrations/                     # SQL migrations
│
├── go.mod
├── go.sum
├── Makefile                       # build, lint, test, run
├── Dockerfile                     # containerisation
└── README.md
```

---

## ⚙️ Composants avec Huma dans cette architecture

Voici comment les pièces s’articulent :

| Élément                     | Avec Huma, rôle / ce que Huma prend en charge                                                                                                                                           | Ce que tu dois fournir                                                                                                                                                                        |
| --------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Router / adapter**        | Choisir un adapter fourni (ex : humachi pour Chi, humago pour `net/http`, etc.). ([Huma][2])                                                                                            | Importer l’adapter, initialiser l’API Huma + router (avec config) dans `app/app.go` ou `main.go`.                                                                                             |
| **Handlers / Opérations**   | Huma permet de déclarer des opérations via `huma.Operation{...}` : méthodes HTTP, chemins, entrées/sorties, réponses, erreurs. Huma génère le schéma et gère la validation. ([Huma][4]) | Dans `internal/api/v1/handlers.go`, définir les handlers / fonctions métiers qui seront appelées depuis les opérations Huma. Créer les modèles d’entrée/sortie.                               |
| **Validation / Schémas**    | Pris en charge automatiquement par Huma via JSON Schema généré à partir des structs Go. ([Go Forum][3])                                                                                 | Annoter les structs avec les tags nécessaires pour validation (ex : `json:"..."`, `huma` ou autres si besoin). Gérer les cas particuliers de validation plus complexes dans la couche métier. |
| **Documentation / OpenAPI** | Huma génère la spec OpenAPI, les schémas, les docs interatives (selon config). ([Huma][4])                                                                                              | Ajuster la configuration (titre, version, chemins des docs), peut-être personnaliser ou étendre le schéma ou les descriptions.                                                                |
| **Gestion des erreurs**     | Huma fournit un modèle d’erreur conforme RFC (par exemple `application/problem+json`) et des réponses d’erreur standardisées. ([Huma][4])                                               | Définir les erreurs métiers, leur mapping, les cas d’erreur dans les handlers/métier.                                                                                                         |
| **Tests**                   | Tester les handlers, la logique métier, possiblement tester les opérations via des requêtes HTTP (local), en mockant les dépôts.                                                        | Mettre en place des tests unitaires / d’intégration.                                                                                                                                          |

---

## 🛠 Exemple de squelette concret avec Huma

Voici un exemple de code minimal pour démarrer :

### `cmd/myapp/main.go`

```go
package main

import (
    "log"
    "net/http"

    "myapp/internal/app"
)

func main() {
    a, err := app.New()
    if err != nil {
        log.Fatal(err)
    }
    // Suppose que app.New retourne aussi le handler http.Handler
    log.Println("Starting server on :8080")
    if err := http.ListenAndServe(":8080", a.Handler()); err != nil {
        log.Fatal(err)
    }
}
```

### `internal/app/app.go`

```go
package app

import (
    "context"
    "myapp/internal/config"
    "myapp/internal/api/v1"
    "myapp/internal/repository/sql"
    "myapp/internal/utils"

    "github.com/danielgtaylor/huma/v2"
    "github.com/danielgtaylor/huma/v2/adapters/humachi"  // Exemple avec Chi
)

type App struct {
    api huma.API
}

func New() (*App, error) {
    cfg, err := config.Load()
    if err != nil {
        return nil, err
    }

    // Setup DB, logger, etc.
    db, err := sql.NewPostgres(cfg.DBURL)
    if err != nil {
        return nil, err
    }

    logger := utils.NewLogger(cfg)

    // Créer le router
    r := /* un router Chi, initialisé selon ton choix */
    // Exemple :
    // mux := chi.NewMux()
    // ou autre

    apiConfig := huma.DefaultConfig("MyApp API", "v1.0.0")
    apiConfig.Servers = []*huma.Server{
        {URL: cfg.ServerURL},
    }
    api := humachi.New(r, apiConfig)

    // On injecte les dépendances nécessaires dans les handlers
    userRepo := sql.NewUserRepo(db)
    userService := v1.NewUserService(userRepo)

    // Enregistrer les opérations
    v1.RegisterUserRoutes(api, userService)

    return &App{api: api}, nil
}

func (a *App) Handler() http.Handler {
    return a.api.Handler()  // ou unwrap pour le router si besoin
}
```

### `internal/api/v1/handlers.go`

```go
package v1

import (
    "context"
    "net/http"

    "myapp/internal/domain/user"

    "github.com/danielgtaylor/huma/v2"
)

func RegisterUserRoutes(api huma.API, svc *user.Service) {
    // Exemple : Get user by ID
    api.Register(huma.Operation{
        OperationID:   "getUserByID",
        Description:   "Get a user by its ID",
        Method:         http.MethodGet,
        Path:           "/users/{id}",
        PathParams:     []*huma.Parameter{huma.PathParameter("id", "User ID")},
        Responses: []*huma.Response{
            huma.OKResponse(&UserResponse{}),
            // éventuellement erreurs
        },
        Handler: func(ctx context.Context, input *struct{
            ID string `path:"id"`
        }) (*UserResponse, error) {
            user, err := svc.GetByID(ctx, input.ID)
            if err != nil {
                return nil, err
            }
            return &UserResponse{
                ID: user.ID,
                Name: user.Name,
                Email: user.Email,
            }, nil
        },
    })

    // Autres routes : list, create, etc.
}
```

### `internal/api/v1/models.go`

```go
package v1

type UserResponse struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

type CreateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}
```

### Domaine / service / repository

Comme dans une architecture classique : la couche `domain/user/service.go` contient la logique métier, `repository/sql/user_repo.go` l’accès DB, etc.

---

## ✅ Avantages / à surveiller avec Huma

**Avantages :**

* Moins de boilerplate pour la validation, la doc, les erreurs.
* Spec OpenAPI toujours à jour (tant que le code est la source de vérité).
* Clair, déclaratif ce que l’API expose via `Operation{}`.
* Intégration possible avec différents routers — flexible.
* Bon pour équipes front-back ou microservices : documentation générée = gains de communication.

**À surveiller / limitations :**

* La courbe d’apprentissage pour bien utiliser les features (schemas, configuration, erreurs).
* Quand la logique métier est complexe, il faut veiller à ne pas tout coller dans le handler. Handler doit rester mince.
* Performance : bien choisir adapter, surveiller overhead (mais Huma semble bien optimisé). ([Huma][5])
* Versioning de l’API : il faudra clairement organiser les versions (v1, v2, …) dans dossiers / chemins.
* Gestion des middlewares (auth, logging, etc.) : bien voir où les attacher, s’ils doivent agir avant ou après Huma, etc.

---

Si tu veux, je peux te proposer une architecture *réelle*, avec Docker, CI, tests, et prête à être clonée, adaptée à ton besoin (par exemple API + auth + base Postgres) ? Tu veux que je fasse ça ?

[1]: https://huma.rocks/?utm_source=chatgpt.com "Huma Introduction - Huma"
[2]: https://huma.rocks/features/bring-your-own-router/?utm_source=chatgpt.com "Bring Your Own Router - Huma"
[3]: https://forum.golangbridge.org/t/huma-v2-13-0-released/35080?utm_source=chatgpt.com "Huma v2.13.0 Released - Releases - Go Forum"
[4]: https://huma.rocks/features/?utm_source=chatgpt.com "Features Overview - Huma"
[5]: https://huma.rocks/why/benchmarks/?utm_source=chatgpt.com "Benchmarks - Huma"

Parfait — j’ai créé un **squelette complet** (fichiers, Docker, CI, tests, migrations, exemples de code) dans un document intitulé **“Squelette-huma-go”** que tu peux consulter et cloner depuis le canvas à côté de la conversation.


# Squelette de projet Go + Huma (prêt à cloner)

Ce document fournit un **squelette complet** prêt à démarrer : code minimal, Docker, docker-compose (Postgres), Makefile, GitHub Actions CI, et tests. Copie/colle les fichiers dans ton repo, `go mod init github.com/you/myapp` et `make run`.

---

## Arborescence recommandée

```
myapp/
├── cmd/myapp/main.go
├── internal/app/app.go
├── internal/api/v1/handlers.go
├── internal/api/v1/models.go
├── internal/domain/user/service.go
├── internal/domain/user/model.go
├── internal/repository/sql/user_repo.go
├── internal/config/config.go
├── internal/utils/logger.go
├── migrations/
├── .env.example
├── Dockerfile
├── docker-compose.yml
├── Makefile
├── go.mod
├── go.sum
├── .github/workflows/ci.yml
└── README.md
```

---

## Fichiers essentiels

### `go.mod` (exemple)

```go
module github.com/you/myapp

go 1.21

require (
    github.com/danielgtaylor/huma/v2 v2.17.0
    github.com/jmoiron/sqlx v1.4.7
    github.com/lib/pq v1.10.10
    github.com/go-chi/chi/v5 v5.1.0
    github.com/stretchr/testify v1.8.4
)
```

(adapte les versions si nécessaire)

---

### `cmd/myapp/main.go`

```go
package main

import (
    "log"
    "net/http"

    "github.com/you/myapp/internal/app"
)

func main() {
    a, err := app.New()
    if err != nil {
        log.Fatalf("app init: %v", err)
    }

    log.Println("Server listening on :8080")
    if err := http.ListenAndServe(":8080", a.Handler()); err != nil {
        log.Fatalf("server: %v", err)
    }
}
```

---

### `internal/app/app.go`

```go
package app

import (
    "github.com/you/myapp/internal/config"
    v1 "github.com/you/myapp/internal/api/v1"
    "github.com/you/myapp/internal/repository/sql"
    "github.com/you/myapp/internal/utils"

    "github.com/danielgtaylor/huma/v2"
    "github.com/danielgtaylor/huma/v2/adapters/humachi"
    "github.com/go-chi/chi/v5"
)

type App struct {
    api huma.API
}

func New() (*App, error) {
    cfg, err := config.Load()
    if err != nil {
        return nil, err
    }

    logger := utils.NewLogger(cfg)

    // DB
    db, err := sql.NewPostgres(cfg.DatabaseURL)
    if err != nil {
        return nil, err
    }

    // Router (chi)
    r := chi.NewRouter()

    // Huma API config
    apiConfig := huma.DefaultConfig("MyApp API", "v1.0.0")
    api := humachi.New(r, apiConfig)

    // deps
    userRepo := sql.NewUserRepo(db)
    userService := v1.NewUserService(userRepo, logger)

    // register routes
    v1.RegisterUserRoutes(api, userService)

    return &App{api: api}, nil
}

func (a *App) Handler() http.Handler {
    return a.api.Handler()
}
```

---

### `internal/api/v1/models.go`

```go
package v1

type UserResponse struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

type CreateUserRequest struct {
    Name  string `json:"name" validate:"required"`
    Email string `json:"email" validate:"required,email"`
}
```

---

### `internal/api/v1/handlers.go`

```go
package v1

import (
    "context"
    "net/http"

    "github.com/danielgtaylor/huma/v2"
    "github.com/you/myapp/internal/domain/user"
)

func RegisterUserRoutes(api huma.API, svc *user.Service) {
    api.Register(huma.Operation{
        OperationID: "createUser",
        Method:      http.MethodPost,
        Path:        "/users",
        Request:     &CreateUserRequest{},
        Responses: []*huma.Response{
            huma.CreatedResponse(&UserResponse{}),
        },
        Handler: func(ctx context.Context, input *CreateUserRequest) (*UserResponse, error) {
            u := &user.User{ Name: input.Name, Email: input.Email }
            if err := svc.Create(ctx, u); err != nil {
                return nil, err
            }
            return &UserResponse{ ID: u.ID, Name: u.Name, Email: u.Email }, nil
        },
    })

    api.Register(huma.Operation{
        OperationID: "getUser",
        Method:      http.MethodGet,
        Path:        "/users/{id}",
        PathParams:  []*huma.Parameter{huma.PathParameter("id", "User ID")},
        Responses: []*huma.Response{
            huma.OKResponse(&UserResponse{}),
        },
        Handler: func(ctx context.Context, input *struct{ ID string `path:"id"` }) (*UserResponse, error) {
            u, err := svc.GetByID(ctx, input.ID)
            if err != nil {
                return nil, err
            }
            return &UserResponse{ ID: u.ID, Name: u.Name, Email: u.Email }, nil
        },
    })
}
```

---

### `internal/domain/user/model.go`

```go
package user

type User struct {
    ID    string `db:"id"`
    Name  string `db:"name"`
    Email string `db:"email"`
}
```

---

### `internal/domain/user/service.go`

```go
package user

import (
    "context"
)

type Repository interface {
    Create(ctx context.Context, u *User) error
    GetByID(ctx context.Context, id string) (*User, error)
}

type Service struct {
    repo Repository
    logger interface{ Info(msg string) }
}

func NewUserService(r Repository, logger interface{ Info(msg string) }) *Service {
    return &Service{repo: r, logger: logger}
}

func (s *Service) Create(ctx context.Context, u *User) error {
    // TODO: générer ID, hash passwords si besoin, validations métiers
    return s.repo.Create(ctx, u)
}

func (s *Service) GetByID(ctx context.Context, id string) (*User, error) {
    return s.repo.GetByID(ctx, id)
}
```

---

### `internal/repository/sql/user_repo.go`

```go
package sql

import (
    "context"
    "github.com/jmoiron/sqlx"
    "github.com/you/myapp/internal/domain/user"
)

type UserRepo struct {
    db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo { return &UserRepo{db: db} }

func (r *UserRepo) Create(ctx context.Context, u *user.User) error {
    // simple insert, adapte selon ton schéma
    _, err := r.db.NamedExecContext(ctx, `INSERT INTO users (id, name, email) VALUES (:id, :name, :email)`, u)
    return err
}

func (r *UserRepo) GetByID(ctx context.Context, id string) (*user.User, error) {
    var u user.User
    if err := r.db.GetContext(ctx, &u, `SELECT id, name, email FROM users WHERE id = $1`, id); err != nil {
        return nil, err
    }
    return &u, nil
}
```

---

### `internal/config/config.go`

```go
package config

import (
    "os"
)

type Config struct {
    DatabaseURL string
    ServerURL   string
}

func Load() (*Config, error) {
    return &Config{
        DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@db:5432/postgres?sslmode=disable"),
        ServerURL:   getEnv("SERVER_URL", "http://localhost:8080"),
    }, nil
}

func getEnv(k, def string) string {
    if v := os.Getenv(k); v != "" { return v }
    return def
}
```

---

### `internal/utils/logger.go`

```go
package utils

import "log"

type Logger struct{}

func NewLogger(_ interface{}) *Logger { return &Logger{} }
func (l *Logger) Info(msg string) { log.Println("INFO:", msg) }
```

---

## Docker & docker-compose

### `Dockerfile`

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 go build -o /bin/myapp ./cmd/myapp

FROM alpine:3.18
COPY --from=builder /bin/myapp /bin/myapp
EXPOSE 8080
ENTRYPOINT ["/bin/myapp"]
```

---

### `docker-compose.yml`

```yaml
version: '3.8'
services:
  db:
    image: postgres:15
    environment:
      POSTGRES_PASSWORD: postgres
    ports:
      - "5432:5432"
    volumes:
      - db_data:/var/lib/postgresql/data

  api:
    build: .
    depends_on:
      - db
    environment:
      DATABASE_URL: postgres://postgres:postgres@db:5432/postgres?sslmode=disable
    ports:
      - "8080:8080"

volumes:
  db_data:
```

---

## Makefile

```makefile
.PHONY: run build test lint format

run:
	go run ./cmd/myapp

build:
	go build -o bin/myapp ./cmd/myapp

docker-build:
	docker build -t myapp:local .

up:
	docker-compose up --build

test:
	go test ./... -v

lint:
	golangci-lint run || true
```

---

## CI: `.github/workflows/ci.yml`

```yaml
name: CI
on: [push, pull_request]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - name: Build
        run: go build ./...
      - name: Test
        run: go test ./... -v
```

---

## Exemple de test unitaire basique

`internal/domain/user/service_test.go`

```go
package user_test

import (
    "context"
    "testing"

    "github.com/stretchr/testify/require"
    "github.com/you/myapp/internal/domain/user"
)

type mockRepo struct{}
func (m *mockRepo) Create(ctx context.Context, u *user.User) error { u.ID = "123"; return nil }
func (m *mockRepo) GetByID(ctx context.Context, id string) (*user.User, error) { return &user.User{ID: id, Name: "Bob", Email: "b@x"}, nil }

func TestCreate(t *testing.T) {
    svc := user.NewUserService(&mockRepo{}, nil)
    u := &user.User{Name: "Alice", Email: "a@x"}
    err := svc.Create(context.Background(), u)
    require.NoError(t, err)
    require.NotEmpty(t, u.ID)
}
```

---

## `.env.example`

```
DATABASE_URL=postgres://postgres:postgres@db:5432/postgres?sslmode=disable
SERVER_URL=http://localhost:8080
```

---

## Migrations

Met en place `migrations/` et utilise `golang-migrate/migrate` ou `pressly/goose`. Exemple SQL minimal `migrations/001_create_users.sql`:

```sql
CREATE TABLE IF NOT EXISTS users (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  email TEXT UNIQUE NOT NULL
);
```

---

## README minimal (README.md)

Explique comment lancer :

```md
# MyApp

## Run locally

cp .env.example .env
make run

## Run with Docker

docker-compose up --build

## Tests

make test
```

---

## Étapes suivantes & personnalisations suggérées

1. Ajouter gestion d'auth (JWT), middlewares (logging, recovery) et monitoring.
2. Compléter les validations métier et les erreurs personnalisées.
3. Ajouter des tests d'intégration avec une base Postgres en CI (service Docker).
4. Ajouter `go:generate` + `mockgen` pour générer des mocks pour les dépôts.

---

Si tu veux, je peux :

* Générer directement les fichiers complets (je peux injecter chaque fichier dans le canvas séparément si tu veux)
* Adapter le squelette pour SQLite (local dev), ou pour utiliser un autre router
* Ajouter un exemple complet d’authentification JWT et les tests d’intégration

Dis‑moi quelle option tu veux que je fasse maintenant et je te fournis les fichiers prêts à copier.


Que veux-tu que je fasse maintenant ? Choisis une option et j’exécute tout de suite :

1. Générer les fichiers sur mesure pour **SQLite** (dev local) au lieu de Postgres.
2. Ajouter **authentification JWT** complète (routes d’auth, middleware, stockage tokens).
3. Ajouter **tests d’intégration** (avec Postgres dockerisé en CI).
4. Te fournir les fichiers prêts à copier **un par un** (je peux injecter chaque fichier dans le canvas pour que tu les télécharges).
5. Autre (dis-moi directement ce que tu veux).

Fait — j’ai ajouté un document **“Myapp-files-ready-to-copy”** contenant tous les fichiers principaux prêts à copier-coller (main.go, app.go, handlers, modèles, Dockerfile, Makefile, CI, tests, migrations, etc.).

Tu peux maintenant :

* Copier-coller ces fichiers dans ton repo et lancer `make run` ou `docker-compose up --build`.
* Me dire quelle amélioration tu veux en **priorité** : SQLite, JWT auth, tests d’intégration, ou que je génère chaque fichier dans un document séparé pour téléchargement.

### example complet

# Fichiers prêts à copier — Projet Go + Huma

Ci-dessous tu trouveras les fichiers principaux prêts à copier-coller dans ton repo. Chaque fichier est présenté dans une section nommée — crée les dossiers correspondants (`cmd/myapp`, `internal/...`, etc.) et colle le contenu.

---

## `cmd/myapp/main.go`

```go
package main

import (
    "log"
    "net/http"

    "github.com/you/myapp/internal/app"
)

func main() {
    a, err := app.New()
    if err != nil {
        log.Fatalf("app init: %v", err)
    }

    log.Println("Server listening on :8080")
    if err := http.ListenAndServe(":8080", a.Handler()); err != nil {
        log.Fatalf("server: %v", err)
    }
}
```

---

## `internal/app/app.go`

```go
package app

import (
    "net/http"

    "github.com/you/myapp/internal/config"
    v1 "github.com/you/myapp/internal/api/v1"
    "github.com/you/myapp/internal/repository/sql"
    "github.com/you/myapp/internal/utils"

    "github.com/danielgtaylor/huma/v2"
    "github.com/danielgtaylor/huma/v2/adapters/humachi"
    "github.com/go-chi/chi/v5"
)

type App struct {
    api huma.API
}

func New() (*App, error) {
    cfg, err := config.Load()
    if err != nil {
        return nil, err
    }

    logger := utils.NewLogger(cfg)

    db, err := sql.NewPostgres(cfg.DatabaseURL)
    if err != nil {
        return nil, err
    }

    r := chi.NewRouter()

    apiConfig := huma.DefaultConfig("MyApp API", "v1.0.0")
    api := humachi.New(r, apiConfig)

    userRepo := sql.NewUserRepo(db)
    userService := v1.NewUserService(userRepo, logger)

    v1.RegisterUserRoutes(api, userService)

    return &App{api: api}, nil
}

func (a *App) Handler() http.Handler {
    return a.api.Handler()
}
```

---

## `internal/api/v1/models.go`

```go
package v1

type UserResponse struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

type CreateUserRequest struct {
    Name  string `json:"name" validate:"required"`
    Email string `json:"email" validate:"required,email"`
}
```

---

## `internal/api/v1/handlers.go`

```go
package v1

import (
    "context"
    "net/http"

    "github.com/danielgtaylor/huma/v2"
    "github.com/you/myapp/internal/domain/user"
)

func RegisterUserRoutes(api huma.API, svc *user.Service) {
    api.Register(huma.Operation{
        OperationID: "createUser",
        Method:      http.MethodPost,
        Path:        "/users",
        Request:     &CreateUserRequest{},
        Responses: []*huma.Response{
            huma.CreatedResponse(&UserResponse{}),
        },
        Handler: func(ctx context.Context, input *CreateUserRequest) (*UserResponse, error) {
            u := &user.User{ Name: input.Name, Email: input.Email }
            if err := svc.Create(ctx, u); err != nil {
                return nil, err
            }
            return &UserResponse{ ID: u.ID, Name: u.Name, Email: u.Email }, nil
        },
    })

    api.Register(huma.Operation{
        OperationID: "getUser",
        Method:      http.MethodGet,
        Path:        "/users/{id}",
        PathParams:  []*huma.Parameter{huma.PathParameter("id", "User ID")},
        Responses: []*huma.Response{
            huma.OKResponse(&UserResponse{}),
        },
        Handler: func(ctx context.Context, input *struct{ ID string `path:"id"` }) (*UserResponse, error) {
            u, err := svc.GetByID(ctx, input.ID)
            if err != nil {
                return nil, err
            }
            return &UserResponse{ ID: u.ID, Name: u.Name, Email: u.Email }, nil
        },
    })
}
```

---

## `internal/domain/user/model.go`

```go
package user

type User struct {
    ID    string `db:"id"`
    Name  string `db:"name"`
    Email string `db:"email"`
}
```

---

## `internal/domain/user/service.go`

```go
package user

import (
    "context"
)

type Repository interface {
    Create(ctx context.Context, u *User) error
    GetByID(ctx context.Context, id string) (*User, error)
}

type Service struct {
    repo Repository
    logger interface{ Info(msg string) }
}

func NewUserService(r Repository, logger interface{ Info(msg string) }) *Service {
    return &Service{repo: r, logger: logger}
}

func (s *Service) Create(ctx context.Context, u *User) error {
    // TODO: générer ID, validations métiers
    return s.repo.Create(ctx, u)
}

func (s *Service) GetByID(ctx context.Context, id string) (*User, error) {
    return s.repo.GetByID(ctx, id)
}
```

---

## `internal/repository/sql/user_repo.go`

```go
package sql

import (
    "context"
    "github.com/jmoiron/sqlx"
    "github.com/you/myapp/internal/domain/user"
)

type UserRepo struct {
    db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo { return &UserRepo{db: db} }

func (r *UserRepo) Create(ctx context.Context, u *user.User) error {
    _, err := r.db.NamedExecContext(ctx, `INSERT INTO users (id, name, email) VALUES (:id, :name, :email)`, u)
    return err
}

func (r *UserRepo) GetByID(ctx context.Context, id string) (*user.User, error) {
    var u user.User
    if err := r.db.GetContext(ctx, &u, `SELECT id, name, email FROM users WHERE id = $1`, id); err != nil {
        return nil, err
    }
    return &u, nil
}
```

---

## `internal/config/config.go`

```go
package config

import (
    "os"
)

type Config struct {
    DatabaseURL string
    ServerURL   string
}

func Load() (*Config, error) {
    return &Config{
        DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@db:5432/postgres?sslmode=disable"),
        ServerURL:   getEnv("SERVER_URL", "http://localhost:8080"),
    }, nil
}

func getEnv(k, def string) string {
    if v := os.Getenv(k); v != "" { return v }
    return def
}
```

---

## `internal/utils/logger.go`

```go
package utils

import "log"

type Logger struct{}

func NewLogger(_ interface{}) *Logger { return &Logger{} }
func (l *Logger) Info(msg string) { log.Println("INFO:", msg) }
```

---

## `Dockerfile`

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 go build -o /bin/myapp ./cmd/myapp

FROM alpine:3.18
COPY --from=builder /bin/myapp /bin/myapp
EXPOSE 8080
ENTRYPOINT ["/bin/myapp"]
```

---

## `docker-compose.yml`

```yaml
version: '3.8'
services:
  db:
    image: postgres:15
    environment:
      POSTGRES_PASSWORD: postgres
    ports:
      - "5432:5432"
    volumes:
      - db_data:/var/lib/postgresql/data

  api:
    build: .
    depends_on:
      - db
    environment:
      DATABASE_URL: postgres://postgres:postgres@db:5432/postgres?sslmode=disable
    ports:
      - "8080:8080"

volumes:
  db_data:
```

---

## `Makefile`

```makefile
.PHONY: run build test lint format

run:
	go run ./cmd/myapp

build:
	go build -o bin/myapp ./cmd/myapp

docker-build:
	docker build -t myapp:local .

up:
	docker-compose up --build

test:
	go test ./... -v

lint:
	golangci-lint run || true
```

---

## `.github/workflows/ci.yml`

```yaml
name: CI
on: [push, pull_request]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - name: Build
        run: go build ./...
      - name: Test
        run: go test ./... -v
```

---

## `internal/domain/user/service_test.go`

```go
package user_test

import (
    "context"
    "testing"

    "github.com/stretchr/testify/require"
    "github.com/you/myapp/internal/domain/user"
)

type mockRepo struct{}
func (m *mockRepo) Create(ctx context.Context, u *user.User) error { u.ID = "123"; return nil }
func (m *mockRepo) GetByID(ctx context.Context, id string) (*user.User, error) { return &user.User{ID: id, Name: "Bob", Email: "b@x"}, nil }

func TestCreate(t *testing.T) {
    svc := user.NewUserService(&mockRepo{}, nil)
    u := &user.User{Name: "Alice", Email: "a@x"}
    err := svc.Create(context.Background(), u)
    require.NoError(t, err)
    require.NotEmpty(t, u.ID)
}
```

---

## `.env.example`

```
DATABASE_URL=postgres://postgres:postgres@db:5432/postgres?sslmode=disable
SERVER_URL=http://localhost:8080
```

---

## `migrations/001_create_users.sql`

```sql
CREATE TABLE IF NOT EXISTS users (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  email TEXT UNIQUE NOT NULL
);
```

---

Si tu veux que je :

* **convertisse ce squelette pour SQLite** (pratique pour le dev local),
* **ajoute l’auth JWT complète** (routes, middleware, stockage blacklist si nécessaire),
* **génère les fichiers dans des documents séparés** pour téléchargement individuel (chaque fichier dans son propre canvas),
* ou **ajoute des tests d’intégration** et CI Dockerized,

dis‑le et je les ajoute immédiatement dans le canvas.


Fait — j’ai ajouté un document **“Myapp-files-ready-to-copy”** contenant tous les fichiers principaux prêts à copier-coller (main.go, app.go, handlers, modèles, Dockerfile, Makefile, CI, tests, migrations, etc.).

Tu peux maintenant :

* Copier-coller ces fichiers dans ton repo et lancer `make run` ou `docker-compose up --build`.
* Me dire quelle amélioration tu veux en **priorité** : SQLite, JWT auth, tests d’intégration, ou que je génère chaque fichier dans un document séparé pour téléchargement.
