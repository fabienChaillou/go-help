Excellente question 👌 ! Les **génériques** (generics) ont été introduits dans **Go 1.18** pour permettre l’écriture de fonctions, types, et structures **indépendants d’un type concret**, tout en conservant le **typage fort**.

---

## 🧠 Qu’est-ce qu’un *générique* en Go ?

Un **générique** permet d’écrire du code **réutilisable** avec différents types, **sans dupliquer** la logique.
C’est une **abstraction de type** (comme en Java, C++, Rust…).

---

## ✅ Syntaxe de base

```go
func PrintSlice[T any](s []T) {
    for _, v := range s {
        fmt.Println(v)
    }
}
```

* `T` est un **paramètre de type**
* `any` signifie : **n’importe quel type** (équivalent à `interface{}`)
* `T` peut ensuite être utilisé comme un **type concret** dans la fonction

---

## 🧪 Exemple simple

```go
func Max[T constraints.Ordered](a, b T) T {
    if a > b {
        return a
    }
    return b
}
```

* `constraints.Ordered` : restreint `T` aux types qui peuvent être comparés (`int`, `float64`, `string`, etc.)
* Appel : `Max(5, 10)` ou `Max("a", "z")`

---

## 📦 Type générique

```go
type Stack[T any] struct {
    items []T
}

func (s *Stack[T]) Push(v T) {
    s.items = append(s.items, v)
}

func (s *Stack[T]) Pop() T {
    n := len(s.items)
    val := s.items[n-1]
    s.items = s.items[:n-1]
    return val
}
```

👆 Ici :

* `Stack[int]` est une pile d’entiers
* `Stack[string]` est une pile de chaînes

---

## ✅ Avantages des génériques

| Avantage           | Détail                                      |
| ------------------ | ------------------------------------------- |
| 🔁 Réutilisabilité | Une seule fonction/type pour tous les types |
| 🔒 Sécurité        | Pas besoin de `interface{}` + assertions    |
| ⚡️ Performance     | Pas d’allocation inutile ou cast            |

---

## 🧩 Cas d’usage typiques

* Fonctions utilitaires : `Map`, `Filter`, `Max`, `Min`
* Structures communes : `Stack[T]`, `Queue[T]`, `Repository[T]`
* API orientées données génériques (pagination, etc.)

---

## 🧪 Exemple concret : `Map`

```go
func Map[T any, U any](in []T, f func(T) U) []U {
    out := make([]U, len(in))
    for i, v := range in {
        out[i] = f(v)
    }
    return out
}

// Utilisation
nums := []int{1, 2, 3}
squares := Map(nums, func(n int) int { return n * n }) // [1 4 9]
```

---

## 🧠 Résumé

| Élément          | Description                                            |
| ---------------- | ------------------------------------------------------ |
| `T any`          | Type générique sans restriction                        |
| `T comparable`   | Type qui supporte `==`, `!=`                           |
| `T Ordered`      | Type qui supporte `<`, `>` (via `constraints.Ordered`) |
| `func[T any]...` | Fonction générique                                     |
| `type X[T any]`  | Type struct ou interface générique                     |

---

Souhaites-tu un exemple **avancé** avec des **contraintes personnalisées** ou un **repository générique** avec filtre/tri/pagination ?

---------------------------------------------------------------------------------------------------------------------------

Parfait ! Voyons un **exemple avancé** de génériques en Go avec :

* ✅ une **contrainte personnalisée**
* ✅ un **repository générique** (CRUD minimal)
* ✅ des fonctions avec **filtrage / tri / pagination en mémoire**

---

## 🧠 Étape 1 : Créer une contrainte personnalisée

On définit une **interface générique** qui impose que le type ait un champ `ID int`.

### 📁 `constraints.go`

```go
package constraints

type Entity interface {
    GetID() int
}
```

Exemple :

```go
type User struct {
    ID   int
    Name string
}

func (u User) GetID() int { return u.ID }
```

---

## 🧱 Étape 2 : Repository générique en mémoire

### 📁 `repository.go`

```go
package repository

import (
    "yourmodule/constraints"
    "sort"
)

type InMemoryRepo[T constraints.Entity] struct {
    data []T
}

func NewInMemoryRepo[T constraints.Entity]() *InMemoryRepo[T] {
    return &InMemoryRepo[T]{}
}

func (r *InMemoryRepo[T]) Add(item T) {
    r.data = append(r.data, item)
}

func (r *InMemoryRepo[T]) GetByID(id int) (T, bool) {
    for _, v := range r.data {
        if v.GetID() == id {
            return v, true
        }
    }
    var zero T
    return zero, false
}

func (r *InMemoryRepo[T]) Filter(f func(T) bool) []T {
    out := []T{}
    for _, v := range r.data {
        if f(v) {
            out = append(out, v)
        }
    }
    return out
}

func (r *InMemoryRepo[T]) Sort(less func(a, b T) bool) {
    sort.Slice(r.data, func(i, j int) bool {
        return less(r.data[i], r.data[j])
    })
}

func (r *InMemoryRepo[T]) Paginate(page, size int) []T {
    start := (page - 1) * size
    end := start + size
    if start >= len(r.data) {
        return []T{}
    }
    if end > len(r.data) {
        end = len(r.data)
    }
    return r.data[start:end]
}
```

