En Go (Golang), array, slice, et map sont trois types fondamentaux pour stocker des collections de données. Voici une explication claire avec exemples :

1. Array (Tableau)


Un array est une structure avec une taille fixe et un type homogène.

```go
var a [3]int // tableau de 3 entiers
a[0] = 10
a[1] = 20
a[2] = 30
fmt.Println(a) // [10 20 30]
```

Taille fixe (ne peut pas être redimensionné).

Type unique pour tous les éléments.

On peut itérer avec for.

2. Slice


Un slice est une vue dynamique sur un tableau : il peut changer de taille.

```go
s := []int{1, 2, 3} // slice de 3 éléments
s = append(s, 4)    // on ajoute un élément
fmt.Println(s)      // [1 2 3 4]
```

Taille dynamique.

Très utilisé en Go.

On peut trancher (slice[1:3]) et ajouter (append).

C’est en fait un pointeur vers un tableau sous-jacent.

3. Map


Une map est une structure clé-valeur, comme un dictionnaire.

```go
m := make(map[string]int)
m["age"] = 30
m["year"] = 2025
fmt.Println(m["age"]) // 30
```

Clé unique, de n’importe quel type comparable (ex: string, int).

Valeur de n’importe quel type.

Accès rapide par clé.

Exemple combiné

```go
func main() {
    // Array
    var arr = [3]int{10, 20, 30}

    // Slice
    sl := []int{1, 2, 3}
    sl = append(sl, 4)

    // Map
    m := map[string]string{
        "name": "Alice",
        "city": "Paris",
    }

    fmt.Println("Array:", arr)
    fmt.Println("Slice:", sl)
    fmt.Println("Map:", m)
}
```

Souhaites-tu un exemple plus pratique avec des structures ou en interaction avec une API ?

Example 

Parfait ! Voici un exemple pratique combinant slice et map dans une petite API REST en Go avec Gin, simulant une liste d’utilisateurs.

Exemple : Mini API d’utilisateurs (en mémoire)

```go
package main

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

// Structure User
type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

// Slice pour stocker les utilisateurs
var users = []User{
    {ID: 1, Name: "Alice", Email: "alice@example.com"},
    {ID: 2, Name: "Bob", Email: "bob@example.com"},
}

// Map pour indexer les utilisateurs par ID
var userIndex = map[int]User{
    1: {ID: 1, Name: "Alice", Email: "alice@example.com"},
    2: {ID: 2, Name: "Bob", Email: "bob@example.com"},
}

func main() {
    r := gin.Default()

    // GET /users - retourne tous les utilisateurs
    r.GET("/users", func(c *gin.Context) {
        c.JSON(http.StatusOK, users) // Utilise le slice
    })

    // GET /users/:id - retourne un utilisateur via map
    r.GET("/users/:id", func(c *gin.Context) {
        idParam := c.Param("id")
        var id int
        fmt.Sscanf(idParam, "%d", &id)

        user, exists := userIndex[id]
        if !exists {
            c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
            return
        }

        c.JSON(http.StatusOK, user)
    })

    // POST /users - ajoute un utilisateur (dans le slice + map)
    r.POST("/users", func(c *gin.Context) {
        var newUser User
        if err := c.ShouldBindJSON(&newUser); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        newUser.ID = len(users) + 1
        users = append(users, newUser)
        userIndex[newUser.ID] = newUser

        c.JSON(http.StatusCreated, newUser)
    })

    r.Run() // par défaut : localhost:8080
}
```

À retenir :
    * slice (users) est utilisé pour stocker la liste ordonnée.

    * map (userIndex) permet un accès rapide par ID.

    * L’API expose trois routes simples : /users, /users/:id, et un POST pour ajouter.

----------------------------------------------------------------------------------------------------------------------------------

Bien sûr ! Voici des tests unitaires et d’intégration basiques pour l’API REST d’utilisateurs en mémoire, en utilisant le package net/http/httptest et Gin.

Fichier : main_test.go

```go
package main

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
)

// Helper pour configurer le routeur pour les tests
func setupRouter() *gin.Engine {
    r := gin.Default()

    r.GET("/users", func(c *gin.Context) {
        c.JSON(http.StatusOK, users)
    })

    r.GET("/users/:id", func(c *gin.Context) {
        var id int
        fmt.Sscanf(c.Param("id"), "%d", &id)
        user, exists := userIndex[id]
        if !exists {
            c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
            return
        }
        c.JSON(http.StatusOK, user)
    })

    r.POST("/users", func(c *gin.Context) {
        var newUser User
        if err := c.ShouldBindJSON(&newUser); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        newUser.ID = len(users) + 1
        users = append(users, newUser)
        userIndex[newUser.ID] = newUser
        c.JSON(http.StatusCreated, newUser)
    })

    return r
}
```

Tests :

