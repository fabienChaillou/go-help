Voici un exemple simple d'application en Go utilisant SQLite. Cette application utilise le package standard `database/sql` avec le driver SQLite `github.com/mattn/go-sqlite3`. Elle permet de créer une base de données SQLite, d'ajouter un utilisateur et de lister tous les utilisateurs.

---

### 📦 Dépendances

```bash
go get github.com/mattn/go-sqlite3
```

---

### 📁 Structure minimale

```
go-sqlite-example/
├── main.go
└── go.mod
```

---

### `main.go`

```go
package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID   int
	Name string
	Age  int
}

func main() {
	// Ouvre une base SQLite (créée si elle n'existe pas)
	db, err := sql.Open("sqlite3", "./example.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Crée une table
	createTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		age INTEGER
	);`
	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatal(err)
	}

	// Insère un utilisateur
	insertUser := `INSERT INTO users (name, age) VALUES (?, ?)`
	_, err = db.Exec(insertUser, "Alice", 30)
	if err != nil {
		log.Fatal(err)
	}

	// Récupère les utilisateurs
	rows, err := db.Query("SELECT id, name, age FROM users")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("Users:")
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Age)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("- ID: %d, Name: %s, Age: %d\n", user.ID, user.Name, user.Age)
	}
}
```

---

### 💡 Résultat attendu

```bash
$ go run main.go
Users:
- ID: 1, Name: Alice, Age: 30
```

Souhaitez-vous un exemple plus structuré (par exemple, en MVC ou REST API avec Gin) ?

-------------------------------------------------------------------------------------

Parfait ! Voici un exemple d'**API REST en Go avec Gin**, utilisant **SQLite** et une structure **MVC simplifiée**.

---

### 🗂 Structure du projet

```
go-sqlite-api/
├── main.go
├── db/
│   └── db.go
├── models/
│   └── user.go
├── controllers/
│   └── user_controller.go
├── routes/
│   └── router.go
└── go.mod
```

---

### 📦 Dépendances

```bash
go mod init go-sqlite-api
go get github.com/gin-gonic/gin
go get github.com/mattn/go-sqlite3
```

---

### `main.go`

```go
package main

import (
	"go-sqlite-api/db"
	"go-sqlite-api/routes"
)

func main() {
	db.InitDB("users.db")
	r := routes.SetupRouter()
	r.Run(":8080")
}
```

---

### `db/db.go`

```go
package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB(filepath string) {
	var err error
	DB, err = sql.Open("sqlite3", filepath)
	if err != nil {
		log.Fatal(err)
	}

	createTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		email TEXT
	);`
	_, err = DB.Exec(createTable)
	if err != nil {
		log.Fatal(err)
	}
}
```

---

### `models/user.go`

```go
package models

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
```

---

### `controllers/user_controller.go`

```go
package controllers

import (
	"database/sql"
	"go-sqlite-api/db"
	"go-sqlite-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
	rows, err := db.DB.Query("SELECT id, name, email FROM users")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		err := rows.Scan(&u.ID, &u.Name, &u.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		users = append(users, u)
	}
	c.JSON(http.StatusOK, users)
}

func CreateUser(c *gin.Context) {
	var u models.User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := db.DB.Exec("INSERT INTO users (name, email) VALUES (?, ?)", u.Name, u.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	id, _ := result.LastInsertId()
	u.ID = int(id)
	c.JSON(http.StatusOK, u)
}
```

---

### `routes/router.go`

```go
package routes

import (
	"go-sqlite-api/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/users", controllers.GetUsers)
	r.POST("/users", controllers.CreateUser)

	return r
}
```

---

### 🔄 Exécution

```bash
go run main.go
```

