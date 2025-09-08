Les tests en Go (Golang) sont intégrés de manière native dans le langage via le package standard `testing`. Cela permet d’écrire, d’exécuter et de structurer des tests unitaires, d’intégration et de performance de manière simple et idiomatique.

---

### 🧪 1. **Structure de base d’un test unitaire**

Un fichier de test en Go :

* Doit se trouver dans le même package que le code à tester (ou dans `xxx_test` pour tester en boîte noire).
* Son nom doit se terminer par `_test.go`.
* Chaque fonction de test doit commencer par `Test` et prendre un argument `*testing.T`.

**Exemple :**

```go
// fichier: math.go
package mathutils

func Add(a, b int) int {
    return a + b
}
```

```go
// fichier: math_test.go
package mathutils

import "testing"

func TestAdd(t *testing.T) {
    result := Add(2, 3)
    if result != 5 {
        t.Errorf("Add(2, 3) = %d; want 5", result)
    }
}
```

---

### ✅ 2. **Exécuter les tests**

Dans le terminal :

```bash
go test ./...
```

Ou dans un dossier spécifique :

```bash
go test
```

Avec plus de détails :

```bash
go test -v
```

---

### 🔁 3. **Table-driven tests** (très idiomatique en Go)

Permet de tester plusieurs cas dans une boucle :

```go
func TestAdd(t *testing.T) {
    tests := []struct {
        a, b, want int
    }{
        {1, 1, 2},
        {2, 3, 5},
        {10, -2, 8},
    }

    for _, tt := range tests {
        got := Add(tt.a, tt.b)
        if got != tt.want {
            t.Errorf("Add(%d, %d) = %d; want %d", tt.a, tt.b, got, tt.want)
        }
    }
}
```

---

### 🧵 4. **Tests d’erreurs et cas limites**

Utilise souvent `t.Fatal`, `t.Log`, `t.Helper` pour clarifier les erreurs :

```go
func TestDivideByZero(t *testing.T) {
    defer func() {
        if r := recover(); r == nil {
            t.Error("expected panic on divide by zero")
        }
    }()
    _ = 10 / 0
}
```

---

### 🧰 5. **Autres types de tests**

* **Benchmarks** : commencent par `BenchmarkXxx(b *testing.B)`
* **Examples** : commencent par `ExampleXxx()` et incluent un commentaire `// Output:` pour valider la sortie

```go
func BenchmarkAdd(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Add(1, 2)
    }
}
```

---

Souhaites-tu un exemple de tests pour un handler HTTP, un repository, ou autre chose plus spécifique ?

