Très bien ! Voici un **exemple de benchmark** appliqué à **une logique métier plus réaliste**, typique dans une API (par exemple un calcul de statistiques sur des données utilisateur).

---

### 📦 Contexte métier : calculer une moyenne d’âges

#### `services/stats_service.go`

```go
package services

type User struct {
	ID   int
	Name string
	Age  int
}

type StatsService struct{}

func (s *StatsService) AverageAge(users []User) float64 {
	if len(users) == 0 {
		return 0
	}
	total := 0
	for _, u := range users {
		total += u.Age
	}
	return float64(total) / float64(len(users))
}
```

---

### 🧪 Benchmark de cette fonction

#### `services/stats_service_test.go`

```go
package services

import (
	"math/rand"
	"testing"
)

func generateUsers(n int) []User {
	users := make([]User, n)
	for i := 0; i < n; i++ {
		users[i] = User{
			ID:   i + 1,
			Name: "User",
			Age:  rand.Intn(100),
		}
	}
	return users
}

func BenchmarkAverageAge100(b *testing.B) {
	service := StatsService{}
	users := generateUsers(100)

	for i := 0; i < b.N; i++ {
		_ = service.AverageAge(users)
	}
}

func BenchmarkAverageAge10000(b *testing.B) {
	service := StatsService{}
	users := generateUsers(10000)

	for i := 0; i < b.N; i++ {
		_ = service.AverageAge(users)
	}
}
```

---

### ▶️ Lancer le benchmark

```bash
go test -bench=.
```

---

### ✅ Résultat typique

```
BenchmarkAverageAge100-8         1150000	       1055 ns/op
BenchmarkAverageAge10000-8        55000	      21783 ns/op
```

---

### 🔬 Interprétation

* Le temps augmente avec la taille du tableau d’utilisateurs.
* Tu peux ainsi profiler des algorithmes ou valider des optimisations dans des services métiers.

---

Souhaites-tu maintenant un **benchmark HTTP** (ex. : test de la vitesse de réponse d’un endpoint Gin) ou un benchmark avec **accès base de données simulé** ?

Super ! Voici un **exemple de benchmark d’un endpoint HTTP** utilisant **Gin**. L’objectif est de tester la **performance d’un handler** simulant une logique métier (ex. : `/ping` ou `/users/average-age`), sans base de données.

---

## 🌐 Benchmark HTTP avec Gin

### 1. 🎯 Exemple de handler simple

#### `handlers/ping_handler.go`

```go
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func PingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}
```

---

### 2. 🧪 Benchmark du handler

#### `handlers/ping_handler_test.go`

```go
package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func BenchmarkPingHandler(b *testing.B) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/ping", PingHandler)

	req, _ := http.NewRequest("GET", "/ping", nil)

	b.ResetTimer() // Ne compte pas l'initialisation
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			b.Fatalf("unexpected status: %d", w.Code)
		}
	}
}
```

---

### ▶️ Lancer le benchmark

```bash
go test -bench=.
```

---

### ✅ Résultat attendu

```
BenchmarkPingHandler-8     500000	      2700 ns/op
```

Cela indique que ton endpoint `/ping` prend en moyenne 2.7 µs par requête dans un environnement simulé.

---

## 🧠 Astuce

Si tu veux benchmarker un vrai service métier dans un handler (comme `GET /users/average-age`), injecte un mock du service et applique la même technique avec `httptest`.

---

Souhaites-tu que je t’ajoute un benchmark complet avec un **mock de service métier dans un handler HTTP** (ex. : `/users/average-age`) ?