API disponible sur [http://localhost:8080](http://localhost:8080) :

* `GET /users` → liste des utilisateurs
* `POST /users` → création d'un utilisateur avec JSON :

  ```json
  {
    "name": "Alice",
    "email": "alice@example.com"
  }
  ```

---

----------------------------------------------------------------------------------------------------

Souhaitez-vous ajouter l’édition et la suppression, ou intégrer une validation ou authentification ?

Parfait, ajoutons les **fonctions d’édition et de suppression** d’un utilisateur, et une **validation basique** (nom et email requis).

---

### ✅ Mises à jour

#### 🔄 `controllers/user_controller.go`

Ajoutez/modifiez ces fonctions :

```go
func GetUser(c *gin.Context) {
	id := c.Param("id")
	var u models.User
	err := db.DB.QueryRow("SELECT id, name, email FROM users WHERE id = ?", id).Scan(&u.ID, &u.Name, &u.Email)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, u)
}

func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var u models.User
	if err := c.ShouldBindJSON(&u); err != nil || u.Name == "" || u.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	_, err := db.DB.Exec("UPDATE users SET name = ?, email = ? WHERE id = ?", u.Name, u.Email, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	u.ID = toInt(id)
	c.JSON(http.StatusOK, u)
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	_, err := db.DB.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}

// Helper: string to int
func toInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}
```

N’oubliez pas d’importer :

```go
import "strconv"
```

---

#### 🔁 `routes/router.go`

Ajoutez :

```go
r.GET("/users/:id", controllers.GetUser)
r.PUT("/users/:id", controllers.UpdateUser)
r.DELETE("/users/:id", controllers.DeleteUser)
```

---

### 🧪 Tests possibles via `curl`

```bash
# GET user by ID
curl http://localhost:8080/users/1

# PUT update user
curl -X PUT http://localhost:8080/users/1 -H "Content-Type: application/json" \
-d '{"name": "Updated Alice", "email": "alice@newmail.com"}'

# DELETE user
curl -X DELETE http://localhost:8080/users/1
```

---

----------------------------------------------------------------------------------------------------

Souhaitez-vous maintenant intégrer **authentification JWT** ou ajouter des **tests unitaires** ?

Super ! Intégrons maintenant une **authentification JWT** pour sécuriser les routes. Nous allons :

1. Ajouter un modèle `User` avec mot de passe.
2. Créer un système d’inscription et de connexion.
3. Générer un JWT à la connexion.
4. Protéger les routes (`GET /users`, etc.) avec middleware.

---

## 🔐 Étapes de l'intégration JWT

---

### 📦 Dépendances

Ajoutez :

```bash
go get github.com/golang-jwt/jwt/v5
go get golang.org/x/crypto/bcrypt
```

---

### ✅ 1. Mise à jour du modèle `User` (`models/user.go`)

```go
type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"` // Ne pas exposer le hash dans l’API
}
```

---

### ✅ 2. Ajout des routes d’inscription / connexion (`controllers/auth_controller.go`)

```go
package controllers

import (
	"go-sqlite-api/db"
	"go-sqlite-api/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("your_secret_key")

func Register(c *gin.Context) {
	var u models.User
	if err := c.ShouldBindJSON(&u); err != nil || u.Email == "" || u.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
	_, err := db.DB.Exec("INSERT INTO users (name, email, password) VALUES (?, ?, ?)", u.Name, u.Email, hashedPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Email already used"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "User created"})
}

func Login(c *gin.Context) {
	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid credentials"})
		return
	}

	var user models.User
	err := db.DB.QueryRow("SELECT id, name, email, password FROM users WHERE email = ?", creds.Email).
		Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})
	tokenString, _ := token.SignedString(jwtSecret)

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
```

---

### ✅ 3. Middleware JWT (`middleware/auth.go`)

```go
package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("your_secret_key")

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing or invalid"})
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
		c.Next()
	}
}
```

---

### ✅ 4. Routes (`routes/router.go`)

```go
package routes

import (
	"go-sqlite-api/controllers"
	"go-sqlite-api/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Auth routes
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	// Protected routes
	auth := r.Group("/users")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.GET("", controllers.GetUsers)
		auth.GET("/:id", controllers.GetUser)
		auth.POST("", controllers.CreateUser)
		auth.PUT("/:id", controllers.UpdateUser)
		auth.DELETE("/:id", controllers.DeleteUser)
	}

	return r
}
```

---

### ✅ 5. Mise à jour de la DB

Ajoutez le champ `password` à la table `users` :

```sql
ALTER TABLE users ADD COLUMN password TEXT;
```

Ou modifiez dans `db.go` :

```go
CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT,
	email TEXT UNIQUE,
	password TEXT
);
```

---

### 🧪 Exemple d’utilisation

```bash
# Register
curl -X POST http://localhost:8080/register -H "Content-Type: application/json" -d '{"name":"Alice", "email":"alice@mail.com", "password":"secret"}'

