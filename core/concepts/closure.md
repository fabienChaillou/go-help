En Go, une **closure** est une **fonction anonyme** (ou littérale de fonction) qui **capture des variables** de son **contexte extérieur** — c’est-à-dire qu’elle peut accéder aux variables définies **hors de son corps**, même après la fin de leur portée originale.

---

## 🔹 Définition simple

Une **closure** est une **fonction qui se souvient de son environnement** au moment de sa création.

---

## 🔹 Exemple simple

```go
func compteur() func() int {
    i := 0
    return func() int {
        i++
        return i
    }
}

func main() {
    c := compteur()
    fmt.Println(c()) // 1
    fmt.Println(c()) // 2
    fmt.Println(c()) // 3
}
```

🟢 **Explication :**

* `compteur` retourne une **fonction anonyme**.
* Cette fonction interne **utilise** la variable `i` définie dans la fonction `compteur`.
* Même après que `compteur()` ait terminé, la fonction retournée **garde une référence à `i`**, qui continue à vivre en mémoire.

---

## 🔹 Pourquoi c'est utile ?

* **Fonctions génératrices** (comme un compteur ou un générateur de données)
* **Encapsulation d’état** sans utiliser de struct
* **Callbacks** ou fonctions de traitement personnalisées

---

## 🔹 Exemple avec fonction littérale

```go
func main() {
    message := "Bonjour"

    f := func() {
        fmt.Println(message)
    }

    f() // Affiche : Bonjour
}
```

Si tu modifies `message` après avoir défini la closure, celle-ci utilisera la **dernière valeur connue** (car elle y accède par **référence**).

---

## 🔹 Closure avec arguments

```go
func multiplier(facteur int) func(int) int {
    return func(x int) int {
        return x * facteur
    }
}

func main() {
    doubler := multiplier(2)
    tripler := multiplier(3)

    fmt.Println(doubler(5)) // 10
    fmt.Println(tripler(5)) // 15
}
```

---

## 🔹 À retenir

* Les closures sont **des fonctions avec mémoire**.
* Elles permettent d'**éviter les variables globales** tout en conservant un **état privé**.
* Elles sont très puissantes dans des contextes comme les **fonctions retournant des fonctions**, les **filtres**, les **callbacks**, etc.

---


Bien sûr ! Voici un **exemple complet** de test d’une **fonction closure** (fonction anonyme qui capture des variables de son environnement) en Go.

---

## 🎯 Objectif

On va créer une fonction qui retourne une **closure compteur**, puis écrire un test unitaire avec `testing`.

---

### 📦 Fichier : `counter.go`

```go
package counter

// NewCounter retourne une fonction closure qui incrémente un compteur
func NewCounter() func() int {
    count := 0
    return func() int {
        count++
        return count
    }
}
```

---

### ✅ Fichier de test : `counter_test.go`

```go
package counter

import "testing"

func TestNewCounter(t *testing.T) {
    c := NewCounter()

    if c() != 1 {
        t.Errorf("expected 1, got %d", c())
    }

    if c() != 2 {
        t.Errorf("expected 2, got %d", c())
    }

    if c() != 3 {
        t.Errorf("expected 3, got %d", c())
    }
}
```

---

## 💡 Explication :

* `NewCounter` retourne une fonction qui **capture** la variable `count`.
* À chaque appel, `count` est **incrémenté et conservé** (comportement typique d’une **closure**).
* Le test appelle 3 fois la fonction et vérifie que les valeurs sont bien 1, 2, 3.

---

## ▶️ Pour exécuter le test :

Dans le terminal :

```bash
go test
```

---

Souhaites-tu voir un exemple avec une **closure qui prend des arguments** ou qui modifie un état externe (comme un logger, une slice, etc.) ?

### autre example

Parfait ! Voici un **exemple de closure en Go qui prend un argument**, modifie un **état interne**, et qu’on peut tester.

---

## 🎯 Objectif

Créer une **closure avec état** qui accumule les valeurs données à chaque appel (comme un `sum()` progressif), et écrire un test pour s’assurer que l’accumulation fonctionne.

---

