Une **race condition** (ou **condition de concurrence**) en Go — comme dans beaucoup de langages concurrents — se produit lorsque **deux ou plusieurs goroutines accèdent à une même variable partagée en même temps**, **et qu'au moins une d'entre elles effectue une écriture**. Cela peut provoquer un comportement non déterministe, c'est-à-dire que le résultat du programme peut varier à chaque exécution.

### Expliqué simplement

Imagine que deux personnes essaient de modifier un même fichier en même temps sans coordination. L'une commence à écrire pendant que l'autre lit ou écrit aussi. Le contenu final du fichier dépendra de qui a été le plus rapide — c'est une race condition.

---

### Exemple simple en Go

```go
package main

import (
    "fmt"
    "time"
)

var counter int

func increment() {
    for i := 0; i < 1000; i++ {
        counter++
    }
}

func main() {
    go increment()
    go increment()

    time.Sleep(1 * time.Second)
    fmt.Println("Counter:", counter)
}
```

#### Que se passe-t-il ici ?

* Deux goroutines appellent `increment()`, chacune essayant d’augmenter `counter`.
* `counter++` n’est **pas une opération atomique** — elle lit, incrémente, puis écrit.
* Résultat : **le compteur n’atteindra probablement pas 2000**, car des mises à jour sont écrasées.

---

### Comment détecter les race conditions en Go ?

Go fournit un outil très utile :

```bash
go run -race main.go
```

Ce flag permet de détecter les accès concurrents à une variable et vous alertera s’il y a un problème.

---

### Comment éviter les race conditions ?

Voici plusieurs méthodes :

#### 1. **Mutex (verrou)**

```go
import "sync"

var mu sync.Mutex

func increment() {
    for i := 0; i < 1000; i++ {
        mu.Lock()
        counter++
        mu.Unlock()
    }
}
```

#### 2. **Canaux (channels)** – modèle de communication par message

```go
func increment(ch chan int) {
    for i := 0; i < 1000; i++ {
        ch <- 1
    }
}

func main() {
    ch := make(chan int)
    go increment(ch)
    go increment(ch)

    go func() {
        for i := 0; i < 2000; i++ {
            counter += <-ch
        }
    }()

    time.Sleep(1 * time.Second)
    fmt.Println("Counter:", counter)
}
```

#### 3. **sync/atomic** – opérations atomiques rapides

```go
import "sync/atomic"

var counter int32

func increment() {
    for i := 0; i < 1000; i++ {
        atomic.AddInt32(&counter, 1)
    }
}
```

---

### À retenir

* Une race condition rend le programme imprévisible.
* Utilisez `-race` pendant le développement pour les détecter.
* Protégez les accès concurrents avec `sync.Mutex`, `sync/atomic` ou des `chan`.

Souhaites-tu un exemple avec un scénario plus concret (ex : accès concurrent à une base de données ou à une map) ?

Prenons un exemple **concret** : plusieurs goroutines accèdent à une **map partagée** (simulant un cache mémoire), ce qui peut provoquer une race condition. Ensuite, on corrigera ce problème avec un `sync.Mutex`.

---

## 🔥 Exemple 1 : Accès concurrent à une map (race condition)

```go
package main

import (
    "fmt"
    "time"
)

var cache = make(map[string]string)

func writeToCache(key, value string) {
    cache[key] = value
}

func readFromCache(key string) string {
    return cache[key]
}

func main() {
    for i := 0; i < 10; i++ {
        go writeToCache(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))
        go readFromCache(fmt.Sprintf("key%d", i))
    }

    time.Sleep(1 * time.Second)
    fmt.Println("Cache:", cache)
}
```

### 🧨 Problème

* La map `cache` est **modifiée** et **lue** en même temps par plusieurs goroutines.
* En Go, **les maps ne sont pas thread-safe**, donc ce code peut panic ou donner des résultats imprévisibles.
* Lance le avec `go run -race` pour détecter les conflits.

---

## ✅ Exemple 2 : Version corrigée avec `sync.Mutex`

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

var (
    cache = make(map[string]string)
    mu    sync.Mutex
)

func writeToCache(key, value string) {
    mu.Lock()
    cache[key] = value
    mu.Unlock()
}

func readFromCache(key string) string {
    mu.Lock()
    defer mu.Unlock()
    return cache[key]
}

func main() {
    for i := 0; i < 10; i++ {
        go writeToCache(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))
        go readFromCache(fmt.Sprintf("key%d", i))
    }

    time.Sleep(1 * time.Second)
    fmt.Println("Cache:", cache)
}
```

### ✅ Comportement

* Le `sync.Mutex` empêche les accès simultanés à la map.
* Le programme devient **sûr et déterministe**.

---

Tu veux que je te montre comment faire ça avec des `channels` ou avec `sync.Map` (une map concurrente intégrée à Go) ?

Parfait, voyons deux alternatives élégantes à `sync.Mutex` : **`channels`** (modèle CSP de Go) et **`sync.Map`**, une map concurrente native.

---

## ✅ Option 1 : Utiliser des **channels** (approche orientée message)

Ici, on encapsule l’accès à la map dans une goroutine "gestionnaire" qui reçoit les opérations via des canaux.

```go
package main

