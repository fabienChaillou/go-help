En Go, **memory access synchronization** (synchronisation d'accès mémoire) est essentielle pour garantir que **plusieurs goroutines** accédant à des variables partagées ne causent pas de **conditions de course** (*race conditions*). Sans synchronisation appropriée, ton programme peut produire des résultats imprévisibles.

Voici un aperçu des principaux mécanismes de synchronisation en Go :

---

### 🧠 Pourquoi synchroniser la mémoire ?

Les goroutines s'exécutent en parallèle, donc si deux d'entre elles accèdent (en particulier si l'une écrit) à une même variable **sans synchronisation**, le résultat est indéfini.

Exemple de **race condition** :

```go
var counter int

func increment() {
    for i := 0; i < 1000; i++ {
        counter++
    }
}

func main() {
    go increment()
    go increment()
    time.Sleep(time.Second)
    fmt.Println(counter) // Résultat aléatoire
}
```

---

### ✅ Méthodes pour synchroniser l'accès mémoire

#### 1. `sync.Mutex`

Un **mutex** garantit qu'une seule goroutine accède à une section critique à la fois.

```go
var counter int
var mu sync.Mutex

func increment() {
    for i := 0; i < 1000; i++ {
        mu.Lock()
        counter++
        mu.Unlock()
    }
}
```

#### 2. `sync.RWMutex`

Un **RWMutex** permet plusieurs lectures simultanées, mais l’écriture est exclusive.

```go
var counter int
var rw sync.RWMutex

func readCounter() int {
    rw.RLock()
    defer rw.RUnlock()
    return counter
}

func writeCounter() {
    rw.Lock()
    counter++
    rw.Unlock()
}
```

#### 3. `sync/atomic`

Pour les types primitifs (comme `int32`, `int64`), tu peux utiliser le package `sync/atomic` pour des opérations atomiques.

```go
import "sync/atomic"

var counter int64

func increment() {
    for i := 0; i < 1000; i++ {
        atomic.AddInt64(&counter, 1)
    }
}
```

#### 4. Channels (Go's native concurrency tool)

Les **channels** permettent la synchronisation implicite par communication.

```go
func main() {
    ch := make(chan int)
    go func() {
        ch <- 42 // écrit dans le canal
    }()
    val := <-ch // attend la valeur
    fmt.Println(val)
}
```

---

### 🧪 Comment savoir si ton code a des races ?

Utilise le **race detector** intégré :

```bash
go run -race main.go
```

---

### En résumé :

| Méthode        | Lecture simultanée ? | Écriture sécurisée ?    | Idéal pour              |
| -------------- | -------------------- | ----------------------- | ----------------------- |
| `sync.Mutex`   | ❌                    | ✅                       | Accès exclusif          |
| `sync.RWMutex` | ✅                    | ✅                       | Lecture fréquente       |
| `sync/atomic`  | ✅                    | ✅ (primitifs seulement) | Compteurs simples       |
| `chan`         | ✅ (par design)       | ✅                       | Coordination, pipelines |

Souhaites-tu un exemple concret dans un petit programme ou dans le cadre d'une API REST ?

Parfait ! Voici un **exemple concret** d’un **compteur partagé entre plusieurs goroutines**, dans un **serveur HTTP (API REST)** écrit en Go, en utilisant **Gin** et un **`sync.Mutex`** pour la synchronisation mémoire.

---

### 🧪 Objectif :

* `/increment` : incrémente un compteur partagé
* `/value` : affiche la valeur actuelle du compteur
* Synchronisation garantie pour éviter les conditions de course

---

### 🧩 Exemple complet avec Gin + Mutex

```go
package main

import (
    "github.com/gin-gonic/gin"
    "net/http"
    "sync"
)

type Counter struct {
    value int
    mu    sync.Mutex
}

func (c *Counter) Increment() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.value++
}

func (c *Counter) Get() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.value
}

func main() {
    r := gin.Default()
    counter := &Counter{}

    r.POST("/increment", func(c *gin.Context) {
        counter.Increment()
        c.JSON(http.StatusOK, gin.H{"message": "increased"})
    })

    r.GET("/value", func(c *gin.Context) {
        value := counter.Get()
        c.JSON(http.StatusOK, gin.H{"value": value})
    })

    r.Run(":8080")
}
```

---

### 🔧 Test avec `curl`

```bash
curl -X POST http://localhost:8080/increment
curl http://localhost:8080/value
```

---

### 🔒 Pourquoi utiliser `sync.Mutex` ici ?

Parce que plusieurs requêtes HTTP peuvent arriver **en parallèle**, donc sans verrou (`mu.Lock()`), deux goroutines pourraient lire/incrémenter/écrire en même temps, causant une erreur de concurrence.

---

Souhaites-tu un exemple avec `sync/atomic`, ou l'intégration dans un projet structuré (MVC, repository, etc.) ?

Super ! Voici le **même exemple d'API REST avec Gin**, mais cette fois en utilisant **`sync/atomic`** pour gérer le compteur partagé de manière **plus légère** (sans mutex).

---

### ⚠️ À noter :

`sync/atomic` ne fonctionne qu'avec certains types primitifs (`int32`, `int64`, etc.). Il est très rapide et efficace, mais n'est pas adapté pour gérer des structures complexes.

---

### 🧩 Exemple complet : compteur avec `sync/atomic`

```go
package main

import (
    "github.com/gin-gonic/gin"
    "net/http"
    "sync/atomic"
)

type AtomicCounter struct {
    value int64
}

func (c *AtomicCounter) Increment() {
    atomic.AddInt64(&c.value, 1)
}

func (c *AtomicCounter) Get() int64 {
    return atomic.LoadInt64(&c.value)
}

func main() {
    r := gin.Default()
    counter := &AtomicCounter{}

    r.POST("/increment", func(c *gin.Context) {
        counter.Increment()
        c.JSON(http.StatusOK, gin.H{"message": "increased"})
    })

    r.GET("/value", func(c *gin.Context) {
        value := counter.Get()
        c.JSON(http.StatusOK, gin.H{"value": value})
    })

    r.Run(":8080")
}
```

---

### ✅ Avantages de `sync/atomic`

* Plus rapide que `sync.Mutex`
* Idéal pour des **opérations simples** et **fréquentes** (compteurs, flags, etc.)
* Évite les blocages (`Lock`/`Unlock`)

---

### 🧪 Test rapide

```bash
curl -X POST http://localhost:8080/increment
curl http://localhost:8080/value
```

---

Souhaites-tu que je t'intègre ça dans une structure MVC complète (controllers, services, etc.) ?