### 📦 Fichier : `accumulator.go`

```go
package accumulator

// NewAccumulator retourne une closure qui accumule les valeurs passées
func NewAccumulator() func(int) int {
    total := 0
    return func(x int) int {
        total += x
        return total
    }
}
```

---

### ✅ Fichier de test : `accumulator_test.go`

```go
package accumulator

import "testing"

func TestNewAccumulator(t *testing.T) {
    acc := NewAccumulator()

    tests := []struct {
        input    int
        expected int
    }{
        {3, 3},
        {5, 8},
        {-2, 6},
        {10, 16},
    }

    for i, tt := range tests {
        got := acc(tt.input)
        if got != tt.expected {
            t.Errorf("test %d: expected %d, got %d", i+1, tt.expected, got)
        }
    }
}
```

---

## 💡 Ce que fait cette closure :

* À chaque appel, elle prend un `int`, l’ajoute à un total interne.
* Elle retourne le total actuel.
* Le test vérifie que ce total évolue correctement.

---

### ▶️ Pour exécuter :

```bash
go test
```

---

---------------------------------------------------------------------------------------------------------------------------

Souhaites-tu un exemple où la closure modifie un **état externe** (comme une `slice` partagée, ou un logger personnalisé) ?

## change etat externe

Parfait ! Voici un exemple plus avancé d’une **closure qui modifie un état externe**, en l’occurrence une `slice` partagée. C’est une technique parfois utilisée pour construire des **collecteurs ou des traceurs** personnalisés.

---

## 🎯 Objectif

Créer une closure qui ajoute des messages dans une **slice externe** (`[]string`) à chaque appel, comme un **logger simple**.

---

### 📦 Fichier : `logger.go`

```go
package logger

// NewLogger retourne une closure qui ajoute un message dans une slice externe
func NewLogger(storage *[]string) func(string) {
    return func(msg string) {
        *storage = append(*storage, msg)
    }
}
```

---

### ✅ Fichier de test : `logger_test.go`

```go
package logger

import (
    "testing"
)

func TestNewLogger(t *testing.T) {
    var logs []string
    log := NewLogger(&logs)

    log("Start")
    log("Running")
    log("Done")

    expected := []string{"Start", "Running", "Done"}

    if len(logs) != len(expected) {
        t.Fatalf("expected %d logs, got %d", len(expected), len(logs))
    }

    for i, msg := range expected {
        if logs[i] != msg {
            t.Errorf("log[%d]: expected %q, got %q", i, msg, logs[i])
        }
    }
}
```

---

### 💡 Fonctionnement :

* `NewLogger` prend un pointeur vers une `slice` de `string`.
* La closure retournée **ajoute un message dans cette slice externe**.
* Le test vérifie que la slice contient bien tous les messages dans le bon ordre.

---

### ✅ Intérêts d’une closure avec état externe :

* Partage d’état sans globales.
* Personnalisation légère (logger, accumulateur, tampon, etc.).
* Réutilisable et testable.

---

---------------------------------------------------------------------------------------------------------------------------

Souhaites-tu maintenant un exemple de **closure concurrente** (utilisée dans une goroutine avec synchronisation) ?

Super ! Voici un exemple de **closure utilisée de manière concurrente** en Go, avec gestion d’un **état partagé** de façon **sécurisée via un mutex**.

---

## 🎯 Objectif

Créer une closure **thread-safe** (sécurisée pour les goroutines), qui incrémente un compteur partagé entre plusieurs goroutines.

---

### 📦 Fichier : `safe_counter.go`

```go
package safecounter

import "sync"

// NewSafeCounter retourne une closure qui incrémente un compteur en toute sécurité
func NewSafeCounter() func() int {
    var mu sync.Mutex
    count := 0

    return func() int {
        mu.Lock()
        defer mu.Unlock()
        count++
        return count
    }
}
```

---

### ✅ Fichier de test : `safe_counter_test.go`

