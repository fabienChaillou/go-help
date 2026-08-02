Excellente question 👌 — voyons ensemble **un exemple concret** d’utilisation de [`pgxmock`](https://github.com/pashagolub/pgxmock) pour tester une fonction qui prend un **pool `pgxpool.Pool`** (ou une interface compatible) en paramètre.

---

## 🎯 Objectif

On va écrire une fonction simple :

```go
func GetUserByEmail(ctx context.Context, pool *pgxpool.Pool, email string) (*User, error)
```

et ensuite un **test unitaire** avec `pgxmock` sans vraie base de données.

---

## 📦 Installation

```bash
go get github.com/pashagolub/pgxmock/v4
go get github.com/jackc/pgx/v5/pgxpool
```

---

## 🧱 Exemple de code à tester (`internal/repository/user_repo.go`)

```go
package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID    int64
	Email string
	Name  string
}

// GetUserByEmail fetches a user from the DB by email.
func GetUserByEmail(ctx context.Context, pool *pgxpool.Pool, email string) (*User, error) {
	row := pool.QueryRow(ctx, "SELECT id, email, name FROM users WHERE email=$1", email)

	var u User
	err := row.Scan(&u.ID, &u.Email, &u.Name)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
```

---

## 🧪 Test avec `pgxmock` (`internal/repository/user_repo_test.go`)

```go
package repository

import (
	"context"
	"testing"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/require"
)

// Exemple de test avec pgxmock
func TestGetUserByEmail(t *testing.T) {
	ctx := context.Background()

	// 1️⃣ Crée le mock
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	// 2️⃣ Prépare les données attendues
	expectedUser := User{
		ID:    1,
		Email: "john@example.com",
		Name:  "John Doe",
	}

	// 3️⃣ Configure l’attente du mock
	mock.ExpectQuery("SELECT id, email, name FROM users WHERE email=\\$1").
		WithArgs(expectedUser.Email).
		WillReturnRows(
			pgxmock.NewRows([]string{"id", "email", "name"}).
				AddRow(expectedUser.ID, expectedUser.Email, expectedUser.Name),
		)

	// 4️⃣ Appelle la fonction
	user, err := GetUserByEmail(ctx, mock, expectedUser.Email)
	require.NoError(t, err)
	require.NotNil(t, user)

	// 5️⃣ Vérifie les résultats
	require.Equal(t, expectedUser.ID, user.ID)
	require.Equal(t, expectedUser.Email, user.Email)
	require.Equal(t, expectedUser.Name, user.Name)

	// 6️⃣ Vérifie que toutes les attentes ont été respectées
	require.NoError(t, mock.ExpectationsWereMet())
}
```

---

## 💡 Détails clés

| Étape                                                | Explication                                                                                      |
| ---------------------------------------------------- | ------------------------------------------------------------------------------------------------ |
| `pgxmock.NewPool()`                                  | Crée un faux `pgxpool.Pool` implémentant les mêmes méthodes (`Query`, `QueryRow`, `Exec`, etc.). |
| `ExpectQuery(...).WithArgs(...).WillReturnRows(...)` | Simule la requête SQL et la réponse du serveur.                                                  |
| `mock.ExpectationsWereMet()`                         | Vérifie que toutes les attentes configurées ont été utilisées.                                   |
| `pgxmock.NewRows(columns).AddRow(...)`               | Définit les colonnes et les données retournées.                                                  |

---

## 🧠 Bonus : si tu veux une interface plutôt qu’un *pgxpool direct*

C’est une bonne pratique pour découpler la dépendance à `pgxpool`.
Tu peux définir une interface générique :

```go
type DBQuerier interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}
```

et l’utiliser ainsi :

```go
func GetUserByEmail(ctx context.Context, db DBQuerier, email string) (*User, error)
```

Puis, dans les tests :

```go
mock, _ := pgxmock.NewPool() // implémente DBQuerier
GetUserByEmail(ctx, mock, "john@example.com")
```

Et dans le code réel :

```go
pool, _ := pgxpool.New(ctx, connString)
GetUserByEmail(ctx, pool, "john@example.com")
```

✅ Cela te permet de **mock facilement** sans changer ta logique.

---

## 🧩 Exemple de test échoué (pour comprendre)

```go
mock.ExpectQuery("SELECT id, email, name FROM users WHERE email=\\$1").
    WithArgs("john@example.com").
    WillReturnError(fmt.Errorf("database timeout"))
```

→ La fonction `GetUserByEmail` retournera une erreur, et tu peux la tester :

```go
_, err := GetUserByEmail(ctx, mock, "john@example.com")
require.Error(t, err)
```

---

Souhaites-tu que je t’ajoute aussi un **exemple avec `Exec()`** (insert/update/delete) utilisant `pgxmock.ExpectExec` — typique pour les fonctions `CreateUser` ou `UpdateUser` ?

## example avec Query() && QueryRaw()

Parfait 💪 — tu veux donc des exemples **complets** avec `pgxmock` pour tester :

* une fonction qui utilise `Query()` (retourne plusieurs lignes),
* une fonction qui utilise `QueryRow()` (retourne une seule ligne).

Je te montre les deux cas, avec des **tests unitaires réalistes**.

---

## 🧱 Exemple de modèle simple

```go
package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID    int64
	Email string
	Name  string
}
```

---

## 🔹 Exemple 1 — Fonction avec `QueryRow()`

Cette fonction récupère **un seul utilisateur** via un `WHERE email=$1`.

```go
// GetUserByEmail fetches a single user by email.
func GetUserByEmail(ctx context.Context, pool *pgxpool.Pool, email string) (*User, error) {
	row := pool.QueryRow(ctx, "SELECT id, email, name FROM users WHERE email=$1", email)

	var u User
	if err := row.Scan(&u.ID, &u.Email, &u.Name); err != nil {
		return nil, err
	}
	return &u, nil
}
```

### ✅ Test avec `pgxmock` (`QueryRow`)

```go
package repository

import (
	"context"
	"testing"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/require"
)

func TestGetUserByEmail(t *testing.T) {
	ctx := context.Background()
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	expected := User{ID: 1, Email: "john@example.com", Name: "John Doe"}

	mock.ExpectQuery(`SELECT id, email, name FROM users WHERE email=\$1`).
		WithArgs(expected.Email).
		WillReturnRows(
			pgxmock.NewRows([]string{"id", "email", "name"}).
				AddRow(expected.ID, expected.Email, expected.Name),
		)

	user, err := GetUserByEmail(ctx, mock, expected.Email)
	require.NoError(t, err)
	require.Equal(t, expected, *user)
	require.NoError(t, mock.ExpectationsWereMet())
}
```

---

## 🔹 Exemple 2 — Fonction avec `Query()`

Cette fonction retourne **tous les utilisateurs** avec `Query()`, qu’on scanne dans une boucle.

```go
// ListUsers fetches all users from the database.
func ListUsers(ctx context.Context, pool *pgxpool.Pool) ([]User, error) {
	rows, err := pool.Query(ctx, "SELECT id, email, name FROM users ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Email, &u.Name); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, rows.Err()
}
```

### ✅ Test avec `pgxmock` (`Query`)

```go
func TestListUsers(t *testing.T) {
	ctx := context.Background()
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	rows := pgxmock.NewRows([]string{"id", "email", "name"}).
		AddRow(1, "john@example.com", "John Doe").
		AddRow(2, "jane@example.com", "Jane Doe")

	mock.ExpectQuery(`SELECT id, email, name FROM users ORDER BY id`).
		WillReturnRows(rows)

	users, err := ListUsers(ctx, mock)
	require.NoError(t, err)
	require.Len(t, users, 2)
	require.Equal(t, "Jane Doe", users[1].Name)
	require.NoError(t, mock.ExpectationsWereMet())
}
```

---

## 🔹 Exemple 3 — Fonction avec `Exec()` (bonus)

Tu l’utiliseras souvent pour des `INSERT` ou `UPDATE`.

```go
// CreateUser inserts a new user.
func CreateUser(ctx context.Context, pool *pgxpool.Pool, email, name string) error {
	_, err := pool.Exec(ctx, "INSERT INTO users (email, name) VALUES ($1, $2)", email, name)
	return err
}
```

### ✅ Test `Exec()` avec `pgxmock`

```go
func TestCreateUser(t *testing.T) {
	ctx := context.Background()
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	mock.ExpectExec(`INSERT INTO users \(email, name\) VALUES \(\$1, \$2\)`).
		WithArgs("john@example.com", "John Doe").
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err = CreateUser(ctx, mock, "john@example.com", "John Doe")
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}
```

---

## 🧠 Astuce — Interface générique

Pour découpler ta logique du `pgxpool`, crée une interface commune :

```go
type DB interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, args ...any) (pgx.CommandTag, error)
}
```

Puis tes fonctions deviennent génériques :

```go
func ListUsers(ctx context.Context, db DB) ([]User, error) { ... }
```

Et dans les tests :

```go
mock, _ := pgxmock.NewPool()
ListUsers(ctx, mock)
```

Et dans la prod :

```go
pool, _ := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
ListUsers(ctx, pool)
```

✅ Ainsi tu peux injecter indifféremment un *mock* ou un vrai *pool* sans modifier le code.

---

## 🔍 Récapitulatif rapide

| Type de test         | Méthode pgxmock              | Exemple Go       |
| -------------------- | ---------------------------- | ---------------- |
| SELECT one row       | `ExpectQuery` + `QueryRow()` | `GetUserByEmail` |
| SELECT many rows     | `ExpectQuery` + `Query()`    | `ListUsers`      |
| INSERT/UPDATE/DELETE | `ExpectExec` + `Exec()`      | `CreateUser`     |

---

Souhaites-tu que je te montre **un exemple complet de repository structuré** (interface + implémentation + tests `pgxmock`) pour une entité `UserRepository` avec `Create`, `GetByEmail`, et `List` ?
Ce serait le squelette parfait pour ton module `internal/app/repository`.

Parfait 😎 — voici un **exemple complet de repository structuré** pour l’entité `User` avec **interface, implémentation avec `pgxpool`**, et **tests unitaires avec `pgxmock`**.

L’idée : tu auras un squelette prêt à réutiliser pour **n’importe quelle entité**.

---

## 1️⃣ Interface `UserRepository`

```go
package repository

import "context"

type User struct {
	ID    int64
	Email string
	Name  string
}

type UserRepository interface {
	Create(ctx context.Context, u *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	List(ctx context.Context) ([]User, error)
}
```

---

## 2️⃣ Implémentation avec `pgxpool`

```go
package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepo struct {
	pool *pgxpool.Pool
}

// NewUserRepository retourne l'implémentation concrète
func NewUserRepository(pool *pgxpool.Pool) UserRepository {
	return &userRepo{pool: pool}
}

// Create insère un nouvel utilisateur
func (r *userRepo) Create(ctx context.Context, u *User) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO users (email, name) VALUES ($1, $2)
	`, u.Email, u.Name)
	return err
}