Parfait. Voici **un exemple de test pour un handler HTTP** avec le framework **Gin** (puisqu'on travaille dessus dans ton projet).

---

### 🎯 **Objectif : tester un handler `GET /users/:id`**

#### Exemple de handler :

```go
// fichier: handlers/user.go
package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

func GetUserHandler(c *gin.Context) {
    id := c.Param("id")
    if id == "1" {
        c.JSON(http.StatusOK, User{ID: 1, Name: "Alice"})
    } else {
        c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
    }
}
```

---

### 🧪 Test du handler

```go
// fichier: handlers/user_test.go
package handlers

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
)

func TestGetUserHandler(t *testing.T) {
    gin.SetMode(gin.TestMode)

    router := gin.Default()
    router.GET("/users/:id", GetUserHandler)

    // Cas utilisateur trouvé
    req, _ := http.NewRequest(http.MethodGet, "/users/1", nil)
    resp := httptest.NewRecorder()
    router.ServeHTTP(resp, req)

    assert.Equal(t, http.StatusOK, resp.Code)
    assert.Contains(t, resp.Body.String(), "Alice")

    // Cas utilisateur non trouvé
    req2, _ := http.NewRequest(http.MethodGet, "/users/99", nil)
    resp2 := httptest.NewRecorder()
    router.ServeHTTP(resp2, req2)

    assert.Equal(t, http.StatusNotFound, resp2.Code)
    assert.Contains(t, resp2.Body.String(), "user not found")
}
```

---

### 🧰 Notes :

* On utilise `httptest.NewRecorder()` pour capturer la réponse HTTP.
* `assert` vient du package [testify](https://github.com/stretchr/testify), pratique pour des assertions plus lisibles.

---

Souhaites-tu que je t’aide à écrire des tests pour un repo, un service, ou un middleware aussi ?

Parfait ! Voici un **exemple de test pour un repository** en Go, adapté à une base SQLite en mémoire (idéal pour les tests rapides).

---

### 📦 **Repository de base : `UserRepository`**

#### Exemple : `repository/user.go`

```go
package repository

import (
    "database/sql"
)

type User struct {
    ID   int
    Name string
}

type UserRepository struct {
    DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{DB: db}
}

func (r *UserRepository) GetByID(id int) (*User, error) {
    user := &User{}
    row := r.DB.QueryRow("SELECT id, name FROM users WHERE id = ?", id)
    err := row.Scan(&user.ID, &user.Name)
    if err != nil {
        return nil, err
    }
    return user, nil
}
```

---

### 🧪 **Test avec SQLite en mémoire**

#### Fichier : `repository/user_test.go`

```go
package repository

import (
    "database/sql"
    "testing"

    _ "github.com/mattn/go-sqlite3"
    "github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) *sql.DB {
    db, err := sql.Open("sqlite3", ":memory:")
    if err != nil {
        t.Fatalf("failed to open db: %v", err)
    }

    _, err = db.Exec(`CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT)`)
    if err != nil {
        t.Fatalf("failed to create table: %v", err)
    }

    _, err = db.Exec(`INSERT INTO users (id, name) VALUES (1, 'Alice'), (2, 'Bob')`)
    if err != nil {
        t.Fatalf("failed to insert users: %v", err)
    }

    return db
}

func TestUserRepository_GetByID(t *testing.T) {
    db := setupTestDB(t)
    repo := NewUserRepository(db)

    user, err := repo.GetByID(1)
    assert.NoError(t, err)
    assert.Equal(t, 1, user.ID)
    assert.Equal(t, "Alice", user.Name)

    _, err = repo.GetByID(99)
    assert.Error(t, err)
}
```

---

### ✅ Points importants :

* On utilise une **base SQLite en mémoire** (`:memory:`) pour des tests propres et jetables.
* `testify/assert` facilite les vérifications.
* On injecte la base dans le repository : c'est testable et modulaire.

---

Souhaites-tu aussi un exemple de **test pour un middleware** (par exemple, une authentification) ?

Parfait. Voici un exemple **complet de test pour un middleware en Go avec Gin**, par exemple un middleware **d’authentification par token**.

---

### 🔐 **1. Middleware d’authentification simple**

#### Fichier : `middleware/auth.go`

```go
package middleware

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token != "Bearer secret-token" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
            return
        }
        c.Next()
    }
}
```

---

### 🧪 **2. Test du middleware**

#### Fichier : `middleware/auth_test.go`

```go
package middleware

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
    gin.SetMode(gin.TestMode)

    r := gin.New()
    r.Use(AuthMiddleware())

    // Handler protégé par le middleware
    r.GET("/protected", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"message": "success"})
    })

    // Cas : pas de token
    req1, _ := http.NewRequest(http.MethodGet, "/protected", nil)
    resp1 := httptest.NewRecorder()
    r.ServeHTTP(resp1, req1)

    assert.Equal(t, http.StatusUnauthorized, resp1.Code)
    assert.Contains(t, resp1.Body.String(), "unauthorized")

    // Cas : mauvais token
    req2, _ := http.NewRequest(http.MethodGet, "/protected", nil)
    req2.Header.Set("Authorization", "Bearer wrong-token")
    resp2 := httptest.NewRecorder()
    r.ServeHTTP(resp2, req2)
    assert.Equal(t, http.StatusUnauthorized, resp2.Code)

    // Cas : bon token
    req3, _ := http.NewRequest(http.MethodGet, "/protected", nil)
    req3.Header.Set("Authorization", "Bearer secret-token")
    resp3 := httptest.NewRecorder()
    r.ServeHTTP(resp3, req3)

    assert.Equal(t, http.StatusOK, resp3.Code)
    assert.Contains(t, resp3.Body.String(), "success")
}
```

---

### ✅ Ce que ça couvre :

* Middleware qui bloque les requêtes sans ou avec un mauvais token.
* Accès autorisé uniquement avec `Authorization: Bearer secret-token`.

---