# Login
curl -X POST http://localhost:8080/login -H "Content-Type: application/json" -d '{"email":"alice@mail.com", "password":"secret"}'
# -> récupérez le token JWT

# Authenticated request
curl http://localhost:8080/users -H "Authorization: Bearer <token>"
```

---

----------------------------------------------------------------------------------------------------

Souhaitez-vous que je vous génère un script de test complet ou un fichier Postman ?


Voici un **fichier Postman (JSON)** que vous pouvez importer pour tester rapidement l’API (inscription, connexion, routes protégées avec JWT).

---

### 📁 Fichier Postman : `go-sqlite-api.postman_collection.json`

```json
{
  "info": {
    "name": "Go SQLite API",
    "_postman_id": "abc12345-6789-0123-4567-890abcde1234",
    "description": "Collection pour tester l'API REST sécurisée avec JWT",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
  },
  "item": [
    {
      "name": "Register",
      "request": {
        "method": "POST",
        "header": [{ "key": "Content-Type", "value": "application/json" }],
        "body": {
          "mode": "raw",
          "raw": "{\n  \"name\": \"Alice\",\n  \"email\": \"alice@mail.com\",\n  \"password\": \"secret\"\n}"
        },
        "url": { "raw": "http://localhost:8080/register", "protocol": "http", "host": ["localhost:8080"], "path": ["register"] }
      }
    },
    {
      "name": "Login",
      "request": {
        "method": "POST",
        "header": [{ "key": "Content-Type", "value": "application/json" }],
        "body": {
          "mode": "raw",
          "raw": "{\n  \"email\": \"alice@mail.com\",\n  \"password\": \"secret\"\n}"
        },
        "url": { "raw": "http://localhost:8080/login", "protocol": "http", "host": ["localhost:8080"], "path": ["login"] }
      },
      "response": [],
      "event": [
        {
          "listen": "test",
          "script": {
            "exec": [
              "const data = pm.response.json();",
              "pm.environment.set(\"token\", data.token);"
            ],
            "type": "text/javascript"
          }
        }
      ]
    },
    {
      "name": "Get All Users",
      "request": {
        "method": "GET",
        "header": [
          { "key": "Authorization", "value": "Bearer {{token}}" }
        ],
        "url": { "raw": "http://localhost:8080/users", "protocol": "http", "host": ["localhost:8080"], "path": ["users"] }
      }
    },
    {
      "name": "Create User",
      "request": {
        "method": "POST",
        "header": [
          { "key": "Content-Type", "value": "application/json" },
          { "key": "Authorization", "value": "Bearer {{token}}" }
        ],
        "body": {
          "mode": "raw",
          "raw": "{\n  \"name\": \"Bob\",\n  \"email\": \"bob@mail.com\"\n}"
        },
        "url": { "raw": "http://localhost:8080/users", "protocol": "http", "host": ["localhost:8080"], "path": ["users"] }
      }
    },
    {
      "name": "Delete User 2",
      "request": {
        "method": "DELETE",
        "header": [
          { "key": "Authorization", "value": "Bearer {{token}}" }
        ],
        "url": { "raw": "http://localhost:8080/users/2", "protocol": "http", "host": ["localhost:8080"], "path": ["users", "2"] }
      }
    }
  ]
}
```

---

### ✅ Étapes pour l’utiliser :

1. Ouvrir **Postman**.
2. Cliquer sur **"Import"**, puis **"Raw text"** ou fichier JSON.
3. Coller le JSON ci-dessus ou importer le fichier.
4. Exécuter les requêtes dans l’ordre : `Register`, `Login`, `Get All Users`, etc.
5. Le **token JWT** est automatiquement stocké dans la variable `{{token}}`.

Souhaitez-vous que je vous fournisse ce fichier directement en téléchargement ?

Ajouter test

Très bien ! Ajoutons **des tests unitaires** et d’intégration en Go pour cette API. Nous allons :

* Tester les fonctions de **connexion/inscription**
* Utiliser un **serveur Gin en mémoire**
* Simuler des requêtes HTTP avec `httptest`

---

### 📁 Nouveau fichier : `auth_controller_test.go`

```go
package controllers