```go
package safecounter

import (
    "sync"
    "testing"
)

func TestSafeCounter(t *testing.T) {
    counter := NewSafeCounter()

    var wg sync.WaitGroup
    n := 1000

    for i := 0; i < n; i++ {
        wg.Add(1)
        go func() {
            counter()
            wg.Done()
        }()
    }

    wg.Wait()

    final := counter()
    if final != n+1 {
        t.Errorf("expected %d, got %d", n+1, final)
    }
}
```

---

## 💡 Explication :

* `NewSafeCounter` retourne une **closure avec mutex** pour éviter les accès concurrents non synchronisés.
* Le test lance **1000 goroutines** qui appellent toutes la même closure.
* À la fin, on vérifie que le compteur a bien été incrémenté **exactement 1000 fois**.
* On appelle une fois de plus la closure pour récupérer la valeur finale (`n + 1`).

---

### 🔒 Pourquoi le `sync.Mutex` est essentiel ici :

Sans verrou (`mu.Lock()`), plusieurs goroutines pourraient lire et modifier `count` en même temps → **condition de course** (`race condition`), ce qui donne des résultats incorrects.

---

Tu peux tester la sécurité de ce code avec :

```bash
go test -race
```

Cela détectera toute condition de course potentielle.

---

Souhaites-tu un exemple de closure qui **retourne un canal** (chan) ou qui **filtre des valeurs** de manière fonctionnelle, comme un pipeline ?

---------------------------------------------------------------------------------------------------------------------------

Souhaites-tu un exemple d'utilisation de closures dans un cas réel comme un filtre, un trie ou un handler HTTP ?

Excellent choix 👌 ! Les **closures sont très utiles dans des cas réels** en Go, surtout quand tu veux **personnaliser du comportement à la volée**, comme :

* Des **filtres de données**
* Des **fonctions de tri (sort)**
* Des **handlers HTTP avec paramètres dynamiques**

---

## 🧪 Voici 3 exemples réels et pratiques d'utilisation des closures en Go :

---

### ✅ 1. **Closure comme filtre personnalisé (ex. sur un slice)**

```go
func Filter(nums []int, predicate func(int) bool) []int {
    var result []int
    for _, n := range nums {
        if predicate(n) {
            result = append(result, n)
        }
    }
    return result
}

func main() {
    data := []int{1, 2, 3, 4, 5, 6}

    // Utilisation d'une closure comme filtre pair
    evens := Filter(data, func(n int) bool {
        return n%2 == 0
    })

    fmt.Println(evens) // ➜ [2 4 6]
}
```

🔍 Ici, la fonction `Filter` reçoit une **fonction en paramètre** — la closure définie inline (`func(n int) bool`).

---

### ✅ 2. **Closure pour tri personnalisé (`sort.Slice`)**

```go
import (
    "fmt"
    "sort"
)

type Person struct {
    Name string
    Age  int
}

func main() {
    people := []Person{
        {"Alice", 30},
        {"Bob", 25},
        {"Charlie", 35},
    }

    // Tri croissant par âge (closure)
    sort.Slice(people, func(i, j int) bool {
        return people[i].Age < people[j].Age
    })

    fmt.Println(people)
}
```

🔍 `sort.Slice` prend une closure qui compare deux éléments. Très pratique pour **trier dynamiquement** selon n’importe quel champ.

---

### ✅ 3. **Closure dans un handler HTTP avec paramètre externe**

```go
package main

import (
    "fmt"
    "net/http"
)

func makeHandler(greeting string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        name := r.URL.Query().Get("name")
        if name == "" {
            name = "world"
        }
        fmt.Fprintf(w, "%s, %s!", greeting, name)
    }
}

func main() {
    http.HandleFunc("/hello", makeHandler("Hello"))
    http.HandleFunc("/hi", makeHandler("Hi"))
    http.ListenAndServe(":8080", nil)
}
```

🧠 Ici, `makeHandler` est une **factory de handlers** : chaque handler retourne une closure qui **capture le paramètre `greeting`**.

Accède à :

* `http://localhost:8080/hello?name=Alice` → `Hello, Alice!`
* `http://localhost:8080/hi` → `Hi, world!`

---

## ✅ Conclusion : Pourquoi utiliser des closures dans ces cas ?

