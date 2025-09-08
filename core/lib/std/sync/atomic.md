En Go, le **package `sync/atomic`** fournit des opérations bas niveau pour effectuer des **manipulations atomiques** sur des variables partagées entre plusieurs goroutines, sans avoir besoin de verrous (`mutex`). Cela permet de garantir la **sécurité concurrente** tout en étant généralement plus rapide que l'utilisation de verrous traditionnels.

---

### 📌 Pourquoi "atomique" ?

Une opération **atomique** est une opération qui se produit **entièrement ou pas du tout**, sans être interrompue. Cela signifie qu'aucune autre goroutine ne peut voir un état intermédiaire de l’opération.

---

### 🔧 Types supportés

Le package `sync/atomic` fonctionne principalement sur des types de base :

* `int32`, `int64`
* `uint32`, `uint64`
* `uintptr`
* `unsafe.Pointer`

---

### ✨ Fonctions courantes

#### 1. `atomic.AddInt64` / `atomic.AddInt32`

Additionne une valeur à une variable de manière atomique.

```go
var counter int64

atomic.AddInt64(&counter, 1) // Incrémente de 1
```

#### 2. `atomic.LoadInt64` / `atomic.LoadInt32`

Lit une valeur de manière atomique.

```go
val := atomic.LoadInt64(&counter)
```

#### 3. `atomic.StoreInt64` / `atomic.StoreInt32`

Assigne une valeur de manière atomique.

```go
atomic.StoreInt64(&counter, 0)
```

#### 4. `atomic.CompareAndSwapInt64`

Compare la valeur actuelle à une valeur attendue, et la remplace si elles sont égales (CAS : Compare And Swap).

```go
ok := atomic.CompareAndSwapInt64(&counter, 0, 1)
```

---

### ✅ Exemple simple

```go
package main

import (
    "fmt"
    "sync"
    "sync/atomic"
)

func main() {
    var counter int64 = 0
    var wg sync.WaitGroup

    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func() {
            atomic.AddInt64(&counter, 1)
            wg.Done()
        }()
    }

    wg.Wait()
    fmt.Println("Counter:", counter) // Affiche 100
}
```

---

### ⚠️ Attention

* Le package `sync/atomic` ne remplace pas les verrous pour toutes les situations.
* Il est adapté pour des **opérations simples** sur des **types de base**.
* Pour des structures complexes ou des enchaînements de plusieurs opérations, il vaut mieux utiliser un `sync.Mutex`.

Souhaites-tu un exemple avec `atomic.Value`, qui permet de stocker des objets de manière atomique ?

Parfait ! Voici une explication et un exemple complet avec **`atomic.Value`** en Go.

---

### 🔍 Qu’est-ce que `atomic.Value` ?

`atomic.Value` est un type fourni par le package `sync/atomic` qui permet de **lire et d'écrire de manière atomique** **n’importe quelle valeur** (même des structures personnalisées), tant que **le type est constant** (c’est-à-dire toujours du même type pour toutes les écritures).

---

### ✅ Méthodes principales

* `Store(interface{})` : écrit une valeur de manière atomique.
* `Load() interface{}` : lit une valeur de manière atomique.

---

### ⚠️ Contraintes

* Tous les `Store()` doivent utiliser **le même type** après la première écriture, sinon un panic se produit.
* Il ne faut pas modifier une valeur lue depuis `Load()` sans copie, sinon cela casse la sécurité concurrente.

---

### 📦 Exemple : Configuration dynamique en mémoire

```go
package main

import (
    "fmt"
    "sync/atomic"
    "time"
)

type Config struct {
    AppName string
    Debug   bool
}

var config atomic.Value

func loadInitialConfig() {
    cfg := Config{
        AppName: "MonApp",
        Debug:   false,
    }
    config.Store(cfg)
}

func updateConfigPeriodically() {
    for {
        time.Sleep(3 * time.Second)
        newCfg := Config{
            AppName: "MonApp",
            Debug:   true, // On imagine une mise à jour ici
        }
        config.Store(newCfg)
        fmt.Println("Configuration mise à jour !")
    }
}

func readConfigLoop() {
    for {
        cfg := config.Load().(Config) // Type assertion sûre
        fmt.Printf("Lecture de la config : AppName=%s, Debug=%v\n", cfg.AppName, cfg.Debug)
        time.Sleep(1 * time.Second)
    }
}

func main() {
    loadInitialConfig()
    go updateConfigPeriodically()
    readConfigLoop()
}
```

---

### 💡 Cas d’usage typiques de `atomic.Value`

* Lire une configuration mise à jour à chaud sans redémarrer un service.
* Partager un cache en lecture seule entre goroutines.
* Utiliser un *flag* ou un *état global* atomique de type struct.

Souhaites-tu un exemple plus avancé, comme un cache en mémoire concurrent sécurisé avec `atomic.Value` ?

Parfait ! Voici un exemple **plus avancé** d'utilisation de `atomic.Value` : un **cache en mémoire** partagé entre goroutines, mis à jour régulièrement, et **sûr pour la concurrence** sans verrou explicite.

---

## 🔧 Scénario : Cache utilisateur en lecture seule

* Le cache contient des données utilisateurs.
* Un processus met à jour périodiquement les données (ex : depuis une base de données ou une API).
* Les goroutines peuvent lire le cache en toute sécurité **sans blocage**.

---

### ✅ Exemple complet

```go
package main

import (
    "fmt"
    "sync/atomic"
    "time"
)

type User struct {
    ID   int
    Name string
}

type UserCache struct {
    Users map[int]User
}

// Cache global
var userCache atomic.Value

func loadInitialUsers() {
    users := map[int]User{
        1: {ID: 1, Name: "Alice"},
        2: {ID: 2, Name: "Bob"},
    }
    userCache.Store(UserCache{Users: users})
    fmt.Println("Cache initialisé.")
}

func updateUserCache() {
    for {
        time.Sleep(5 * time.Second)

        // Simule une mise à jour depuis la DB
        updatedUsers := map[int]User{
            1: {ID: 1, Name: "Alice Updated"},
            2: {ID: 2, Name: "Bob Updated"},
            3: {ID: 3, Name: "Charlie"},
        }

        // Nouvelle copie de la structure complète
        userCache.Store(UserCache{Users: updatedUsers})
        fmt.Println("Cache mis à jour.")
    }
}

func getUser(id int) (User, bool) {
    cache := userCache.Load().(UserCache)
    user, ok := cache.Users[id]
    return user, ok
}

func simulateUserRequests() {
    for {
        user, ok := getUser(2)
        if ok {
            fmt.Printf("[LECTURE] Utilisateur ID=2 : %+v\n", user)
        } else {
            fmt.Println("[LECTURE] Utilisateur non trouvé")
        }
        time.Sleep(2 * time.Second)
    }
}

func main() {
    loadInitialUsers()
    go updateUserCache()
    simulateUserRequests()
}
```

---

### ⚙️ Pourquoi ça marche ?

* `atomic.Value` permet de **remplacer tout le cache d’un coup**.
* Aucune goroutine ne modifie une carte partagée en direct.
* Toutes les lectures utilisent une **copie immuable**.
* Aucun `mutex` n’est nécessaire, donc lecture ultra rapide et concurrente.

---