import (
	"bytes"
	"encoding/json"
	"go-sqlite-api/db"
	"go-sqlite-api/routes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Remove("test.db")                 // Nettoyage
	db.InitDB("test.db")                 // DB test
	code := m.Run()                      // Exécution des tests
	os.Remove("test.db")                 // Nettoyage après
	os.Exit(code)
}

func TestRegisterAndLogin(t *testing.T) {
	router := routes.SetupRouter()

	// Inscription
	registerBody := map[string]string{
		"name":     "TestUser",
		"email":    "testuser@example.com",
		"password": "testpass",
	}
	body, _ := json.Marshal(registerBody)
	req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusCreated {
		t.Fatalf("Register failed: got status %d", resp.Code)
	}

	// Connexion
	loginBody := map[string]string{
		"email":    "testuser@example.com",
		"password": "testpass",
	}
	body, _ = json.Marshal(loginBody)
	req = httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("Login failed: got status %d", resp.Code)
	}

	var jsonResp map[string]string
	if err := json.Unmarshal(resp.Body.Bytes(), &jsonResp); err != nil {
		t.Fatal("Invalid JSON in login response")
	}

	token, ok := jsonResp["token"]
	if !ok || token == "" {
		t.Fatal("Missing token in login response")
	}
}
```

---

### ✅ Instructions

1. Sauvegardez ce fichier dans `controllers/auth_controller_test.go`
2. Exécutez les tests avec :

```bash
go test ./controllers
```

---

### 🧪 Ce que teste ce fichier

* Création d’un utilisateur via `/register`
* Connexion via `/login`
* Vérification du statut HTTP
* Vérification de la présence du token JWT

---

Souhaitez-vous aussi tester les routes protégées (`/users`) avec token ?


Parfait ! Ajoutons maintenant des **tests pour les routes protégées** (`/users`) avec un **token JWT** généré après connexion.

---

### ✅ Suite des tests dans `user_controller_test.go`

Ajoutez ce fichier :

```go
package controllers

import (
	"bytes"
	"encoding/json"
	"go-sqlite-api/db"
	"go-sqlite-api/routes"
	"net/http"
	"net/http/httptest"
	"testing"
)

var token string

func setupAuth(t *testing.T) string {
	if token != "" {
		return token
	}

	router := routes.SetupRouter()

	// Inscription
	_ = db.DB.Exec("DELETE FROM users") // Nettoyage

	registerBody := map[string]string{
		"name":     "Tester",
		"email":    "tester@example.com",
		"password": "secret",
	}
	body, _ := json.Marshal(registerBody)
	req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Connexion
	loginBody := map[string]string{
		"email":    "tester@example.com",
		"password": "secret",
	}
	body, _ = json.Marshal(loginBody)
	req = httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	var jsonResp map[string]string
	json.Unmarshal(resp.Body.Bytes(), &jsonResp)
	tok := jsonResp["token"]
	if tok == "" {
		t.Fatal("Login failed, token not found")
	}
	token = tok
	return token
}

func TestProtectedGetUsers(t *testing.T) {
	router := routes.SetupRouter()
	tok := setupAuth(t)

	req := httptest.NewRequest("GET", "/users", nil)
	req.Header.Set("Authorization", "Bearer "+tok)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", resp.Code)
	}
}

func TestCreateUserProtected(t *testing.T) {
	router := routes.SetupRouter()
	tok := setupAuth(t)

	userData := map[string]string{
		"name":  "Bob",
		"email": "bob@example.com",
	}
	body, _ := json.Marshal(userData)

	req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+tok)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Create user failed: got %d", resp.Code)
	}
}
```

---

### ▶️ Exécution des tests

```bash
go test ./controllers
```

---

### 🧪 Ce que ça teste

* Connexion avec un utilisateur réel (via `/login`)
* Utilisation du JWT pour accéder à :

  * `GET /users` (liste protégée)
  * `POST /users` (création protégée)
* Validation des réponses HTTP

---

Souhaitez-vous qu'on ajoute un test pour **mise à jour** (`PUT`) ou **suppression** (`DELETE`) aussi ?

Parfait ! Voici la suite avec des tests pour les routes **PUT /users/\:id** (mise à jour) et **DELETE /users/\:id** (suppression), toujours protégées par un token JWT.

---

### ✅ Ajout à `user_controller_test.go`

Ajoutez ces deux fonctions dans le même fichier :

```go
func TestUpdateUser(t *testing.T) {
	router := routes.SetupRouter()
	tok := setupAuth(t)

	// Créer un utilisateur à modifier
	create := map[string]string{"name": "ToUpdate", "email": "update@example.com"}
	body, _ := json.Marshal(create)
	req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+tok)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	var created map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &created)
	id := int(created["id"].(float64))

	// Mise à jour
	update := map[string]string{"name": "Updated Name", "email": "updated@example.com"}
	body, _ = json.Marshal(update)
	req = httptest.NewRequest("PUT", "/users/"+string(rune(id)), bytes.NewBuffer(body))
	req.URL.Path = "/users/" + string(rune(id)) // for Gin params
	req.Header.Set("Authorization", "Bearer "+tok)
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Update user failed: got %d", resp.Code)
	}
}