import (
    "fmt"
    "time"
)

type request struct {
    key      string
    value    string
    isWrite  bool
    response chan string
}

func cacheManager(reqChan chan request) {
    cache := make(map[string]string)
    for req := range reqChan {
        if req.isWrite {
            cache[req.key] = req.value
        } else {
            req.response <- cache[req.key]
        }
    }
}

func main() {
    reqChan := make(chan request)

    // Lancement du gestionnaire de cache
    go cacheManager(reqChan)

    // Écriture
    for i := 0; i < 10; i++ {
        reqChan <- request{
            key:     fmt.Sprintf("key%d", i),
            value:   fmt.Sprintf("value%d", i),
            isWrite: true,
        }
    }

    // Lecture
    for i := 0; i < 10; i++ {
        resp := make(chan string)
        reqChan <- request{
            key:      fmt.Sprintf("key%d", i),
            isWrite:  false,
            response: resp,
        }
        fmt.Printf("key%d = %s\n", i, <-resp)
    }

    time.Sleep(500 * time.Millisecond)
}
```

### ✔️ Avantages

* Pas de besoin de mutex : un seul gestionnaire contrôle l’accès.
* Approche très Go-style (inspirée d’Erlang et CSP).

---

## ✅ Option 2 : Utiliser `sync.Map` (map thread-safe intégrée)

`sync.Map` est optimisée pour un grand nombre de lectures et un petit nombre d’écritures.

```go
package main

import (
    "fmt"
    "sync"
)

var cache sync.Map

func writeToCache(key, value string) {
    cache.Store(key, value)
}

func readFromCache(key string) string {
    value, ok := cache.Load(key)
    if ok {
        return value.(string)
    }
    return "not found"
}

func main() {
    for i := 0; i < 10; i++ {
        go writeToCache(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))
        go func(i int) {
            fmt.Println(readFromCache(fmt.Sprintf("key%d", i)))
        }(i)
    }

    // Juste pour attendre les goroutines (dans la vraie vie, utilise sync.WaitGroup)
    fmt.Scanln()
}
```

### ✔️ Avantages

* Pas besoin de gérer la synchronisation soi-même.
* Utile dans les cas de **cache partagé très lu, peu écrit**.
* Interface différente de `map` (méthodes `Store`, `Load`, `Delete`, `Range`).

---

Souhaites-tu que je t’intègre l’un de ces modèles dans un microservice REST (ex. cache utilisateur en mémoire) ?

Très bien ! Voici un petit **microservice REST en Go avec Gin**, qui expose une API de type **cache clé/valeur en mémoire**.

🛠️ Il utilise :

* **Gin** pour le serveur HTTP,
* **`sync.Map`** pour stocker les données en mémoire, de manière **thread-safe**,
* JSON pour les entrées/sorties.

---

## 📦 Structure des endpoints

| Méthode | Endpoint      | Description             |
| ------- | ------------- | ----------------------- |
| `GET`   | `/cache/:key` | Lire une valeur par clé |
| `POST`  | `/cache`      | Écrire une clé/valeur   |

---

## ✅ Code complet avec `sync.Map`

```go
package main

import (
    "github.com/gin-gonic/gin"
    "net/http"
    "sync"
)

var cache sync.Map

type CacheEntry struct {
    Key   string `json:"key" binding:"required"`
    Value string `json:"value" binding:"required"`
}

func main() {
    r := gin.Default()

    // POST /cache
    r.POST("/cache", func(c *gin.Context) {
        var entry CacheEntry
        if err := c.ShouldBindJSON(&entry); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        cache.Store(entry.Key, entry.Value)
        c.JSON(http.StatusOK, gin.H{"message": "stored"})
    })

    // GET /cache/:key
    r.GET("/cache/:key", func(c *gin.Context) {
        key := c.Param("key")
        if value, ok := cache.Load(key); ok {
            c.JSON(http.StatusOK, gin.H{"key": key, "value": value})
        } else {
            c.JSON(http.StatusNotFound, gin.H{"error": "key not found"})
        }
    })

    r.Run(":8080")
}
```

---

## 🧪 Test rapide avec `curl`

```bash
# Stocker une valeur
curl -X POST http://localhost:8080/cache \
    -H "Content-Type: application/json" \
    -d '{"key":"user42", "value":"Alice"}'

# Récupérer la valeur
curl http://localhost:8080/cache/user42
```

---

### 🧩 Variante possible

Souhaites-tu une version avec un **mutex + map** au lieu de `sync.Map`, ou bien avec un **TTL (expiration)** des entrées comme un vrai cache ?
