Le package [`net/http`](https://pkg.go.dev/net/http) de la bibliothèque standard de Go est l’un des plus utilisés pour créer des **serveurs HTTP** (API, sites web, etc.) et pour faire des **requêtes HTTP** en tant que client.

---

## 🔧 Fonctions principales de `net/http`

### 1. **Créer un serveur HTTP (côté serveur)**

```go
package main

import (
    "fmt"
    "net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Bonjour depuis Go !")
}

func main() {
    http.HandleFunc("/", helloHandler)            // Associe la route "/" au handler
    fmt.Println("Serveur sur : http://localhost:8080")
    http.ListenAndServe(":8080", nil)             // Lance le serveur sur le port 8080
}
```

### Explication :

* `http.HandleFunc`: lie un chemin à une fonction.
* `http.ResponseWriter`: permet d’écrire la réponse HTTP (corps, code, headers…).
* `*http.Request`: contient toutes les infos de la requête entrante (URL, méthode, headers, corps, etc.).
* `http.ListenAndServe`: démarre un serveur HTTP.

---

### 2. **Faire des requêtes HTTP (client)**

```go
package main

import (
    "fmt"
    "io"
    "net/http"
)

func main() {
    resp, err := http.Get("https://api.github.com")
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    body, _ := io.ReadAll(resp.Body)
    fmt.Println(string(body))
}
```

### Explication :

* `http.Get`, `http.Post`, `http.NewRequest` : fonctions pour envoyer des requêtes.
* `http.Client` : objet plus avancé pour gérer les timeouts, redirections, etc.
* `resp.Body` : contient la réponse (il faut toujours faire `defer resp.Body.Close()`).

---

## 🧱 Structures importantes

### `http.Request`

Contient :

* `Method` (GET, POST…)
* `URL`
* `Header`
* `Body`
* `Form` / `PostForm` (si `.ParseForm()` est appelée)

### `http.ResponseWriter`

Interface pour écrire :

* Statut avec `WriteHeader(…)`
* Headers avec `Header().Set(...)`
* Corps avec `Write([]byte)` ou `fmt.Fprintln(...)`

---

## 🧩 Exemple avec plusieurs routes

```go
func hello(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Hello!")
}

func goodbye(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Goodbye!")
}

func main() {
    http.HandleFunc("/hello", hello)
    http.HandleFunc("/goodbye", goodbye)
    http.ListenAndServe(":8080", nil)
}
```

---

## 🚨 Limitations de `net/http`

* Il est **basique** : pour des projets plus complexes (routing dynamique, middlewares, validations), on préfère souvent utiliser des frameworks comme **Gin**, **Echo**, **Fiber**, etc.
* Il n'y a pas de notion de groupe de routes, de middleware natif (sauf via `http.Handler` et `http.HandlerFunc`), ni de parsing JSON automatique.

---

### exmpale d'API Rest

Voici un **exemple simple d'API REST** en Go pur avec `net/http` (sans framework externe), qui gère une ressource `Book` avec des routes basiques de type CRUD (Create, Read, Update, Delete).

---

## 📚 Exemple d’API REST avec `net/http`

### 1. **Structure du fichier `main.go`**

```go
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

var books = []Book{
	{ID: 1, Title: "1984", Author: "George Orwell"},
	{ID: 2, Title: "Le Petit Prince", Author: "Antoine de Saint-Exupéry"},
}

// GET /books
func getBooks(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(books)
}

// GET /books/{id}
func getBook(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/books/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}
	for _, book := range books {
		if book.ID == id {
			json.NewEncoder(w).Encode(book)
			return
		}
	}
	http.NotFound(w, r)
}

// POST /books
func createBook(w http.ResponseWriter, r *http.Request) {
	var newBook Book
	if err := json.NewDecoder(r.Body).Decode(&newBook); err != nil {
		http.Error(w, "JSON invalide", http.StatusBadRequest)
		return
	}
	newBook.ID = books[len(books)-1].ID + 1
	books = append(books, newBook)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newBook)
}

// DELETE /books/{id}
func deleteBook(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/books/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}
	for i, book := range books {
		if book.ID == id {
			books = append(books[:i], books[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.NotFound(w, r)
}

func main() {
	http.HandleFunc("/books", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getBooks(w, r)
		case http.MethodPost:
			createBook(w, r)
		default:
			http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/books/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getBook(w, r)
		case http.MethodDelete:
			deleteBook(w, r)
		default:
			http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("API REST démarrée sur : http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
```

---

## ✅ Fonctionnalités

* `GET /books` : liste tous les livres.
* `GET /books/{id}` : récupère un livre par son ID.
* `POST /books` : crée un nouveau livre (ex. via `curl` ou Postman).
* `DELETE /books/{id}` : supprime un livre.

---

## 📌 Exemple de requête POST

```bash
curl -X POST http://localhost:8080/books \
  -H "Content-Type: application/json" \
  -d '{"title": "Clean Code", "author": "Robert C. Martin"}'
```

---