| Cas          | Intérêt de la closure                                          |
| ------------ | -------------------------------------------------------------- |
| Filtrage     | Passer une logique inline, flexible                            |
| Tri          | Comparateur dynamique basé sur des champs différents           |
| Handler HTTP | Créer des handlers configurables sans créer plein de fonctions |

---

---------------------------------------------------------------------------------------------------------------------------

Souhaites-tu que je t’aide à écrire un mini serveur HTTP avec des **middlewares en closures** (ex : logger, auth, etc.) ?

Parfait ! Voici un **exemple clair et complet** d’un **mini serveur HTTP** en Go avec l’utilisation de **middlewares en closures** – une approche très utilisée dans les APIs et backends Go.

---

## 🎯 Objectif

Créer un serveur avec :

* Un **handler principal** (`/hello`)
* Un **middleware logger** qui affiche chaque requête
* Un **middleware auth** qui bloque si un token n’est pas présent

---

## 📦 Code complet (dans `main.go`)

```go
package main

import (
    "fmt"
    "log"
    "net/http"
)

// middlewareLogger est une closure qui wrap un handler pour afficher les logs HTTP
func middlewareLogger(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("%s %s\n", r.Method, r.URL.Path)
        next.ServeHTTP(w, r) // Appelle le handler suivant
    })
}

// middlewareAuth vérifie la présence d’un token ?auth=secret
func middlewareAuth(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Query().Get("auth") != "secret" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        next.ServeHTTP(w, r)
    })
}

// helloHandler est le handler de base
func helloHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Hello, Go web world!")
}

func main() {
    // Handler avec les 2 middlewares appliqués
    finalHandler := http.HandlerFunc(helloHandler)
    wrapped := middlewareLogger(middlewareAuth(finalHandler))

    http.Handle("/hello", wrapped)
    log.Println("Server started on http://localhost:8080")
    http.ListenAndServe(":8080", nil)
}
```

---

## 🧪 Comportement

* `http://localhost:8080/hello` → ❌ 401 Unauthorized
* `http://localhost:8080/hello?auth=secret` → ✅ `Hello, Go web world!`

🧠 Chaque middleware est une **closure** qui capture le handler suivant, et décide :

* de passer la requête à la suite (`next.ServeHTTP`)
* ou de l’interrompre (ex : pour bloquer, logger, modifier)

---

## 🧩 Bonus : chaînage plus propre (si souhaité)

```go
func chainMiddlewares(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
    for i := len(middlewares) - 1; i >= 0; i-- {
        h = middlewares[i](h)
    }
    return h
}

// Utilisation :
wrapped := chainMiddlewares(http.HandlerFunc(helloHandler), middlewareLogger, middlewareAuth)
```

---

---------------------------------------------------------------------------------------------------------------------------

Souhaites-tu que je t’aide à intégrer aussi des middlewares plus avancés comme :

* **mesure de temps de réponse**
* **en-têtes CORS**
* ou une **authentification via token Bearer** ?

Parfait 🙌 ! Voici une version enrichie de notre mini serveur HTTP Go avec **middlewares avancés**, y compris :

1. ✅ **Logger HTTP**
2. 🔐 **Auth via token `Authorization: Bearer <token>`**
3. ⏱️ **Mesure du temps de réponse**
4. 🌐 **CORS** (Cross-Origin Resource Sharing)

---

## 🧠 Structure

Chaque middleware est une **closure** qui "wrappe" un `http.Handler`.

---

## 📦 Code complet (`main.go`)

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    "strings"
    "time"
)

// Middleware 1: Logger
func middlewareLogger(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("[LOG] %s %s", r.Method, r.URL.Path)
        next.ServeHTTP(w, r)
    })
}

// Middleware 2: Auth avec header Authorization: Bearer secret123
func middlewareAuth(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        auth := r.Header.Get("Authorization")
        if !strings.HasPrefix(auth, "Bearer ") || strings.TrimPrefix(auth, "Bearer ") != "secret123" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        next.ServeHTTP(w, r)
    })
}