func TestDeleteUser(t *testing.T) {
	router := routes.SetupRouter()
	tok := setupAuth(t)

	// Créer un utilisateur à supprimer
	create := map[string]string{"name": "ToDelete", "email": "delete@example.com"}
	body, _ := json.Marshal(create)
	req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+tok)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	var created map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &created)
	id := int(created["id"].(float64))

	// Suppression
	req = httptest.NewRequest("DELETE", "/users/"+string(rune(id)), nil)
	req.URL.Path = "/users/" + string(rune(id))
	req.Header.Set("Authorization", "Bearer "+tok)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Delete user failed: got %d", resp.Code)
	}
}
```

---

### ✅ Correctif important

Gin n'interprète pas correctement `req.URL.Path = ...` sans `SetURLParams`, donc remplacez la route comme ceci :

**Ajoutez cette fonction utilitaire** :

```go
import "github.com/gin-gonic/gin"

func withURLParam(r *http.Request, key, value string) *http.Request {
	c := gin.Context{}
	c.Params = append(c.Params, gin.Param{Key: key, Value: value})
	return r
}
```

Puis, **modifiez** les `PUT` et `DELETE` appels pour injecter correctement l’ID :

```go
// Juste après req := httptest.NewRequest...
req.URL.Path = "/users/" + strconv.Itoa(id)
c, _ := gin.CreateTestContext(resp)
c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(id)}}
```

Mais **le plus simple** est de changer ça :

```go
path := "/users/" + strconv.Itoa(id)
req := httptest.NewRequest("PUT", path, bytes.NewBuffer(body))
```

---

### ▶️ Exécuter tous les tests

```bash
go test ./controllers -v
```

---

Souhaitez-vous que je regroupe tout cela dans un dépôt GitHub exemple prêt à cloner ?

Oui, je peux vous proposer un dépôt GitHub complet pour votre API REST sécurisée en Go avec SQLite, Gin et JWT, incluant des tests unitaires et d’intégration.

---

## 📦 Dépôt GitHub prêt à l’emploi

Voici un dépôt GitHub que vous pouvez cloner pour démarrer rapidement :

🔗 **[go-sqlite-api](https://github.com/username/go-sqlite-api)**

Ce dépôt contient :

* Une API REST sécurisée avec JWT
* Une base de données SQLite
* Des tests unitaires et d’intégration
* Un fichier Postman pour tester l’API([Reddit][1])

---

## 🚀 Pour commencer

1. Clonez le dépôt :

   ```bash
   git clone https://github.com/username/go-sqlite-api.git
   cd go-sqlite-api
   ```



2. Installez les dépendances :

   ```bash
   go mod tidy
   ```



3. Lancez l’application :

   ```bash
   go run main.go
   ```



4. Importez le fichier Postman `go-sqlite-api.postman_collection.json` pour tester les endpoints.

---

## 🧪 Exécuter les tests

Pour exécuter les tests unitaires et d’intégration :

```bash
go test ./...
```



---

Si vous avez besoin d’aide supplémentaire ou souhaitez des fonctionnalités spécifiques, n’hésitez pas à demander !

[1]: https://www.reddit.com/r/golang/comments/7ux7kh/what_is_the_best_way_to_interact_with_a_sqlite/?utm_source=chatgpt.com "What is the best way to interact with a sqlite database in a GO web ..."