```go
func TestGetUsers(t *testing.T) {
    router := setupRouter()
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/users", nil)
    router.ServeHTTP(w, req)

    assert.Equal(t, 200, w.Code)
    var response []User
    json.Unmarshal(w.Body.Bytes(), &response)
    assert.True(t, len(response) >= 2)
}

func TestGetUserByID_Success(t *testing.T) {
    router := setupRouter()
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/users/1", nil)
    router.ServeHTTP(w, req)

    assert.Equal(t, 200, w.Code)
    var user User
    json.Unmarshal(w.Body.Bytes(), &user)
    assert.Equal(t, "Alice", user.Name)
}

func TestGetUserByID_NotFound(t *testing.T) {
    router := setupRouter()
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/users/999", nil)
    router.ServeHTTP(w, req)

    assert.Equal(t, 404, w.Code)
}

func TestCreateUser(t *testing.T) {
    router := setupRouter()

    newUser := User{Name: "Charlie", Email: "charlie@example.com"}
    jsonValue, _ := json.Marshal(newUser)

    w := httptest.NewRecorder()
    req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonValue))
    req.Header.Set("Content-Type", "application/json")
    router.ServeHTTP(w, req)

    assert.Equal(t, 201, w.Code)

    var createdUser User
    json.Unmarshal(w.Body.Bytes(), &createdUser)
    assert.Equal(t, "Charlie", createdUser.Name)
    assert.NotZero(t, createdUser.ID)
}
```

Dépendances :


Installe testify pour des assertions simples :

`go get github.com/stretchr/testify`

Lancer les tests :

`go test -v`

Souhaites-tu que je structure tout cela dans une arborescence propre (MVC) ou que je prépare les tests avec une base SQLite ?



Avec architecture



Parfait ! Voici une structure propre de type MVC avec SQLite, avec des tests unitaires, basée sur Gin.

Arborescence proposée :

```shell
/userapi
├── main.go
├── go.mod
├── config/
│   └── db.go
├── models/
│   └── user.go
├── controllers/
│   └── user_controller.go
├── routes/
│   └── routes.go
├── tests/
│   └── user_test.go
```

1. config/db.go : Connexion SQLite

```go
package config

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
    "log"
)

var DB *sql.DB

func ConnectDatabase() {
    var err error
    DB, err = sql.Open("sqlite3", "users.db")
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    createTable := `CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT,
        email TEXT
    );`
    _, err = DB.Exec(createTable)
    if err != nil {
        log.Fatal("Failed to create table:", err)
    }
}
```

2. models/user.go

```go
package models

type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}
```

3. controllers/user_controller.go

```go
package controllers

import (
    "database/sql"
    "net/http"
    "userapi/config"
    "userapi/models"

    "github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
    rows, err := config.DB.Query("SELECT id, name, email FROM users")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer rows.Close()

    var users []models.User
    for rows.Next() {
        var user models.User
        rows.Scan(&user.ID, &user.Name, &user.Email)
        users = append(users, user)
    }

    c.JSON(http.StatusOK, users)
}

func GetUser(c *gin.Context) {
    id := c.Param("id")
    row := config.DB.QueryRow("SELECT id, name, email FROM users WHERE id = ?", id)

    var user models.User
    err := row.Scan(&user.ID, &user.Name, &user.Email)
    if err == sql.ErrNoRows {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    } else if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, user)
}

func CreateUser(c *gin.Context) {
    var user models.User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    res, err := config.DB.Exec("INSERT INTO users (name, email) VALUES (?, ?)", user.Name, user.Email)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    id, _ := res.LastInsertId()
    user.ID = int(id)
    c.JSON(http.StatusCreated, user)
}
```
4. routes/routes.go

```go
package routes

import (
    "github.com/gin-gonic/gin"
    "userapi/controllers"
)

func SetupRoutes() *gin.Engine {
    r := gin.Default()

    r.GET("/users", controllers.GetUsers)
    r.GET("/users/:id", controllers.GetUser)
    r.POST("/users", controllers.CreateUser)

    return r
}
````

5. main.go

```go
package main

import (
    "userapi/config"
    "userapi/routes"
)

func main() {
    config.ConnectDatabase()
    r := routes.SetupRoutes()
    r.Run(":8080")
}
```

6. tests/user_test.go (test simple d’intégration)

```go
package tests

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "userapi/config"
    "userapi/models"
    "userapi/routes"

    "github.com/stretchr/testify/assert"
)

func init() {
    config.ConnectDatabase()
}

func TestCreateAndGetUser(t *testing.T) {
    router := routes.SetupRoutes()

    // Créer un utilisateur
    user := models.User{Name: "Test User", Email: "test@example.com"}
    jsonData, _ := json.Marshal(user)
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonData))
    req.Header.Set("Content-Type", "application/json")
    router.ServeHTTP(w, req)

    assert.Equal(t, 201, w.Code)

    var created models.User
    json.Unmarshal(w.Body.Bytes(), &created)
    assert.NotZero(t, created.ID)

    // Récupérer l'utilisateur créé
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("GET", "/users/"+string(rune(created.ID)), nil)
    router.ServeHTTP(w, req)

    assert.Equal(t, 200, w.Code)
}
```

Souhaites-tu que je génère un repository Git complet, ou que je t’aide à ajouter l’authentification ou des validations dans ce projet ?