// Middleware 3: Timer
func middlewareTimer(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        duration := time.Since(start)
        log.Printf("[TIMER] %s took %v", r.URL.Path, duration)
    })
}

// Middleware 4: CORS
func middlewareCORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")

        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }

        next.ServeHTTP(w, r)
    })
}

// Handler final
func helloHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "✅ Hello, secured Go world!")
}

// Helper pour chaîner les middlewares
func chain(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
    for i := len(middlewares) - 1; i >= 0; i-- {
        h = middlewares[i](h)
    }
    return h
}

func main() {
    handler := chain(
        http.HandlerFunc(helloHandler),
        middlewareLogger,
        middlewareTimer,
        middlewareAuth,
        middlewareCORS,
    )

    http.Handle("/hello", handler)

    log.Println("🚀 Server started on http://localhost:8080")
    http.ListenAndServe(":8080", nil)
}
```

---

## 🧪 Exemple de test avec `curl`

```bash
# Échec : pas d'auth
curl http://localhost:8080/hello

# Échec : mauvais token
curl -H "Authorization: Bearer wrong" http://localhost:8080/hello

# ✅ Succès
curl -H "Authorization: Bearer secret123" http://localhost:8080/hello
```

---

## 🧩 Ce que tu peux faire ensuite :

* Ajouter un middleware de **log JSON structuré**
* Créer un middleware pour **limiter le nombre de requêtes par IP (rate limiting)**
* Utiliser un **router externe** comme `chi`, `gorilla/mux` ou `Gin` qui supporte nativement ce style de middlewares

---------------------------------------------------------------------------------------------------------------------------

Souhaites-tu un exemple avec un **router tiers** comme `chi` ou `gin` pour comparer ?

Parfait ! Comparons maintenant notre approche standard (avec `net/http`) à une version équivalente utilisant un **router tiers**. Voici un exemple avec [`chi`](https://github.com/go-chi/chi), un routeur très léger, performant et **parfaitement adapté aux middlewares en Go**.

---

## 🛠️ Installation de `chi`

```bash
go get github.com/go-chi/chi/v5
```

---

## ✅ Objectif

Reproduire notre `/hello` avec :

* Middleware **logger**
* Middleware **auth Bearer**
* Middleware **timer**
* Middleware **CORS**
* Tout en utilisant `chi.Router`

---

## 📦 Code complet avec `chi` (`main.go`)

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    "strings"
    "time"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
)

// Auth middleware personnalisé
func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
        if token != "secret123" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        next.ServeHTTP(w, r)
    })
}

// Timer middleware personnalisé
func timerMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        log.Printf("[TIMER] %s took %v", r.URL.Path, time.Since(start))
    })
}

// CORS middleware (très simple)
func corsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusOK)
            return
        }
        next.ServeHTTP(w, r)
    })
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "✅ Hello from chi router!")
}

func main() {
    r := chi.NewRouter()

    // Middlewares globaux
    r.Use(middleware.Logger)     // chi's built-in logger
    r.Use(timerMiddleware)
    r.Use(corsMiddleware)

    // Route protégée
    r.With(authMiddleware).Get("/hello", helloHandler)

    log.Println("🚀 chi server started on http://localhost:8080")
    http.ListenAndServe(":8080", r)
}
```

---

## 🧪 Tester avec `curl`

```bash
# Échec (pas de token)
curl http://localhost:8080/hello

# ✅ Succès
curl -H "Authorization: Bearer secret123" http://localhost:8080/hello
```

---

## ✅ Pourquoi `chi` est intéressant ?

| Caractéristique       | Avantage                                     |
| --------------------- | -------------------------------------------- |
| 📌 `.With(...)`       | Ajoute un middleware à une route spécifique  |
| ⚙️ `.Use(...)`        | Ajoute des middlewares globaux               |
| 🔀 Très rapide        | Idéal pour des APIs REST                     |
| ✅ Compatible net/http | Tu peux toujours utiliser `http.HandlerFunc` |

---

Souhaites-tu maintenant voir une version équivalente avec **Gin**, un framework plus haut niveau avec JSON, binding et rendu intégrés ?

---------------------------------------------------------------------------------------------------------------------------
