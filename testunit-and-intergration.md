Excellent sujet ! 👌
`mockgen` (du projet [GoMock](https://github.com/golang/mock)) est un outil très utilisé en Go pour **générer automatiquement des mocks** d’interfaces, afin de faciliter les **tests unitaires** sans dépendances réelles (base de données, API, etc.).

Voici un guide **complet et concret** pour bien t’en servir.

---

## ⚙️ 1. Installation

Installe `mockgen` globalement (ou via `go install` pour Go ≥ 1.17) :

```bash
go install go.uber.org/mock/mockgen@latest
```

Assure-toi que `$GOPATH/bin` est dans ton `PATH` :

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

---

## 🧩 2. Exemple d’interface à mocker

Imaginons une interface `UserRepository` dans ton projet :

```go
// file: internal/repository/user_repository.go
package repository

type User struct {
	ID   int64
	Name string
}

type UserRepository interface {
	GetByID(id int64) (*User, error)
	Create(user *User) error
}
```

---

## 🪄 3. Génération automatique du mock

Deux manières sont possibles :

### a) **Manuellement en ligne de commande :**

```bash
mockgen -source=internal/repository/user_repository.go \
        -destination=internal/mocks/mock_user_repository.go \
        -package=mocks
```

Cela crée un fichier `mock_user_repository.go` contenant une implémentation mock de `UserRepository`.

### b) **Avec `go:generate` (recommandé)**

Ajoute cette ligne en haut du fichier :

```go
//go:generate mockgen -source=user_repository.go -destination=../mocks/mock_user_repository.go -package=mocks
```

Puis exécute :

```bash
go generate ./...
```

👉 Cela permet de régénérer tous tes mocks automatiquement sans te rappeler les commandes exactes.

---

## 🧪 4. Utilisation du mock dans un test

Exemple avec le framework de test standard :

```go
// file: internal/service/user_service_test.go
package service_test

import (
	"errors"
	"testing"

	"example.com/myapp/internal/mocks"
	"example.com/myapp/internal/repository"
	"example.com/myapp/internal/service"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)

	mockRepo.EXPECT().
		GetByID(int64(1)).
		Return(&repository.User{ID: 1, Name: "Alice"}, nil)

	svc := service.NewUserService(mockRepo)

	user, err := svc.GetUser(1)
	assert.NoError(t, err)
	assert.Equal(t, "Alice", user.Name)
}
```

---

## 🧱 5. Exemple de service utilisant l’interface

```go
// file: internal/service/user_service.go
package service

import "example.com/myapp/internal/repository"

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUser(id int64) (*repository.User, error) {
	return s.repo.GetByID(id)
}
```

---

## 🧰 6. Arborescence typique

```shell
internal/
  repository/
    user_repository.go
  service/
    user_service.go
    user_service_test.go
  mocks/
    mock_user_repository.go  <-- généré
```

---

## ✅ 7. Commandes utiles

| Commande                     | Description                |
| ---------------------------- | -------------------------- |
| `go generate ./...`          | Régénère tous les mocks    |
| `go test ./...`              | Lance tous les tests       |
| `mockgen -package=mocks ...` | Génération manuelle rapide |

---

Souhaites-tu que je te montre un **exemple complet fonctionnel** (avec `go.mod`, fichiers, tests exécutables) que tu pourrais copier-coller et lancer ?
Je peux te fournir un mini-projet `UserService` prêt à tester avec `mockgen`.

Parfait 👍
Voici un **mini-projet Go complet et exécutable** montrant comment utiliser `mockgen` pour tester un service qui dépend d’un repository.
Tu pourras le **copier-coller intégralement** dans un dossier et exécuter `go test ./...`.

---

## 🧩 Arborescence du projet

```shell
mockgen-example/
├── go.mod
├── internal/
│   ├── repository/
│   │   └── user_repository.go
│   ├── service/
│   │   ├── user_service.go
│   │   └── user_service_test.go
│   └── mocks/
│       └── (généré automatiquement)
```

---

## 🧱 1. `go.mod`

```go
module mockgen-example

go 1.23

require (
	go.uber.org/mock v0.5.0
	github.com/stretchr/testify v1.10.0
)
```

---

## 📦 2. `internal/repository/user_repository.go`

```go
package repository

type User struct {
	ID   int64
	Name string
}

type UserRepository interface {
	GetByID(id int64) (*User, error)
	Create(user *User) error
}

//go:generate mockgen -source=user_repository.go -destination=../mocks/mock_user_repository.go -package=mocks
```

💡 La ligne `//go:generate` permet de générer le mock automatiquement.

---

## ⚙️ 3. Générer le mock

À la racine du projet :

```bash
go generate ./...
```

Cela va créer automatiquement le fichier suivant :

```
internal/mocks/mock_user_repository.go
```

*(si `mockgen` n’est pas installé : `go install go.uber.org/mock/mockgen@latest`)*

---

## 🧰 4. `internal/service/user_service.go`

```go
package service

import "mockgen-example/internal/repository"

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUser(id int64) (*repository.User, error) {
	return s.repo.GetByID(id)
}

func (s *UserService) CreateUser(name string) (*repository.User, error) {
	user := &repository.User{Name: name}
	err := s.repo.Create(user)
	return user, err
}
```

---

## 🧪 5. `internal/service/user_service_test.go`

```go
package service_test

import (
	"errors"
	"testing"

	"mockgen-example/internal/mocks"
	"mockgen-example/internal/repository"
	"mockgen-example/internal/service"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)

	mockRepo.EXPECT().
		GetByID(int64(1)).
		Return(&repository.User{ID: 1, Name: "Alice"}, nil)

	svc := service.NewUserService(mockRepo)

	user, err := svc.GetUser(1)
	assert.NoError(t, err)
	assert.Equal(t, "Alice", user.Name)
}

func TestGetUser_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)

	mockRepo.EXPECT().
		GetByID(int64(2)).
		Return(nil, errors.New("not found"))

	svc := service.NewUserService(mockRepo)

	user, err := svc.GetUser(2)
	assert.Error(t, err)
	assert.Nil(t, user)
}
```

---

## ▶️ 6. Exécution

Une fois le projet en place :

```bash
go mod tidy
go generate ./...
go test ./... -v
```

✅ Résultat attendu :

```shell
=== RUN   TestGetUser_Success
--- PASS: TestGetUser_Success (0.00s)
=== RUN   TestGetUser_NotFound
--- PASS: TestGetUser_NotFound (0.00s)
PASS
ok      mockgen-example/internal/service  0.002s
```

---

Souhaites-tu que je te fournisse un **script `Makefile`** (avec `make test`, `make mock`, etc.) pour automatiser la génération et les tests ?
C’est souvent très pratique dans les projets réels.

Tu veux dire **`pgx`** (le driver PostgreSQL moderne pour Go) ?
Si oui 👉 parfait, je peux te montrer comment **mockgen** s’utilise dans un projet **Go + pgx + repository pattern**, pour tester la logique sans accéder à une vraie base PostgreSQL.

---

## 🎯 Objectif

On va créer un mini-projet :

* avec un `UserRepository` qui utilise `pgxpool.Pool` pour exécuter des requêtes SQL,
* un `UserService` qui dépend de ce repository,
* et des tests unitaires avec **GoMock** (`mockgen`) sans base réelle.

---

## 🧩 Structure du projet

```shell
mockgen-pgx/
├── go.mod
├── internal/
│   ├── repository/
│   │   └── user_repository.go
│   ├── service/
│   │   ├── user_service.go
│   │   └── user_service_test.go
│   └── mocks/
│       └── (mock généré)
```

---

## ⚙️ 1. `go.mod`

```go
module mockgen-pgx

go 1.23

require (
	github.com/jackc/pgx/v5 v5.6.0
	go.uber.org/mock v0.5.0
	github.com/stretchr/testify v1.10.0
)
```

---

## 🧱 2. `internal/repository/user_repository.go`

```go
package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID   int64
	Name string
}

type DB interface {
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgxpool.Row
	Exec(ctx context.Context, sql string, args ...interface{}) (pgxpool.CommandTag, error)
}

// UserRepository encapsule l’accès à la base
type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*User, error) {
	row := r.db.QueryRow(ctx, "SELECT id, name FROM users WHERE id=$1", id)

	var user User
	if err := row.Scan(&user.ID, &user.Name); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Create(ctx context.Context, name string) (*User, error) {
	row := r.db.QueryRow(ctx,
		"INSERT INTO users(name) VALUES($1) RETURNING id, name",
		name,
	)

	var user User
	if err := row.Scan(&user.ID, &user.Name); err != nil {
		return nil, err
	}
	return &user, nil
}
```

💡 On pourrait mocker `pgxpool.Pool`, mais c’est un struct concret (pas une interface).
➡️ Donc, pour les tests, on **abstrait** la dépendance via une **interface `UserStore`**.

---

## 🧰 3. Interface à mocker : `internal/repository/store.go`

```go
package repository

import (
	"context"
)

type UserStore interface {
	GetByID(ctx context.Context, id int64) (*User, error)
	Create(ctx context.Context, name string) (*User, error)
}

//go:generate mockgen -source=store.go -destination=../mocks/mock_user_store.go -package=mocks
```

---

## 🪄 4. Générer le mock

```bash
go install go.uber.org/mock/mockgen@latest
go generate ./...
```

Cela crée `internal/mocks/mock_user_store.go`.

---

## 🧩 5. `internal/service/user_service.go`

```go
package service

import (
	"context"
	"mockgen-pgx/internal/repository"
)

type UserService struct {
	store repository.UserStore
}

func NewUserService(store repository.UserStore) *UserService {
	return &UserService{store: store}
}

func (s *UserService) RegisterUser(ctx context.Context, name string) (*repository.User, error) {
	return s.store.Create(ctx, name)
}

func (s *UserService) GetProfile(ctx context.Context, id int64) (*repository.User, error) {
	return s.store.GetByID(ctx, id)
}
```

---

## 🧪 6. `internal/service/user_service_test.go`

```go
package service_test

import (
	"context"
	"errors"
	"testing"

	"mockgen-pgx/internal/mocks"
	"mockgen-pgx/internal/repository"
	"mockgen-pgx/internal/service"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestRegisterUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockUserStore(ctrl)

	mockStore.EXPECT().
		Create(gomock.Any(), "Alice").
		Return(&repository.User{ID: 1, Name: "Alice"}, nil)

	svc := service.NewUserService(mockStore)

	user, err := svc.RegisterUser(context.Background(), "Alice")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), user.ID)
	assert.Equal(t, "Alice", user.Name)
}

func TestGetProfile_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockUserStore(ctrl)

	mockStore.EXPECT().
		GetByID(gomock.Any(), int64(99)).
		Return(nil, errors.New("not found"))

	svc := service.NewUserService(mockStore)

	user, err := svc.GetProfile(context.Background(), 99)
	assert.Error(t, err)
	assert.Nil(t, user)
}
```

---

## ▶️ 7. Exécution

```bash
go mod tidy
go generate ./...
go test ./... -v
```

✅ Sortie attendue :

```shell
=== RUN   TestRegisterUser_Success
--- PASS: TestRegisterUser_Success (0.00s)
=== RUN   TestGetProfile_NotFound
--- PASS: TestGetProfile_NotFound (0.00s)
PASS
ok      mockgen-pgx/internal/service  0.002s
```

---

## 💡 Notes avancées

* Dans les vrais projets, `pgxpool.Pool` reste au niveau de `main` ou `repository`, jamais dans les tests.
* `mockgen` permet de **générer des mocks pour n’importe quelle interface**, donc tu peux aussi mocker une interface de type `DB` si tu veux tester le repository lui-même.
* Si tu veux un repository générique, on peut combiner **Squirrel** (pour builder les requêtes) et `sqlx`/`pgx` avec le même principe de mock.

---

Souhaites-tu que je te montre la **version complète avec `pgx + Squirrel`** et un **repository générique mockable** (comme dans un vrai projet de type REST API avec Gin) ?
C’est souvent la meilleure approche pour des tests unitaires propres.

Oui, très bon choix — intégrer **Huma CLI** dans un projet Go + `pgx` + tests avec `mockgen` est tout à fait faisable. Je vais te proposer un exemple **bout à bout** : structure de projet, `main.go` intégrant `humacli`, un repository utilisant `pgx`, une interface mockée, un service, un handler Huma, et un test unitaire avec mock.

Tu pourras t’en inspirer et l’adapter à ton contexte.

---

## 🏗️ Structure projet proposée

```shell
myapp/
├── go.mod
├── main.go
├── internal/
│   ├── repository/
│   │   ├── store.go
│   │   └── pgx_user_repository.go
│   ├── service/
│   │   └── user_service.go
│   ├── handler/
│   │   └── user_handler.go
│   └── mocks/
│       └── mock_store.go  ← généré par mockgen
```

On utilisera `humacli` (le CLI fourni par Huma) pour lancer le serveur, potentiellement ajouter des commandes comme `openapi`, etc. Huma a un module `humacli` pour cela. ([huma.rocks][1])

---

## 📦 Dépendances (go.mod)

```go
module myapp

go 1.23

require (
    github.com/jackc/pgx/v5 v5.6.0
    github.com/danielgtaylor/huma/v2 v2.34.1  // version possible
    github.com/danielgtaylor/huma/v2/humacli v2.34.1
    go.uber.org/mock v0.5.0
    github.com/stretchr/testify v1.10.0
)
```

---

## 📂 internal/repository/store.go

Ici on définit une **interface abstraite** que l’on pourra mocker :

```go
package repository

import (
    "context"
)

type User struct {
    ID   int64
    Name string
}

type Store interface {
    GetByID(ctx context.Context, id int64) (*User, error)
    Create(ctx context.Context, name string) (*User, error)
}

// optionnel : pour générer le mock automatiquement
//go:generate mockgen -source=store.go -destination=../mocks/mock_store.go -package=mocks
```

---

## 📂 internal/repository/pgx_user_repository.go

Implémentation concrète de `Store` en utilisant `pgxpool.Pool` :

```go
package repository

import (
    "context"
    "fmt"

    "github.com/jackc/pgx/v5/pgxpool"
)

type PgxUserRepository struct {
    pool *pgxpool.Pool
}

func NewPgxUserRepository(pool *pgxpool.Pool) *PgxUserRepository {
    return &PgxUserRepository{pool: pool}
}

func (r *PgxUserRepository) GetByID(ctx context.Context, id int64) (*User, error) {
    row := r.pool.QueryRow(ctx, "SELECT id, name FROM users WHERE id=$1", id)
    var u User
    if err := row.Scan(&u.ID, &u.Name); err != nil {
        return nil, err
    }
    return &u, nil
}

func (r *PgxUserRepository) Create(ctx context.Context, name string) (*User, error) {
    row := r.pool.QueryRow(ctx,
        "INSERT INTO users(name) VALUES($1) RETURNING id, name",
        name,
    )
    var u User
    if err := row.Scan(&u.ID, &u.Name); err != nil {
        return nil, err
    }
    return &u, nil
}
```

---

## 📂 internal/service/user_service.go

Service qui fait la logique métier :

```go
package service

import (
    "context"
    "myapp/internal/repository"
)

type UserService struct {
    store repository.Store
}

func NewUserService(store repository.Store) *UserService {
    return &UserService{store: store}
}

func (s *UserService) RegisterUser(ctx context.Context, name string) (*repository.User, error) {
    return s.store.Create(ctx, name)
}

func (s *UserService) GetUser(ctx context.Context, id int64) (*repository.User, error) {
    return s.store.GetByID(ctx, id)
}
```

---

## 📂 internal/handler/user_handler.go

Handler HTTP pour Huma :

```go
package handler

import (
    "context"
    "net/http"
    "strconv"

    "myapp/internal/service"

    "github.com/danielgtaylor/huma/v2"
)

type CreateUserInput struct {
    // Huma permet de déclarer que ce champ est dans le body
    Body struct {
        Name string `json:"name"`
    }
}

type CreateUserOutput struct {
    Body struct {
        ID   int64  `json:"id"`
        Name string `json:"name"`
    }
}

type GetUserInput struct {
    // id en path
    Path struct {
        ID string `path:"id"`
    }
}

type GetUserOutput struct {
    Body struct {
        ID   int64  `json:"id"`
        Name string `json:"name"`
    }
}

func RegisterUserRoutes(api huma.API, svc *service.UserService) {
    huma.Post[CreateUserInput, CreateUserOutput](api, "/users", func(ctx context.Context, in *CreateUserInput) (*CreateUserOutput, error) {
        out := &CreateUserOutput{}
        u, err := svc.RegisterUser(ctx, in.Body.Name)
        if err != nil {
            return nil, err
        }
        out.Body.ID = u.ID
        out.Body.Name = u.Name
        return out, nil
    })

    huma.Get[GetUserInput, GetUserOutput](api, "/users/{id}", func(ctx context.Context, in *GetUserInput) (*GetUserOutput, error) {
        out := &GetUserOutput{}
        // convert path string to int64
        id64, err := strconv.ParseInt(in.Path.ID, 10, 64)
        if err != nil {
            return nil, err
        }
        u, err := svc.GetUser(ctx, id64)
        if err != nil {
            return nil, err
        }
        out.Body.ID = u.ID
        out.Body.Name = u.Name
        return out, nil
    })
}
```

---

## 🧰 main.go

C’est ici que tu utilises **humacli**, configures la CLI + serveur HTTP :

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "time"

    "myapp/internal/handler"
    "myapp/internal/repository"
    "myapp/internal/service"

    "github.com/danielgtaylor/huma/v2"
    "github.com/danielgtaylor/huma/v2/humacli"
    "github.com/danielgtaylor/huma/v2/adapters/humachi"
    "github.com/jackc/pgx/v5/pgxpool"
)

type Options struct {
    Port int `doc:"HTTP listen port" short:"p" default:"8888"`
    DSN  string `doc:"PostgreSQL DSN (connection string)"`
}

func main() {
    cli := humacli.New(func(hooks humacli.Hooks, opts *Options) {
        // Set up router + API
        router := humachi.New(nil, huma.DefaultConfig("MyApp", "1.0.0"))
        api := router  // router est aussi un API

        // Connexion à la base
        pool, err := pgxpool.New(context.Background(), opts.DSN)
        if err != nil {
            log.Fatalf("failed to connect to DB: %v", err)
        }

        // Crée repository, service, et enregistre handlers
        store := repository.NewPgxUserRepository(pool)
        svc := service.NewUserService(store)
        handler.RegisterUserRoutes(api, svc)

        // Démarrer le serveur HTTP quand la CLI démarre
        hooks.OnStart(func() {
            addr := fmt.Sprintf(":%d", opts.Port)
            log.Printf("Listening on %s", addr)
            go func() {
                http.ListenAndServe(addr, router)  // router implémente http.Handler via Huma
            }()
        })
        // shutdown hook si besoin
        hooks.OnStop(func() {
            ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
            defer cancel()
            pool.Close()  // fermer le pool
            // ici, tu pourrais shutdown le serveur HTTP proprement
        })
    })

    cli.Run()
}
```

Avec ça, tu peux lancer :

```bash
go run main.go --dsn="postgres://user:pass@host:port/dbname" --port=8888
```

et tu auras tes endpoints `/users` (POST) et `/users/{id}` (GET), avec documentation OpenAPI automatiquement exposée par Huma (ex. `/openapi.json`). ([Go Packages][2])

---

## 🔄 Génération du mock

Dans `internal/repository/store.go`, on a le commentaire `//go:generate`.
Ensuite, dans la racine :

```bash
go install go.uber.org/mock/mockgen@latest
go generate ./internal/repository
```

Cela crée `internal/mocks/mock_store.go`.

---

## 🧪 Test unitaire du service avec mock

Par exemple, tester `UserService` :

```go
package service_test

import (
    "context"
    "errors"
    "testing"

    "myapp/internal/mocks"
    "myapp/internal/repository"
    "myapp/internal/service"

    "github.com/golang/mock/gomock"
    "github.com/stretchr/testify/assert"
)

func TestRegisterUser_Success(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockStore := mocks.NewMockStore(ctrl)

    mockStore.EXPECT().
        Create(gomock.Any(), "Alice").
        Return(&repository.User{ID: 1, Name: "Alice"}, nil)

    svc := service.NewUserService(mockStore)

    u, err := svc.RegisterUser(context.Background(), "Alice")
    assert.NoError(t, err)
    assert.Equal(t, int64(1), u.ID)
    assert.Equal(t, "Alice", u.Name)
}

func TestGetUser_NotFound(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockStore := mocks.NewMockStore(ctrl)

    mockStore.EXPECT().
        GetByID(gomock.Any(), int64(42)).
        Return(nil, errors.New("not found"))

    svc := service.NewUserService(mockStore)

    u, err := svc.GetUser(context.Background(), 42)
    assert.Error(t, err)
    assert.Nil(t, u)
}
```

Si tu veux, tu peux aussi **tester le handler Huma** à un niveau un peu plus haut : injecter un `UserService` mocké dans le handler, puis faire une requête HTTP à ton router Huma (via `httptest.NewRequest`) pour vérifier la réponse JSON.

---

## ✅ Résumé & conseils

* L’essentiel est de **séparer les couches** : interface (Store), impl. concrète (pgx), service, handler/HTTP.
* Avec cette séparation, tu peux mocker la couche `Store` (via mockgen) et tester le service ou le handler sans base de données réelle.
* `humacli` te donne un wrapper CLI / serveur prêt à l’utilisation, avec hooks `OnStart`, `OnStop`, etc. pour lancer ton serveur via la CLI. ([huma.rocks][1])
* Huma gère automatiquement l’OpenAPI, la validation, etc., ce qui rend ton handler plus concis. ([Go Packages][2])

Si tu veux, je peux te donner le code complet dans un dépôt GitHub (ou fichier zip) pour que tu puisses cloner et tester directement. Veux-tu cela ?

[1]: https://huma.rocks/features/cli/?utm_source=chatgpt.com "Service CLI - Huma"
[2]: https://pkg.go.dev/github.com/danielgtaylor/huma/v2?utm_source=chatgpt.com "huma package - github.com/danielgtaylor/huma/v2"

💯 Excellente question — **oui, absolument !**
Tu peux **tester les routes Huma CLI** très facilement, car sous le capot, Huma repose sur un **routeur HTTP standard** (`humachi` ou `humaserve`) compatible avec le package `net/http/httptest` de la lib standard Go.

👉 Donc, tu peux :

* lancer ton API Huma en mémoire (sans `main()` ni serveur réel),
* injecter un **mock de service** (via `mockgen`),
* exécuter des requêtes HTTP (`httptest.NewRequest`),
* et vérifier les **réponses JSON** et **statuts HTTP**.

Je vais te montrer un **exemple complet** avec :

* un handler Huma (`/users` POST + `/users/{id}` GET),
* un mock de `UserService`,
* un test qui vérifie les routes **sans base ni réseau**.

---

## 🧩 Structure minimale

```shell
myapp/
├── go.mod
├── internal/
│   ├── handler/
│   │   └── user_handler.go
│   ├── service/
│   │   ├── user_service.go
│   │   └── user_service_test.go
│   └── mocks/
│       └── mock_service.go  ← généré par mockgen
```

---

## 📦 go.mod

```go
module myapp

go 1.23

require (
    github.com/danielgtaylor/huma/v2 v2.34.1
    github.com/danielgtaylor/huma/v2/adapters/humachi v2.34.1
    github.com/golang/mock v0.5.0
    github.com/stretchr/testify v1.10.0
)
```

---

## ⚙️ 1. Interface de service

```go
// internal/service/user_service.go
package service

import "context"

type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type UserService interface {
	RegisterUser(ctx context.Context, name string) (*User, error)
	GetUser(ctx context.Context, id int64) (*User, error)
}

//go:generate mockgen -source=user_service.go -destination=../mocks/mock_service.go -package=mocks
```

---

## 🧱 2. Handler Huma

```go
// internal/handler/user_handler.go
package handler

import (
	"context"
	"net/http"
	"strconv"

	"myapp/internal/service"

	"github.com/danielgtaylor/huma/v2"
)

type CreateUserInput struct {
	Body struct {
		Name string `json:"name"`
	}
}

type CreateUserOutput struct {
	Body struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	}
}

type GetUserInput struct {
	Path struct {
		ID string `path:"id"`
	}
}

type GetUserOutput struct {
	Body service.User
}

func RegisterUserRoutes(api huma.API, svc service.UserService) {
	huma.Post[CreateUserInput, CreateUserOutput](api, "/users", func(ctx context.Context, in *CreateUserInput) (*CreateUserOutput, error) {
		u, err := svc.RegisterUser(ctx, in.Body.Name)
		if err != nil {
			return nil, huma.Error500InternalServerError(err)
		}
		out := &CreateUserOutput{}
		out.Body.ID = u.ID
		out.Body.Name = u.Name
		return out, nil
	})

	huma.Get[GetUserInput, GetUserOutput](api, "/users/{id}", func(ctx context.Context, in *GetUserInput) (*GetUserOutput, error) {
		id, err := strconv.ParseInt(in.Path.ID, 10, 64)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid ID")
		}
		u, err := svc.GetUser(ctx, id)
		if err != nil {
			return nil, huma.Error404NotFound("user not found")
		}
		return &GetUserOutput{Body: *u}, nil
	})
}
```

---

## 🧪 3. Test des routes Huma (sans serveur réel)

```go
// internal/handler/user_handler_test.go
package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"myapp/internal/handler"
	"myapp/internal/mocks"
	"myapp/internal/service"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateUserRoute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockUserService(ctrl)

	mockSvc.EXPECT().
		RegisterUser(gomock.Any(), "Alice").
		Return(&service.User{ID: 1, Name: "Alice"}, nil)

	// Crée un routeur Huma en mémoire
	router := humachi.New(nil, huma.DefaultConfig("Test API", "1.0"))
	handler.RegisterUserRoutes(router, mockSvc)

	body := []byte(`{"name":"Alice"}`)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "Alice", resp["name"])
}

func TestGetUserRoute_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockUserService(ctrl)
	mockSvc.EXPECT().
		GetUser(gomock.Any(), int64(99)).
		Return(nil, assert.AnError)

	router := humachi.New(nil, huma.DefaultConfig("Test API", "1.0"))
	handler.RegisterUserRoutes(router, mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/users/99", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
```

---

## ⚙️ 4. Génération du mock

```bash
go install go.uber.org/mock/mockgen@latest
go generate ./...
```

Cela va créer `internal/mocks/mock_service.go`.

---

## ▶️ 5. Exécution

```bash
go test ./... -v
```

✅ Résultat attendu :

```shell
=== RUN   TestCreateUserRoute
--- PASS: TestCreateUserRoute (0.00s)
=== RUN   TestGetUserRoute_NotFound
--- PASS: TestGetUserRoute_NotFound (0.00s)
PASS
```

---

## 🚀 En résumé

| Élément            | Rôle                                    |
| ------------------ | --------------------------------------- |
| `huma` + `humachi` | Fournit l’API HTTP testable             |
| `httptest`         | Simule des requêtes HTTP sans serveur   |
| `mockgen`          | Fournit un mock de ton service métier   |
| `testify/assert`   | Vérifie les réponses JSON et codes HTTP |

---

Souhaites-tu que je t’ajoute aussi un **test d’intégration complet** (Huma + humacli + PostgreSQL en mémoire ou Docker) pour simuler une vraie requête `go run main.go` + test HTTP ?
C’est la suite naturelle de ce setup (utile en CI/CD).

Parfait 🚀
On va maintenant monter un **exemple complet de test d’intégration** pour un projet **Go + Huma CLI + pgx + mockgen**, où :

* ton API tourne réellement (via `humacli` + HTTP),
* la base PostgreSQL tourne dans **Docker (testcontainer)**,
* et les tests font de vraies requêtes HTTP (via `net/http`) contre le serveur,
* tout en gardant tes **tests unitaires isolés avec mockgen**.

---

## 🧩 Architecture complète

```shell
huma-pgx-example/
├── go.mod
├── main.go
├── internal/
│   ├── handler/
│   │   ├── user_handler.go
│   │   └── user_handler_test.go   ← tests unitaires HTTP avec mock
│   ├── service/
│   │   ├── user_service.go
│   │   └── user_service_test.go   ← tests unitaires purs
│   ├── repository/
│   │   ├── store.go
│   │   └── pgx_user_repository.go
│   └── mocks/
│       └── mock_store.go
├── integration/
│   └── integration_test.go        ← test d’intégration Huma + DB
├── Dockerfile                     ← (facultatif pour CI)
└── Makefile
```

---

## ⚙️ `go.mod`

```go
module huma-pgx-example

go 1.23

require (
    github.com/danielgtaylor/huma/v2 v2.34.1
    github.com/danielgtaylor/huma/v2/adapters/humachi v2.34.1
    github.com/danielgtaylor/huma/v2/humacli v2.34.1
    github.com/jackc/pgx/v5 v5.6.0
    github.com/stretchr/testify v1.10.0
    go.uber.org/mock v0.5.0
    github.com/testcontainers/testcontainers-go v0.31.0
)
```

---

## 🧱 1. `internal/repository/store.go`

```go
package repository

import "context"

type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Store interface {
	GetByID(ctx context.Context, id int64) (*User, error)
	Create(ctx context.Context, name string) (*User, error)
}

//go:generate mockgen -source=store.go -destination=../mocks/mock_store.go -package=mocks
```

---

## 🧱 2. `internal/repository/pgx_user_repository.go`

```go
package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxUserRepository struct {
	pool *pgxpool.Pool
}

func NewPgxUserRepository(pool *pgxpool.Pool) *PgxUserRepository {
	return &PgxUserRepository{pool: pool}
}

func (r *PgxUserRepository) GetByID(ctx context.Context, id int64) (*User, error) {
	row := r.pool.QueryRow(ctx, "SELECT id, name FROM users WHERE id=$1", id)
	var u User
	if err := row.Scan(&u.ID, &u.Name); err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *PgxUserRepository) Create(ctx context.Context, name string) (*User, error) {
	row := r.pool.QueryRow(ctx,
		"INSERT INTO users(name) VALUES($1) RETURNING id, name",
		name,
	)
	var u User
	if err := row.Scan(&u.ID, &u.Name); err != nil {
		return nil, err
	}
	return &u, nil
}
```

---

## 🧱 3. `internal/service/user_service.go`

```go
package service

import (
	"context"
	"huma-pgx-example/internal/repository"
)

type UserService struct {
	store repository.Store
}

func NewUserService(store repository.Store) *UserService {
	return &UserService{store: store}
}

func (s *UserService) RegisterUser(ctx context.Context, name string) (*repository.User, error) {
	return s.store.Create(ctx, name)
}

func (s *UserService) GetUser(ctx context.Context, id int64) (*repository.User, error) {
	return s.store.GetByID(ctx, id)
}
```

---

## 🧱 4. `internal/handler/user_handler.go`

```go
package handler

import (
	"context"
	"net/http"
	"strconv"

	"huma-pgx-example/internal/service"

	"github.com/danielgtaylor/huma/v2"
)

type CreateUserInput struct {
	Body struct {
		Name string `json:"name"`
	}
}

type CreateUserOutput struct {
	Body struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	}
}

type GetUserInput struct {
	Path struct {
		ID string `path:"id"`
	}
}

type GetUserOutput struct {
	Body service.User
}

func RegisterUserRoutes(api huma.API, svc *service.UserService) {
	huma.Post[CreateUserInput, CreateUserOutput](api, "/users", func(ctx context.Context, in *CreateUserInput) (*CreateUserOutput, error) {
		u, err := svc.RegisterUser(ctx, in.Body.Name)
		if err != nil {
			return nil, huma.Error500InternalServerError(err)
		}
		out := &CreateUserOutput{}
		out.Body.ID = u.ID
		out.Body.Name = u.Name
		return out, nil
	})

	huma.Get[GetUserInput, GetUserOutput](api, "/users/{id}", func(ctx context.Context, in *GetUserInput) (*GetUserOutput, error) {
		id, err := strconv.ParseInt(in.Path.ID, 10, 64)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid ID")
		}
		u, err := svc.GetUser(ctx, id)
		if err != nil {
			return nil, huma.Error404NotFound("user not found")
		}
		return &GetUserOutput{Body: *u}, nil
	})
}
```

---

## 🧱 5. `main.go` (Huma CLI)

```go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"huma-pgx-example/internal/handler"
	"huma-pgx-example/internal/repository"
	"huma-pgx-example/internal/service"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Options struct {
	Port int    `doc:"HTTP port" default:"8080"`
	DSN  string `doc:"PostgreSQL DSN"`
}

func main() {
	cli := humacli.New(func(hooks humacli.Hooks, opts *Options) {
		router := humachi.New(nil, huma.DefaultConfig("HumaApp", "1.0"))
		api := router

		pool, err := pgxpool.New(context.Background(), opts.DSN)
		if err != nil {
			log.Fatalf("DB error: %v", err)
		}

		store := repository.NewPgxUserRepository(pool)
		svc := service.NewUserService(store)
		handler.RegisterUserRoutes(api, svc)

		hooks.OnStart(func() {
			addr := fmt.Sprintf(":%d", opts.Port)
			log.Printf("Listening on %s", addr)
			go http.ListenAndServe(addr, router)
		})

		hooks.OnStop(func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			pool.Close()
		})
	})
	cli.Run()
}
```

---

## 🧪 6. `integration/integration_test.go`

Voici le **test d’intégration complet** utilisant **testcontainers-go** (PostgreSQL réel dans Docker) :

```go
package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"huma-pgx-example/internal/handler"
	"huma-pgx-example/internal/repository"
	"huma-pgx-example/internal/service"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestIntegration_HumaAndPostgres(t *testing.T) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:15-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_PASSWORD": "test",
			"POSTGRES_USER":     "test",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}
	pgContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	assert.NoError(t, err)
	defer pgContainer.Terminate(ctx)

	host, _ := pgContainer.Host(ctx)
	port, _ := pgContainer.MappedPort(ctx, "5432")

	dsn := fmt.Sprintf("postgres://test:test@%s:%s/testdb?sslmode=disable", host, port.Port())
	time.Sleep(2 * time.Second) // laisse le temps à Postgres de démarrer

	pool, err := pgxpool.New(ctx, dsn)
	assert.NoError(t, err)

	_, err = pool.Exec(ctx, `
		CREATE TABLE users (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL
		);
	`)
	assert.NoError(t, err)

	router := humachi.New(nil, huma.DefaultConfig("IntegrationTest", "1.0"))
	store := repository.NewPgxUserRepository(pool)
	svc := service.NewUserService(store)
	handler.RegisterUserRoutes(router, svc)

	server := &http.Server{
		Addr:    ":8090",
		Handler: router,
	}
	go server.ListenAndServe()
	defer server.Shutdown(ctx)

	// 🔹 Teste POST /users
	body := []byte(`{"name":"Alice"}`)
	resp, err := http.Post("http://localhost:8090/users", "application/json", bytes.NewReader(body))
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	data, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(data, &result)
	assert.Equal(t, "Alice", result["name"])

	// 🔹 Teste GET /users/1
	resp2, err := http.Get("http://localhost:8090/users/1")
	assert.NoError(t, err)
	defer resp2.Body.Close()
	assert.Equal(t, http.StatusOK, resp2.StatusCode)

	var user map[string]interface{}
	json.NewDecoder(resp2.Body).Decode(&user)
	assert.Equal(t, "Alice", user["name"])
}
```

---

## 🧰 7. Makefile

```makefile
mock:
	go generate ./...

test:
	go test ./... -v

integration:
	go test ./integration -v
```

---

## ✅ Pour exécuter

```bash
go mod tidy
make mock
make integration
```

Résultat attendu :

```shell
=== RUN   TestIntegration_HumaAndPostgres
--- PASS: TestIntegration_HumaAndPostgres (4.21s)
PASS
ok  	huma-pgx-example/integration	4.227s
```

---

## 🎯 Résumé

| Type de test      | Description                                                                   |
| ----------------- | ----------------------------------------------------------------------------- |
| **Unitaires**     | Utilisent `mockgen` pour simuler le repository ou le service                  |
| **Handlers Huma** | Testent les routes avec `httptest.NewRequest` sans serveur réel               |
| **Intégration**   | Lancement d’un vrai serveur Huma + PostgreSQL Docker avec `testcontainers-go` |

---

Souhaites-tu que je te prépare la version **Dockerisée complète (Dockerfile + docker-compose + CI GitHub Actions)** de ce projet Huma + pgx + tests ?
C’est parfait pour une pipeline de déploiement ou de validation CI/CD.

Excellent 👌 — tu veux donc une **version GitLab CI/CD complète** pour ton projet Huma + Go + PostgreSQL + tests (unitaires + intégration Docker).
Je vais te fournir **toute la structure prête à l’emploi**, avec :

* un `Dockerfile` léger pour builder et exécuter ton app Huma,
* un `.gitlab-ci.yml` propre pour :

  * builder l’image,
  * exécuter les tests unitaires (mockgen inclus),
  * exécuter les tests d’intégration avec PostgreSQL (via Testcontainers),
* et une **optimisation du cache Go** pour GitLab.

---

## 🧩 Arborescence finale

```shell
huma-pgx-example/
├── .gitlab-ci.yml
├── Dockerfile
├── go.mod
├── go.sum
├── main.go
├── Makefile
├── internal/
│   ├── handler/
│   ├── repository/
│   ├── service/
│   └── mocks/
└── integration/
    └── integration_test.go
```

---

## 🐋 `Dockerfile`

Ce Dockerfile sert à **builder** et **exécuter** ton app (ou à builder pour les tests).

```dockerfile
# Étape 1 : Build
FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o huma-app main.go

# Étape 2 : Runtime minimal
FROM gcr.io/distroless/base-debian12

WORKDIR /app
COPY --from=builder /app/huma-app .

EXPOSE 8080
ENTRYPOINT ["/app/huma-app"]
```

---

## 🧪 `Makefile`

```makefile
mock:
	go generate ./...

test:
	go test ./... -v -coverprofile=coverage.out

integration:
	go test ./integration -v -coverprofile=integration.out
```

---

## 🧰 `.gitlab-ci.yml`

Voici le pipeline complet GitLab CI :

```yaml
stages:
  - lint
  - test
  - integration
  - build

variables:
  GO_VERSION: "1.23"
  GOCACHE: "$CI_PROJECT_DIR/.cache/go-build"
  GOPATH: "$CI_PROJECT_DIR/.cache/go"
  CGO_ENABLED: "1"
  DOCKER_DRIVER: overlay2

cache:
  key: ${CI_COMMIT_REF_SLUG}
  paths:
    - .cache/go-build
    - .cache/go
    - go.sum
    - go.mod

before_script:
  - mkdir -p .cache/go-build .cache/go
  - go version
  - go mod download

lint:
  stage: lint
  image: golang:${GO_VERSION}
  script:
    - go fmt ./...
    - go vet ./...
    - echo "✅ Lint OK"

unit_test:
  stage: test
  image: golang:${GO_VERSION}
  services:
    - name: docker:dind
      command: ["--tls=false"]
  variables:
    DOCKER_HOST: tcp://docker:2375/
    DOCKER_TLS_CERTDIR: ""
  script:
    - go install go.uber.org/mock/mockgen@latest
    - make mock
    - make test
  artifacts:
    reports:
      junit: report.xml
    paths:
      - coverage.out
    expire_in: 1 week

integration_test:
  stage: integration
  image: golang:${GO_VERSION}
  services:
    - name: docker:dind
      command: ["--tls=false"]
  variables:
    DOCKER_HOST: tcp://docker:2375/
    DOCKER_TLS_CERTDIR: ""
  script:
    - apt-get update && apt-get install -y docker-ce-cli
    - go test ./integration -v -p 1
  allow_failure: false

build_docker:
  stage: build
  image: docker:25.0
  services:
    - name: docker:dind
  script:
    - docker build -t $CI_REGISTRY_IMAGE:latest .
    - echo "$CI_REGISTRY_PASSWORD" | docker login -u "$CI_REGISTRY_USER" $CI_REGISTRY --password-stdin
    - docker push $CI_REGISTRY_IMAGE:latest
  only:
    - main
```

---

## ⚙️ Explication du pipeline

| Stage                | Description                                                                       |
| -------------------- | --------------------------------------------------------------------------------- |
| **lint**             | Vérifie format et erreurs de compilation                                          |
| **unit_test**        | Lance les tests unitaires (avec mocks)                                            |
| **integration_test** | Lance les tests d’intégration via Testcontainers + Docker-in-Docker               |
| **build_docker**     | Construit et pousse l’image Docker sur le registre GitLab (uniquement sur `main`) |

---

## 💡 Tips GitLab Runners

* **Runner Docker** requis :
  active l’option “*Privileged*” (Docker-in-Docker).
  Dans GitLab, va dans :
  **Admin → CI/CD → Runners → Edit → Enable privileged mode** ✅
* Si tu veux un PostgreSQL natif sans testcontainers, tu peux ajouter :

```yaml
services:
  - name: postgres:15-alpine
    alias: db
variables:
  POSTGRES_USER: test
  POSTGRES_PASSWORD: test
  POSTGRES_DB: testdb
```

et adapter le DSN dans tes tests.

---

## ✅ Résumé de l’intégration GitLab CI

| Composant                        | Rôle                                                |
| -------------------------------- | --------------------------------------------------- |
| `.gitlab-ci.yml`                 | Orchestration des tests, lint et build Docker       |
| `Dockerfile`                     | Build production minimal (distroless)               |
| `make mock`                      | Génère les mocks avec `mockgen`                     |
| `make test` / `make integration` | Lancement local ou CI                               |
| Testcontainers                   | Fournit un PostgreSQL réel pour tests d’intégration |
| Cache GitLab                     | Accélère les builds Go                              |

---

Souhaites-tu que je t’ajoute aussi un **job GitLab de déploiement automatique** (ex. sur un serveur ou Kubernetes via Helm chart après build) ?
Ça s’intègre très bien après le job `build_docker`.

## CI de deploy

Parfait 👌 tu veux donc étendre ton pipeline GitLab CI pour inclure un **déploiement automatique** après le build — parfait pour un **environnement staging** ou **production (Kubernetes)**.

On va le faire de manière professionnelle et modulaire 🔥

---

## 🎯 Objectif

Ton pipeline GitLab CI va désormais :

1. ✅ **Linter** le code
2. ✅ **Tester** (unitaires + intégration)
3. ✅ **Builder et pousser l’image Docker**
4. 🚀 **Déployer automatiquement sur Kubernetes** (via Helm)

---

## 🧩 Arborescence mise à jour

```shell
huma-pgx-example/
├── .gitlab-ci.yml
├── helm/
│   └── huma-app/
│       ├── Chart.yaml
│       ├── values.yaml
│       └── templates/
│           ├── deployment.yaml
│           └── service.yaml
├── Dockerfile
├── go.mod
├── main.go
├── internal/
│   ├── handler/
│   ├── repository/
│   ├── service/
│   └── mocks/
└── integration/
    └── integration_test.go
```

---

## ⚙️ 1. Helm Chart minimal

### `helm/huma-app/Chart.yaml`

```yaml
apiVersion: v2
name: huma-app
description: Helm chart for Huma Go API
type: application
version: 0.1.0
appVersion: "1.0.0"
```

### `helm/huma-app/values.yaml`

```yaml
image:
  repository: registry.gitlab.com/<YOUR_NAMESPACE>/huma-pgx-example
  tag: latest
  pullPolicy: Always

service:
  port: 8080

resources: {}
```

### `helm/huma-app/templates/deployment.yaml`

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: huma-app
spec:
  replicas: 1
  selector:
    matchLabels:
      app: huma-app
  template:
    metadata:
      labels:
        app: huma-app
    spec:
      containers:
        - name: huma-app
          image: "{{ .Values.image.repository }}:{{ .Values.image.tag }}"
          imagePullPolicy: "{{ .Values.image.pullPolicy }}"
          ports:
            - containerPort: {{ .Values.service.port }}
          env:
            - name: DSN
              value: "postgres://user:password@postgres:5432/db?sslmode=disable"
```

### `helm/huma-app/templates/service.yaml`

```yaml
apiVersion: v1
kind: Service
metadata:
  name: huma-app
spec:
  selector:
    app: huma-app
  ports:
    - port: {{ .Values.service.port }}
      targetPort: {{ .Values.service.port }}
      protocol: TCP
      name: http
  type: ClusterIP
```

---

## 🧰 2. Mise à jour `.gitlab-ci.yml`

Ajout du **stage “deploy”** :

```yaml
stages:
  - lint
  - test
  - integration
  - build
  - deploy

variables:
  GO_VERSION: "1.23"
  KUBE_CONTEXT: "staging"
  KUBE_NAMESPACE: "huma"
  HELM_RELEASE_NAME: "huma-app"
  DOCKER_DRIVER: overlay2

cache:
  key: ${CI_COMMIT_REF_SLUG}
  paths:
    - .cache/go-build
    - .cache/go
    - go.sum
    - go.mod

before_script:
  - mkdir -p .cache/go-build .cache/go
  - go version
  - go mod download

lint:
  stage: lint
  image: golang:${GO_VERSION}
  script:
    - go fmt ./...
    - go vet ./...

unit_test:
  stage: test
  image: golang:${GO_VERSION}
  script:
    - go install go.uber.org/mock/mockgen@latest
    - make mock
    - make test
  artifacts:
    paths:
      - coverage.out

integration_test:
  stage: integration
  image: golang:${GO_VERSION}
  services:
    - name: docker:dind
      command: ["--tls=false"]
  variables:
    DOCKER_HOST: tcp://docker:2375/
    DOCKER_TLS_CERTDIR: ""
  script:
    - apt-get update && apt-get install -y docker-ce-cli
    - go test ./integration -v -p 1
  allow_failure: false

build_docker:
  stage: build
  image: docker:25.0
  services:
    - name: docker:dind
  script:
    - docker build -t $CI_REGISTRY_IMAGE:$CI_COMMIT_SHA .
    - docker tag $CI_REGISTRY_IMAGE:$CI_COMMIT_SHA $CI_REGISTRY_IMAGE:latest
    - echo "$CI_REGISTRY_PASSWORD" | docker login -u "$CI_REGISTRY_USER" $CI_REGISTRY --password-stdin
    - docker push $CI_REGISTRY_IMAGE:$CI_COMMIT_SHA
    - docker push $CI_REGISTRY_IMAGE:latest
  artifacts:
    name: "docker-image"
    paths:
      - helm/huma-app/values.yaml
    expire_in: 1 week
  only:
    - main

deploy_staging:
  stage: deploy
  image:
    name: alpine/helm:3.14.3
    entrypoint: [""]
  script:
    - echo "🔐 Configure kubeconfig"
    - mkdir -p ~/.kube
    - echo "$KUBECONFIG_CONTENT" | base64 -d > ~/.kube/config
    - helm upgrade --install $HELM_RELEASE_NAME ./helm/huma-app \
        --namespace $KUBE_NAMESPACE \
        --create-namespace \
        --set image.repository=$CI_REGISTRY_IMAGE \
        --set image.tag=$CI_COMMIT_SHA \
        --set image.pullPolicy=Always
  only:
    - main
  environment:
    name: staging
    url: https://staging.example.com
```

---

## 🔐 3. Variables GitLab à définir dans le projet

| Variable               | Description                            |
| ---------------------- | -------------------------------------- |
| `CI_REGISTRY_USER`     | Utilisateur du registre GitLab         |
| `CI_REGISTRY_PASSWORD` | Token d’accès au registre              |
| `KUBECONFIG_CONTENT`   | Ton kubeconfig encodé en base64        |
| `KUBE_NAMESPACE`       | Namespace cible sur le cluster         |
| `KUBE_CONTEXT`         | (optionnel) Nom du contexte Kubernetes |

---

### 💡 Pour encoder ton kubeconfig :

```bash
cat ~/.kube/config | base64 -w 0
```

Et colle ce contenu dans la variable GitLab `KUBECONFIG_CONTENT`.

---

## 🚀 Déploiement

À chaque **push sur la branche `main`** :

1. Les tests unitaires et d’intégration passent ✅
2. L’image Docker est buildée et poussée sur le **GitLab Container Registry** 🐳
3. Helm déploie automatiquement sur ton cluster Kubernetes (namespace `huma`) 🎯

---

## ✅ Résumé du flux CI/CD

| Étape              | Action                                            | Environnement      |
| ------------------ | ------------------------------------------------- | ------------------ |
| `lint`             | Format & vet                                      | Runner Docker      |
| `unit_test`        | Tests unitaires + mocks                           | Runner Docker      |
| `integration_test` | Tests d’intégration avec PostgreSQL Testcontainer | Docker-in-Docker   |
| `build_docker`     | Build & push image                                | GitLab Registry    |
| `deploy_staging`   | Déploiement automatique via Helm                  | Kubernetes cluster |

---