// GetByEmail récupère un utilisateur par email
func (r *userRepo) GetByEmail(ctx context.Context, email string) (*User, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, email, name FROM users WHERE email=$1
	`, email)

	var u User
	if err := row.Scan(&u.ID, &u.Email, &u.Name); err != nil {
		return nil, err
	}
	return &u, nil
}

// List retourne tous les utilisateurs
func (r *userRepo) List(ctx context.Context) ([]User, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, email, name FROM users ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Email, &u.Name); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, rows.Err()
}
```

---

## 3️⃣ Tests unitaires avec `pgxmock`

### Installation

```bash
go get github.com/pashagolub/pgxmock/v4
go get github.com/stretchr/testify
```

---

### `user_repo_test.go`

```go
package repository

import (
	"context"
	"testing"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/require"
)

func setupMockRepo(t *testing.T) (UserRepository, *pgxmock.PgxPoolMock) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)

	repo := NewUserRepository(mock)
	return repo, mock
}

func TestUserRepository_Create(t *testing.T) {
	ctx := context.Background()
	repo, mock := setupMockRepo(t)
	defer mock.Close()

	mock.ExpectExec(`INSERT INTO users \(email, name\) VALUES \(\$1, \$2\)`).
		WithArgs("john@example.com", "John Doe").
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err := repo.Create(ctx, &User{Email: "john@example.com", Name: "John Doe"})
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetByEmail(t *testing.T) {
	ctx := context.Background()
	repo, mock := setupMockRepo(t)
	defer mock.Close()

	mock.ExpectQuery(`SELECT id, email, name FROM users WHERE email=\$1`).
		WithArgs("john@example.com").
		WillReturnRows(
			pgxmock.NewRows([]string{"id", "email", "name"}).
				AddRow(1, "john@example.com", "John Doe"),
		)

	u, err := repo.GetByEmail(ctx, "john@example.com")
	require.NoError(t, err)
	require.Equal(t, int64(1), u.ID)
	require.Equal(t, "John Doe", u.Name)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_List(t *testing.T) {
	ctx := context.Background()
	repo, mock := setupMockRepo(t)
	defer mock.Close()

	mock.ExpectQuery(`SELECT id, email, name FROM users ORDER BY id`).
		WillReturnRows(
			pgxmock.NewRows([]string{"id", "email", "name"}).
				AddRow(1, "john@example.com", "John Doe").
				AddRow(2, "jane@example.com", "Jane Doe"),
		)

	users, err := repo.List(ctx)
	require.NoError(t, err)
	require.Len(t, users, 2)
	require.Equal(t, "Jane Doe", users[1].Name)
	require.NoError(t, mock.ExpectationsWereMet())
}
```

---

## 🔹 Explications

* `setupMockRepo` : crée le mock et injecte dans le repository.
* `ExpectExec` : simule les commandes `INSERT`, `UPDATE`, `DELETE`.
* `ExpectQuery` : simule une seule ligne (`QueryRow`) ou plusieurs lignes (`Query`).
* `ExpectationsWereMet` : vérifie que toutes les attentes du mock ont été respectées.
* On peut ajouter facilement **de nouveaux tests** pour gérer les erreurs ou cas limites.

---

Si tu veux, je peux te faire une **version complète MVC** pour `User` intégrant :

* `UserRepository` avec `pgxpool`,
* `UserService`,
* `UserHandler` (Gin),
* **tests unitaires et intégration** avec `pgxmock`.

Ça te donnerait **le squelette complet prêt pour production et tests**.

Veux‑tu que je fasse ça ?

Parfait 😎 — je vais te créer un **squelette complet MVC pour `User`** en Go, intégrant :

* `UserRepository` (avec `pgxpool`)
* `UserService` (logique métier)
* `UserHandler` (HTTP avec Gin)
* Tests unitaires **repository** avec `pgxmock`
* Tests unitaires **service** avec `mock repository`
* Tests **HTTP handler** (Gin)

Tout sera prêt pour ton projet modulaire.

---

## 1️⃣ Arborescence du projet

```
myapp/
├── cmd/server/main.go
├── internal/
│   ├── app/
│   │   ├── repository/
│   │   │   ├── user_repo.go
│   │   │   └── user_repo_test.go
│   │   ├── service/
│   │   │   ├── user_service.go
│   │   │   └── user_service_test.go
│   │   └── http/
│   │       ├── user_handler.go
│   │       └── user_handler_test.go
│   └── pkg/
│       └── db/db.go
├── go.mod
└── go.sum
```

---

## 2️⃣ Repository `UserRepository`

### `internal/app/repository/user_repo.go`

```go
package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID    int64
	Email string
	Name  string
}

type UserRepository interface {
	Create(ctx context.Context, u *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	List(ctx context.Context) ([]User, error)
}

type userRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) UserRepository {
	return &userRepo{pool: pool}
}

func (r *userRepo) Create(ctx context.Context, u *User) error {
	_, err := r.pool.Exec(ctx, "INSERT INTO users (email, name) VALUES ($1, $2)", u.Email, u.Name)
	return err
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*User, error) {
	row := r.pool.QueryRow(ctx, "SELECT id, email, name FROM users WHERE email=$1", email)
	var u User
	if err := row.Scan(&u.ID, &u.Email, &u.Name); err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) List(ctx context.Context) ([]User, error) {
	rows, err := r.pool.Query(ctx, "SELECT id, email, name FROM users ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Email, &u.Name); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}
```

---

### `internal/app/repository/user_repo_test.go` (pgxmock)

```go
package repository

import (
	"context"
	"testing"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/require"
)

func setupMockRepo(t *testing.T) (UserRepository, *pgxmock.PgxPoolMock) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	repo := NewUserRepository(mock)
	return repo, mock
}

func TestUserRepository_Create(t *testing.T) {
	ctx := context.Background()
	repo, mock := setupMockRepo(t)
	defer mock.Close()

	mock.ExpectExec(`INSERT INTO users \(email, name\) VALUES \(\$1, \$2\)`).
		WithArgs("john@example.com", "John Doe").
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err := repo.Create(ctx, &User{Email: "john@example.com", Name: "John Doe"})
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}
```

*(Les tests `GetByEmail` et `List` sont identiques à ce que je t’avais montré plus haut)*

---

## 3️⃣ Service `UserService`

### `internal/app/service/user_service.go`

```go
package service

import (
	"context"
	"errors"
	"myapp/internal/app/repository"
)

type UserService interface {
	RegisterUser(ctx context.Context, email, name string) (*repository.User, error)
	GetUser(ctx context.Context, email string) (*repository.User, error)
	ListUsers(ctx context.Context) ([]repository.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) RegisterUser(ctx context.Context, email, name string) (*repository.User, error) {
	existing, _ := s.repo.GetByEmail(ctx, email)
	if existing != nil {
		return nil, errors.New("user already exists")
	}
	u := &repository.User{Email: email, Name: name}
	if err := s.repo.Create(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *userService) GetUser(ctx context.Context, email string) (*repository.User, error) {
	return s.repo.GetByEmail(ctx, email)
}

func (s *userService) ListUsers(ctx context.Context) ([]repository.User, error) {
	return s.repo.List(ctx)
}
```

---

### `internal/app/service/user_service_test.go` (mock repository)

```go
package service

import (
	"context"
	"errors"
	"myapp/internal/app/mocks"
	"myapp/internal/app/repository"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/mock"
)

func TestUserService_RegisterUser(t *testing.T) {
	ctx := context.Background()
	mockRepo := &mocks.MockUserRepository{}
	svc := NewUserService(mockRepo)

	mockRepo.On("GetByEmail", ctx, "john@example.com").Return(nil, nil)
	mockRepo.On("Create", ctx, mock.AnythingOfType("*repository.User")).Return(nil)

	user, err := svc.RegisterUser(ctx, "john@example.com", "John Doe")
	require.NoError(t, err)
	require.Equal(t, "John Doe", user.Name)
	mockRepo.AssertExpectations(t)
}

func TestUserService_RegisterUser_AlreadyExists(t *testing.T) {
	ctx := context.Background()
	mockRepo := &mocks.MockUserRepository{}
	svc := NewUserService(mockRepo)

	existing := &repository.User{ID: 1, Email: "john@example.com", Name: "John Doe"}
	mockRepo.On("GetByEmail", ctx, "john@example.com").Return(existing, nil)

	user, err := svc.RegisterUser(ctx, "john@example.com", "John Doe")
	require.Error(t, err)
	require.Nil(t, user)
	mockRepo.AssertExpectations(t)
}
```

---

## 4️⃣ Handler HTTP avec Gin

### `internal/app/http/user_handler.go`

```go
package http

import (
	"net/http"
	"myapp/internal/app/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	svc service.UserService
}

func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) RegisterUser(c *gin.Context) {
	var req struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := h.svc.RegisterUser(c, req.Email, req.Name)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, u)
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	users, err := h.svc.ListUsers(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}
```

---

### `internal/app/http/user_handler_test.go`

```go
package http

import (
	"bytes"
	"encoding/json"
	"myapp/internal/app/mocks"
	"myapp/internal/app/repository"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRegisterUserHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &mocks.MockUserService{}
	handler := NewUserHandler(mockSvc)

	body := map[string]string{"email": "john@example.com", "name": "John Doe"}
	data, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(data))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	mockSvc.On("RegisterUser", mock.Anything, "john@example.com", "John Doe").
		Return(&repository.User{ID: 1, Email: "john@example.com", Name: "John Doe"}, nil)

	handler.RegisterUser(c)
	require.Equal(t, http.StatusCreated, w.Code)
	mockSvc.AssertExpectations(t)
}
```

---

## 5️⃣ Main `cmd/server/main.go`

```go
package main

import (
	"log"
	"myapp/internal/app/http"
	"myapp/internal/app/repository"
	"myapp/internal/app/service"
	"myapp/internal/pkg/db"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	pool, err := db.NewSQLX(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("DB init failed: %v", err)
	}
	defer pool.Close()

	userRepo := repository.NewUserRepository(pool)
	userSvc := service.NewUserService(userRepo)
	userHandler := http.NewUserHandler(userSvc)

	r := gin.Default()
	r.POST("/users", userHandler.RegisterUser)
	r.GET("/users", userHandler.ListUsers)

	log.Println("Server running on :8080")
	r.Run(":8080")
}
```

---

✅ Avec cette structure, tu as :

* **Repository** testé avec `pgxmock`
* **Service** testé avec mock repository
* **Handler HTTP** testé avec mock service
* **Main** prêt pour Gin et PostgreSQL

Tout est modulaire et testable.

---

Si tu veux, je peux te fournir **la version finale Docker + docker-compose + golang-migrate** prête pour CI/CD avec GitLab, intégrant ce squelette MVC.
Veux‑tu que je fasse ça ?

Parfait 😎 — je vais te créer **le projet complet prêt pour Docker + docker-compose + golang-migrate + GitLab CI/CD**, intégrant le squelette MVC `User` que je t’ai montré.

Tout sera prêt pour **développement local, tests unitaires et intégration, migrations et CI/CD GitLab**.

---

## 1️⃣ Arborescence finale du projet

```
myapp/
├── cmd/server/main.go
├── internal/
│   ├── app/
│   │   ├── repository/
│   │   │   ├── user_repo.go
│   │   │   └── user_repo_test.go
│   │   ├── service/
│   │   │   ├── user_service.go
│   │   │   └── user_service_test.go
│   │   └── http/
│   │       ├── user_handler.go
│   │       └── user_handler_test.go
│   └── pkg/
│       └── db/db.go
├── migrations/
│   ├── 001_create_users.up.sql
│   └── 001_create_users.down.sql
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
└── .gitlab-ci.yml
```

---

## 2️⃣ Dockerfile

```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /bin/myapp ./cmd/server

# Runtime
FROM gcr.io/distroless/static
COPY --from=builder /bin/myapp /bin/myapp
EXPOSE 8080
ENTRYPOINT ["/bin/myapp"]
```

---

## 3️⃣ docker-compose.yml

```yaml
version: "3.8"
services:
  db:
    image: postgres:15
    environment:
      POSTGRES_USER: myapp
      POSTGRES_PASSWORD: myapppw
      POSTGRES_DB: myapp
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
      DATABASE_URL: postgres://myapp:myapppw@db:5432/myapp?sslmode=disable
    ports:
      - "8080:8080"
    depends_on:
      db:
        condition: service_healthy

volumes:
  db_data:
```

---

## 4️⃣ Migrations

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

## 5️⃣ DB helper

### `internal/pkg/db/db.go`

```go
package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewSQLX(databaseURL string) (*pgxpool.Pool, error) {
	ctx := context.Background()
	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		return nil, err
	}
	return pool, nil
}
```

---

## 6️⃣ golang-migrate integration

### `internal/pkg/migrate/migrate.go`

```go
package migrate

import (
	"fmt"
	"log"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(databaseURL, migrationsPath string) {
	m, err := migrate.New(fmt.Sprintf("file://%s", migrationsPath), databaseURL)
	if err != nil {
		log.Fatalf("Failed to load migrations: %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Migration failed: %v", err)
	}
	log.Println("✅ Migrations applied successfully")
}
```

---

## 7️⃣ GitLab CI/CD

### `.gitlab-ci.yml`

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

services:
  - name: postgres:15
    alias: postgres

before_script:
  - apt-get update -qq && apt-get install -y curl unzip
  - go mod download
  - until pg_isready -h postgres -U $POSTGRES_USER; do echo "Waiting for postgres..."; sleep 1; done

lint:
  stage: lint
  image: golang:${GO_VERSION}-alpine
  script:
    - go install golang.org/x/lint/golint@latest
    - golint ./...
    - go vet ./...
  allow_failure: true

unit_tests:
  stage: test
  image: golang:${GO_VERSION}-alpine
  script:
    - go test ./internal/... -v -coverprofile=coverage.out
  artifacts:
    paths:
      - coverage.out

integration_tests:
  stage: test
  image: golang:${GO_VERSION}
  script:
    - wget -q https://github.com/golang-migrate/migrate/releases/download/v4.17.1/migrate.linux-amd64.tar.gz
    - tar -xzf migrate.linux-amd64.tar.gz -C /usr/local/bin
    - chmod +x /usr/local/bin/migrate
    - migrate -path migrations -database "$DATABASE_URL" up
    - go test ./tests -v -count=1

build:
  stage: build
  image: golang:${GO_VERSION}
  script:
    - go build -o bin/myapp ./cmd/server
  artifacts:
    paths:
      - bin/myapp
```

---

## 8️⃣ Commandes utiles

```bash
# Lancer la stack localement
docker-compose up --build

# Appliquer les migrations localement
docker-compose exec app migrate -path migrations -database $DATABASE_URL up

# Lancer l'API
go run ./cmd/server

# Tests unitaires
go test ./internal/...

# Tests d'intégration (via container Postgres)
go test ./tests -v
```

---

✅ Avec ça, tu as un **projet complet prêt pour :**

* Développement local avec PostgreSQL via Docker
* Migrations versionnées via golang-migrate
* Architecture MVC pour `User`
* Tests unitaires et d’intégration (pgxmock / testcontainers)
* CI/CD GitLab prêt à l’emploi

---

Si tu veux, je peux te faire maintenant **un exemple prêt à copier-coller** avec **tests d’intégration réels et pgxmock** inclus dans le projet, pour que tu puisses lancer tout avec un seul `docker-compose up`.

## example avec testcontainers-go

Parfait 😎 — je vais te montrer **un exemple complet de test d’intégration avec `testcontainers-go`**, utilisant PostgreSQL dans un conteneur Docker, sans toucher à ta base locale.

Ce test est **100 % Go**, indépendant de Docker Compose, et fonctionne directement dans ton pipeline CI/CD.

---

## 1️⃣ Installation

```bash
go get github.com/testcontainers/testcontainers-go
go get github.com/testcontainers/testcontainers-go/wait
go get github.com/jackc/pgx/v5/pgxpool
go get github.com/stretchr/testify
```

---

## 2️⃣ Exemple de test d’intégration

### `internal/repository/user_integration_test.go`

```go
package repository

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupPostgresContainer(ctx context.Context) (tc.Container, string, error) {
	req := tc.ContainerRequest{
		Image:        "postgres:15",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "myapp",
			"POSTGRES_PASSWORD": "myapppw",
			"POSTGRES_DB":       "myapp_test",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(30 * time.Second),
	}
	container, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, "", err
	}

	host, err := container.Host(ctx)
	if err != nil {
		return nil, "", err
	}
	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		return nil, "", err
	}

	databaseURL := fmt.Sprintf("postgres://myapp:myapppw@%s:%s/myapp_test?sslmode=disable", host, port.Port())
	return container, databaseURL, nil
}