---

## 🧪 Étape 3 : Utilisation avec `User`

```go
package main

import (
    "fmt"
    "yourmodule/repository"
)

type User struct {
    ID   int
    Name string
}

func (u User) GetID() int { return u.ID }

func main() {
    repo := repository.NewInMemoryRepo[User]()

    repo.Add(User{ID: 1, Name: "Alice"})
    repo.Add(User{ID: 2, Name: "Bob"})
    repo.Add(User{ID: 3, Name: "Charlie"})

    // Filtrer
    usersWithB := repo.Filter(func(u User) bool {
        return u.Name[0] == 'B'
    })
    fmt.Println("Filtrés:", usersWithB)

    // Trier
    repo.Sort(func(a, b User) bool {
        return a.Name < b.Name
    })

    // Paginer
    page := repo.Paginate(1, 2)
    fmt.Println("Page 1:", page)
}
```

---

## ✅ Résultat

| Fonction     | Description           |
| ------------ | --------------------- |
| `Add()`      | Ajout d’un élément    |
| `GetByID()`  | Récupération par ID   |
| `Filter()`   | Filtrage personnalisé |
| `Sort()`     | Tri personnalisé      |
| `Paginate()` | Pagination simple     |

---

---------------------------------------------------------------------------------------------------------------------------

Souhaites-tu que je t’ajoute une **version SQL (Postgres/MySQL) avec pagination dynamique** ou une **version REST avec `gin`/`fiber`** utilisant ce repository ?