func TestUserRepository_Integration(t *testing.T) {
	ctx := context.Background()

	// 1️⃣ Démarre le conteneur Postgres
	container, databaseURL, err := setupPostgresContainer(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	// 2️⃣ Connexion à la DB
	pool, err := pgxpool.New(ctx, databaseURL)
	require.NoError(t, err)
	defer pool.Close()

	// 3️⃣ Crée la table users
	_, err = pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS users (
			id BIGSERIAL PRIMARY KEY,
			email TEXT NOT NULL UNIQUE,
			name TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
	`)
	require.NoError(t, err)

	repo := NewUserRepository(pool)

	// 4️⃣ Test Create
	user := &User{Email: "john@example.com", Name: "John Doe"}
	err = repo.Create(ctx, user)
	require.NoError(t, err)

	// 5️⃣ Test GetByEmail
	got, err := repo.GetByEmail(ctx, "john@example.com")
	require.NoError(t, err)
	require.Equal(t, "John Doe", got.Name)

	// 6️⃣ Test List
	users, err := repo.List(ctx)
	require.NoError(t, err)
	require.Len(t, users, 1)
	require.Equal(t, "john@example.com", users[0].Email)

	log.Println("✅ Test d’intégration Postgres avec Testcontainers terminé")
}
```

---

## 3️⃣ Explications clés

| Étape                                | Explication                                                 |
| ------------------------------------ | ----------------------------------------------------------- |
| `setupPostgresContainer`             | Crée et démarre un conteneur PostgreSQL isolé pour le test. |
| `pgxpool.New`                        | Connexion au conteneur via `pgxpool`.                       |
| `CREATE TABLE`                       | On crée la table pour le test d’intégration.                |
| Tests `Create`, `GetByEmail`, `List` | Exécutés **réellement** sur Postgres.                       |
| `defer container.Terminate(ctx)`     | Nettoie automatiquement le conteneur après le test.         |

✅ Tu n’as pas besoin d’une DB locale, le test est **complètement isolé**.

---

## 4️⃣ Exécution avec Dockerfile

Si tu veux exécuter ces tests dans ton pipeline ou localement via Docker :

```dockerfile
# Dockerfile pour tests
FROM golang:1.21-alpine
WORKDIR /app
RUN apk add --no-cache git bash
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Installer Testcontainers dépendances
RUN go install github.com/testcontainers/testcontainers-go@latest

# Lancer les tests
CMD ["go", "test", "./internal/repository", "-v", "-count=1"]
```

Puis build & run :

```bash
docker build -t myapp-test .
docker run --rm myapp-test
```

Tu verras le test d’intégration complet s’exécuter avec un conteneur Postgres temporaire. ✅

---

Si tu veux, je peux te préparer **un Docker Compose complet** qui combine :

* L’API Go
* PostgreSQL pour dev
* Et un service de test qui exécute automatiquement les tests d’intégration avec `testcontainers-go`.

Veux‑tu que je fasse ça ?

Parfait 😎 — je vais te créer un **Docker Compose complet** pour ton projet Go MVC avec :

* **API Go** (`myapp`)
* **PostgreSQL** pour dev
* **Service de test** qui exécute automatiquement les tests d’intégration avec `testcontainers-go`

Tout sera prêt à copier-coller.

---

## 1️⃣ Structure du projet

```
myapp/
├── cmd/server/main.go
├── internal/...
├── migrations/...
├── go.mod
├── go.sum
├── Dockerfile
└── docker-compose.yml
```

---

## 2️⃣ Dockerfile

```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git bash curl
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /bin/myapp ./cmd/server

# Runtime
FROM gcr.io/distroless/static
COPY --from=builder /bin/myapp /bin/myapp
EXPOSE 8080
ENTRYPOINT ["/bin/myapp"]
```

---

## 3️⃣ docker-compose.yml complet

```yaml
version: "3.8"

services:
  db:
    image: postgres:15
    environment:
      POSTGRES_USER: myapp
      POSTGRES_PASSWORD: myapppw
      POSTGRES_DB: myapp
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
      DATABASE_URL: postgres://myapp:myapppw@db:5432/myapp?sslmode=disable
    ports:
      - "8080:8080"
    depends_on:
      db:
        condition: service_healthy

  tests:
    build: .
    environment:
      DATABASE_URL: postgres://myapp:myapppw@db:5432/myapp?sslmode=disable
    depends_on:
      - db
    command: >
      sh -c "
        go test ./internal/repository -v -count=1 &&
        go test ./internal/service -v -count=1 &&
        go test ./internal/http -v -count=1
      "

volumes:
  db_data:
```

---

## 4️⃣ Explications

| Service        | Description                                                                                      |
| -------------- | ------------------------------------------------------------------------------------------------ |
| **db**         | PostgreSQL pour dev & tests.                                                                     |
| **app**        | Ton API Go compilée depuis Dockerfile.                                                           |
| **tests**      | Container qui exécute tous les tests unitaires et d’intégration. Utilise le même code que `app`. |
| **depends_on** | Assure que `db` est prêt avant l’API ou les tests.                                               |

✅ Avec ça, tu peux :

```bash
docker-compose up --build
```

* Le service `db` démarre PostgreSQL.
* Le service `app` démarre l’API sur le port 8080.
* Le service `tests` exécute **les tests unitaires et d’intégration** avec `testcontainers-go`.

Les conteneurs de test se terminent automatiquement après l’exécution des tests.

---

Si tu veux, je peux te préparer **la version finale GitLab CI/CD** qui combine **Docker Compose + tests d’intégration `testcontainers-go` + golang-migrate** prête à lancer dans ton pipeline.

Veux‑tu que je fasse ça ?

Parfait 😎 — je vais te créer un **`.gitlab-ci.yml` complet** prêt pour ton projet Go MVC avec Docker Compose, `golang-migrate` et tests d’intégration via `testcontainers-go`.

Ce pipeline fera :

1. Lint + vet
2. Tests unitaires (`pgxmock`)
3. Tests d’intégration réels (`testcontainers-go`)
4. Build de l’API Docker

---

## `.gitlab-ci.yml` complet

```yaml
stages:
  - lint
  - test
  - build

variables:
  GO_VERSION: "1.21"
  DATABASE_URL: postgres://myapp:myapppw@db:5432/myapp_test?sslmode=disable

# --- Stage 1 : Lint ---
lint:
  stage: lint
  image: golang:${GO_VERSION}-alpine
  script:
    - go install golang.org/x/lint/golint@latest
    - golint ./...
    - go vet ./...
  allow_failure: true

# --- Stage 2 : Tests unitaires ---
unit_tests:
  stage: test
  image: golang:${GO_VERSION}-alpine
  services:
    - name: postgres:15
      alias: db
  variables:
    POSTGRES_DB: myapp_test
    POSTGRES_USER: myapp
    POSTGRES_PASSWORD: myapppw
  before_script:
    - apk add --no-cache bash
    - until pg_isready -h db -U $POSTGRES_USER; do echo "Waiting for Postgres..."; sleep 1; done
  script:
    - go test ./internal/... -v -coverprofile=coverage.out
  artifacts:
    paths:
      - coverage.out

# --- Stage 3 : Tests d'intégration avec testcontainers-go ---
integration_tests:
  stage: test
  image: golang:${GO_VERSION}
  services:
    - name: docker:dind
      alias: docker
  variables:
    DOCKER_DRIVER: overlay2
  before_script:
    - apk add --no-cache bash git curl
    - go mod download
  script:
    - go test ./internal/repository -v -count=1
    - go test ./internal/service -v -count=1
    - go test ./internal/http -v -count=1

# --- Stage 4 : Build Docker ---
build:
  stage: build
  image: docker:24.0.6
  services:
    - docker:dind
  variables:
    DOCKER_DRIVER: overlay2
  before_script:
    - docker info
  script:
    - docker build -t registry.gitlab.com/$CI_PROJECT_PATH/myapp:latest .
    - docker push registry.gitlab.com/$CI_PROJECT_PATH/myapp:latest
```

---

## 🔹 Explications

| Stage                 | Description                                                                                                         |
| --------------------- | ------------------------------------------------------------------------------------------------------------------- |
| **lint**              | Vérifie la qualité du code avec `golint` et `go vet`.                                                               |
| **unit_tests**        | Tests unitaires sur `pgxmock` avec Postgres minimal.                                                                |
| **integration_tests** | Tests d’intégration réels via `testcontainers-go` ; démarre des conteneurs temporaires PostgreSQL pour chaque test. |
| **build**             | Compile l’API et pousse l’image Docker dans le registry GitLab.                                                     |

---

## 🔹 Points clés

1. **Testcontainers-go fonctionne dans CI/CD** car `docker:dind` est disponible.
2. Tu n’as pas besoin d’une DB fixe pour les tests d’intégration.
3. `golang-migrate` peut être intégré dans le job `integration_tests` si tu veux appliquer des migrations avant les tests :

```bash
wget -q https://github.com/golang-migrate/migrate/releases/download/v4.17.1/migrate.linux-amd64.tar.gz
tar -xzf migrate.linux-amd64.tar.gz -C /usr/local/bin
chmod +x /usr/local/bin/migrate
migrate -path migrations -database "$DATABASE_URL" up
```

---

Avec cette configuration :

* Tu as un **pipeline complet GitLab**
* Tests unitaires et d’intégration isolés
* Build Docker prêt à pousser dans le registry
* Compatible avec `testcontainers-go` et `golang-migrate`

---