Parfait ! Passons à l’étape suivante : **exposer ce repository générique via une API REST** en Go.
On va utiliser le framework **[Gin](https://github.com/gin-gonic/gin)**, simple et rapide à mettre en place.

---

## 🎯 Objectif

Créer une API REST qui expose :

* `GET /users` : liste paginée
* `POST /users` : ajout d’un utilisateur
* `GET /users/:id` : récupère un user par ID

Utilise :
✅ Génériques
✅ Repository en mémoire
✅ Struct `User`

---

## 🧰 Prérequis

```bash
go get github.com/gin-gonic/gin
```

---

## 📁 `/main.go`

```go
package main

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "yourmodule/repository"
)

// --- Struct et implémentation de Entity ---
type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

func (u User) GetID() int { return u.ID }

// --- Repository générique ---
var userRepo = repository.NewInMemoryRepo[User]()

func main() {
    r := gin.Default()

    // Routes
    r.POST("/users", createUser)
    r.GET("/users/:id", getUserByID)
    r.GET("/users", listUsers)

    r.Run(":8080")
}

// --- Handlers ---
func createUser(c *gin.Context) {
    var user User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    userRepo.Add(user)
    c.JSON(http.StatusCreated, user)
}

func getUserByID(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    user, found := userRepo.GetByID(id)
    if !found {
        c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
        return
    }
    c.JSON(http.StatusOK, user)
}

func listUsers(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

    users := userRepo.Paginate(page, size)
    c.JSON(http.StatusOK, users)
}
```

---

## 🧪 Exemple d’appel via `curl`

```bash
curl -X POST http://localhost:8080/users \
    -H "Content-Type: application/json" \
    -d '{"id": 1, "name": "Alice"}'
```

```bash
curl http://localhost:8080/users?page=1&size=2
```

---

## ✅ Ce que tu gagnes

| Élément            | Bénéfice                                             |
| ------------------ | ---------------------------------------------------- |
| 🌐 REST API        | Interface universelle                                |
| 🔁 Génériques      | Réutilisable pour d’autres entités (`Product`, etc.) |
| 🧪 Facile à tester | Pas besoin de base de données au début               |

---

---------------------------------------------------------------------------------------------------------------------------

Souhaites-tu maintenant :

* ➕ Ajouter une **base de données (PostgreSQL)** derrière ce repo ?
* 🛡️ Ajouter une validation (`validator`, `binding` de Gin) ?
* 🧪 Ajouter des **tests HTTP automatisés** avec `httptest` ?

Parfait ! On va maintenant **remplacer le repository en mémoire** par un repository connecté à **PostgreSQL**, en gardant la logique REST avec **Gin** et les **génériques**.

---

## 🎯 Objectif final

🔗 API REST (Gin)
💽 Repository générique (`SQLRepository[T]`)
🗃️ Stockage réel PostgreSQL (via `sqlx`)
✅ Fonctionnalités : `GET`, `POST`, `GET by ID` sur `/users`

---

## 🧰 Prérequis

1. PostgreSQL installé ou via Docker :

   ```bash
   docker run --name pg -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres
   ```

2. Créer la table :

   ```sql
   CREATE DATABASE testdb;

   \c testdb

   CREATE TABLE users (
     id SERIAL PRIMARY KEY,
     name TEXT NOT NULL
   );
   ```

3. Installer les libs :

   ```bash
   go get github.com/gin-gonic/gin
   go get github.com/jmoiron/sqlx
   go get github.com/lib/pq
   go get github.com/Masterminds/squirrel
   ```

---

## 🧩 `repository/sql_repository.go` – version générique PostgreSQL

```go
package repository

import (
    "fmt"
    "reflect"

    sq "github.com/Masterminds/squirrel"
    "github.com/jmoiron/sqlx"
)

type SQLRepository[T any] struct {
    DB    *sqlx.DB
    Table string
}

// GetByID générique
func (r *SQLRepository[T]) GetByID(id int) (*T, error) {
    var t T
    query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", r.Table)
    err := r.DB.Get(&t, query, id)
    if err != nil {
        return nil, err
    }
    return &t, nil
}

// Create générique avec reflection + squirrel
func (r *SQLRepository[T]) Create(entity *T) error {
    val := reflect.Indirect(reflect.ValueOf(entity))
    typ := val.Type()

    values := map[string]interface{}{}
    for i := 0; i < typ.NumField(); i++ {
        field := typ.Field(i)
        dbTag := field.Tag.Get("db")
        if dbTag == "" || dbTag == "-" || dbTag == "id" {
            continue
        }
        values[dbTag] = val.Field(i).Interface()
    }

    query, args, err := sq.
        Insert(r.Table).
        SetMap(values).
        PlaceholderFormat(sq.Dollar).
        ToSql()
    if err != nil {
        return err
    }

    _, err = r.DB.Exec(query, args...)
    return err
}
```

---

## 📁 `main.go` – REST + DB

```go
package main

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"

    "yourmodule/repository"
)

type User struct {
    ID   int    `json:"id" db:"id"`
    Name string `json:"name" db:"name"`
}

var userRepo *repository.SQLRepository[User]

func main() {
    db, err := sqlx.Connect("postgres", "user=postgres password=postgres dbname=testdb sslmode=disable")
    if err != nil {
        panic(err)
    }

    userRepo = &repository.SQLRepository[User]{
        DB:    db,
        Table: "users",
    }

    r := gin.Default()
    r.POST("/users", createUser)
    r.GET("/users/:id", getUserByID)

    r.Run(":8080")
}

func createUser(c *gin.Context) {
    var u User
    if err := c.ShouldBindJSON(&u); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    err := userRepo.Create(&u)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, u)
}

func getUserByID(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    u, err := userRepo.GetByID(id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
        return
    }
    c.JSON(http.StatusOK, u)
}
```

---

## ✅ Résultat : API REST connectée à PostgreSQL avec génériques

* 🔁 **Réutilisable** avec d'autres entités (il suffit de changer la struct et le nom de table)
* 🔐 **Sécurisé** via `sqlx` + `squirrel` (pas d'injection)
* ⚡️ **Performant** et extensible (ajout facile de `Update`, `Delete`, `List`, etc.)

---

---------------------------------------------------------------------------------------------------------------------------

Souhaites-tu que je t’ajoute :

* 🔄 `Update()` et `Delete()` génériques ?
* 🌐 une route `GET /users` avec **pagination SQL** ?
* ✅ des tests API automatisés avec `httptest` ?

Tu choisis !

Parfait, allons-y étape par étape. On va compléter l’API avec :

---

### ✅ Étape 1 : `Update()` et `Delete()` **génériques**

### 🌐 Étape 2 : `GET /users` avec **pagination SQL**

### 🧪 Étape 3 : Tests API automatiques avec `httptest`

---

## ✅ ÉTAPE 1 : `Update()` et `Delete()` génériques

### 🔧 Ajoute dans `sql_repository.go` :

```go
func (r *SQLRepository[T]) Update(id int, entity *T) error {
    val := reflect.Indirect(reflect.ValueOf(entity))
    typ := val.Type()

    values := map[string]interface{}{}
    for i := 0; i < typ.NumField(); i++ {
        field := typ.Field(i)
        dbTag := field.Tag.Get("db")
        if dbTag == "" || dbTag == "-" || dbTag == "id" {
            continue
        }
        values[dbTag] = val.Field(i).Interface()
    }

    query, args, err := sq.
        Update(r.Table).
        SetMap(values).
        Where(sq.Eq{"id": id}).
        PlaceholderFormat(sq.Dollar).
        ToSql()
    if err != nil {
        return err
    }

    _, err = r.DB.Exec(query, args...)
    return err
}

func (r *SQLRepository[T]) Delete(id int) error {
    query, args, err := sq.
        Delete(r.Table).
        Where(sq.Eq{"id": id}).
        PlaceholderFormat(sq.Dollar).
        ToSql()
    if err != nil {
        return err
    }

    _, err = r.DB.Exec(query, args...)
    return err
}
```

---

## ✏️ Dans `main.go`, ajoute ces routes :

```go
r.PUT("/users/:id", updateUser)
r.DELETE("/users/:id", deleteUser)
```

Et les handlers :

```go
func updateUser(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    var u User
    if err := c.ShouldBindJSON(&u); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    err := userRepo.Update(id, &u)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

func deleteUser(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    err := userRepo.Delete(id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
```

---

## 🌐 ÉTAPE 2 : `GET /users` avec **pagination SQL**

Dans `sql_repository.go` :

```go
func (r *SQLRepository[T]) ListPaginated(limit, offset int) ([]T, error) {
    var items []T
    query := fmt.Sprintf("SELECT * FROM %s ORDER BY id LIMIT $1 OFFSET $2", r.Table)
    err := r.DB.Select(&items, query, limit, offset)
    return items, err
}
```

Et dans `main.go` :

```go
r.GET("/users", listUsers)

func listUsers(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
    offset := (page - 1) * size

    users, err := userRepo.ListPaginated(size, offset)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, users)
}
```

---

## 🧪 ÉTAPE 3 : Tests API avec `httptest`

Crée un fichier `main_test.go` :

```go
package main

import (
    "bytes"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
    r := gin.Default()
    r.POST("/users", createUser)
    r.GET("/users/:id", getUserByID)
    return r
}

func TestCreateAndGetUser(t *testing.T) {
    r := setupRouter()

    // Création
    body := []byte(`{"name": "Test User"}`)
    req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)
    assert.Equal(t, 201, w.Code)

    // Récupération (naïf ici, ID supposé connu)
    req = httptest.NewRequest("GET", "/users/1", nil)
    w = httptest.NewRecorder()
    r.ServeHTTP(w, req)
    assert.Equal(t, 200, w.Code)
}
```

---

## ✅ Résumé des fonctionnalités

| Endpoint            | Action REST                  |
| ------------------- | ---------------------------- |
| `POST /users`       | Créer un utilisateur         |
| `GET /users/:id`    | Récupérer un utilisateur     |
| `GET /users`        | Lister avec pagination       |
| `PUT /users/:id`    | Mettre à jour un utilisateur |
| `DELETE /users/:id` | Supprimer un utilisateur     |

---

---------------------------------------------------------------------------------------------------------------------------

Souhaites-tu maintenant :

* intégrer un système de **validation (`binding:"required"`)** avec Gin ?
* ajouter un **middleware d’authentification simple (ex: token Bearer)** ?
* ou packager tout ça en **module structuré (MVC-like)** pour production ?

Parfait ! Passons maintenant à l'intégration de :

---

## ✅ Validation des entrées (via `binding:"required"`)

## 🔐 Authentification simple avec token Bearer

## 🧱 Structure de projet modulaire (type MVC)

---

### ✅ 1. Validation des champs

Gin supporte les tags `binding` via [go-playground/validator](https://github.com/go-playground/validator).

#### 🔧 Ajoute les tags à `User` :

```go
type User struct {
    ID   int    `json:"id" db:"id"`
    Name string `json:"name" binding:"required" db:"name"`
}
```

#### 🔧 Gin le valide automatiquement dans `createUser()` :

```go
if err := c.ShouldBindJSON(&u); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
}
```

---

### 🔐 2. Authentification avec un token Bearer

#### 🔧 Middleware :

```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        auth := c.GetHeader("Authorization")
        if auth != "Bearer secret-token" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
            return
        }
        c.Next()
    }
}
```

#### 🔧 Applique-le aux routes :

```go
auth := r.Group("/")
auth.Use(AuthMiddleware())
auth.POST("/users", createUser)
auth.PUT("/users/:id", updateUser)
auth.DELETE("/users/:id", deleteUser)
```

---

### 🧱 3. Structure MVC recommandée

```bash
yourapp/
│
├── main.go                  # Entry point
├── handler/                 # Gin handlers (controllers)
│   └── user.go
├── model/                   # Domain models (User, Product, etc.)
│   └── user.go
├── repository/              # Generic + entity-specific repos
│   └── sql_repository.go
├── middleware/              # Auth, CORS, logging
│   └── auth.go
├── config/                  # DB, env, etc.
│   └── database.go
└── test/                    # API/unit tests
```

---

Souhaites-tu que je te génère cette **arborescence complète prête à copier**, avec du code dans chaque fichier clé ?
