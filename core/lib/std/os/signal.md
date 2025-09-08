Le package `os/signal` en Go est utilisé pour intercepter les **signaux du système d'exploitation** (comme `SIGINT`, `SIGTERM`, etc.) afin que ton programme puisse réagir proprement (par exemple, en faisant un nettoyage avant de quitter).

---

### 🔧 Utilisation typique

Voici un exemple simple qui montre comment capturer `Ctrl+C` (`SIGINT`) et quitter proprement :

```go
package main

import (
    "fmt"
    "os"
    "os/signal"
    "syscall"
)

func main() {
    // Crée un canal pour recevoir les signaux
    sigs := make(chan os.Signal, 1)

    // Enregistre les signaux qu'on souhaite intercepter
    signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

    // Bloque jusqu'à réception d'un signal
    fmt.Println("Attente du signal (Ctrl+C pour interrompre)...")
    sig := <-sigs

    fmt.Println()
    fmt.Println("Signal reçu :", sig)
    fmt.Println("Arrêt propre du programme.")
}
```

---

### 🔍 Détail des composants

| Élément                   | Description                                                        |
| ------------------------- | ------------------------------------------------------------------ |
| `os.Signal`               | Interface représentant un signal du système.                       |
| `syscall.SIGINT`          | Signal envoyé par `Ctrl+C` (interruption clavier).                 |
| `syscall.SIGTERM`         | Signal de terminaison standard.                                    |
| `signal.Notify()`         | Enregistre un ou plusieurs signaux à écouter sur un canal.         |
| `make(chan os.Signal, 1)` | Crée un canal pour recevoir les signaux (tampon d’un élément ici). |

---

### 🧼 Exemple avec cleanup (nettoyage avant de quitter)

```go
go func() {
    sig := <-sigs
    fmt.Println("\nSignal reçu:", sig)
    fmt.Println("Fermeture des ressources...")
    // ex: fermer la DB, arrêter les workers, etc.
    os.Exit(0)
}()
```

> Cela est souvent utilisé dans des serveurs web ou applications longue durée pour gérer un arrêt proprement (graceful shutdown).

Souhaites-tu un exemple d'arrêt propre avec un serveur HTTP en Go ?

Parfait ! Voici un **exemple complet** d'un serveur HTTP en Go qui utilise `os/signal` pour effectuer un **arrêt propre (graceful shutdown)** :

---

### 🌐 Exemple avec `net/http`, `os/signal`, et `context`

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
)

func main() {
    // Crée un mux simple
    mux := http.NewServeMux()
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "Hello, world!")
    })

    // Crée le serveur
    srv := &http.Server{
        Addr:    ":8080",
        Handler: mux,
    }

    // Canal pour recevoir les signaux
    stop := make(chan os.Signal, 1)
    signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

    // Démarre le serveur dans une goroutine
    go func() {
        log.Println("Serveur démarré sur :8080")
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Erreur serveur: %v", err)
        }
    }()

    // Attente du signal
    <-stop
    log.Println("Signal reçu, arrêt en cours...")

    // Contexte avec timeout pour le shutdown
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        log.Fatalf("Erreur lors de l'arrêt du serveur: %v", err)
    }

    log.Println("Serveur arrêté proprement")
}
```

---

### ✅ Ce que fait ce code :

* Lance un serveur HTTP.
* Écoute les signaux `SIGINT` et `SIGTERM`.
* Quand l’un de ces signaux est reçu (ex : `Ctrl+C`), il :

  * Crée un contexte avec timeout de 5 secondes.
  * Appelle `srv.Shutdown(ctx)` pour arrêter proprement les connexions en cours.
  * Affiche des logs à chaque étape.


### signal.NotifyContext()

Très bon choix ! `signal.NotifyContext` est une **fonction moderne et pratique** introduite dans Go 1.16. Elle permet de créer directement un **`context.Context` annulé automatiquement lorsqu’un signal système est reçu** — ce qui évite de gérer manuellement un canal et un `<-stop`.

---

### 🧠 Concept clé

```go
ctx, stop := signal.NotifyContext(parentCtx, os.Interrupt, syscall.SIGTERM)
```

* `ctx` : un contexte que tu peux passer à des fonctions, serveurs, etc.
* `stop()` : une fonction à appeler pour **libérer les ressources** liées au contexte (bonnes pratiques).

Quand l’un des signaux spécifiés (`SIGINT`, `SIGTERM`, etc.) est reçu :

* `ctx.Done()` est fermé → tu peux détecter ça avec `<-ctx.Done()`
* Le contexte est automatiquement annulé.

---

### 🧪 Exemple d’utilisation dans un serveur HTTP

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
)

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "Bonjour !")
    })

    srv := &http.Server{
        Addr:    ":8080",
        Handler: mux,
    }

    // Crée un contexte qui s'annule lors d'un signal système
    ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
    defer stop()

    // Démarre le serveur
    go func() {
        log.Println("Serveur HTTP démarré sur :8080")
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Erreur serveur: %v", err)
        }
    }()

    // Attend que le contexte soit annulé
    <-ctx.Done()
    log.Println("Signal reçu, arrêt du serveur...")

    // Contexte de timeout pour shutdown
    shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := srv.Shutdown(shutdownCtx); err != nil {
        log.Fatalf("Erreur lors de l'arrêt: %v", err)
    }

    log.Println("Arrêt propre effectué")
}
```

---

### ✅ Avantages de `signal.NotifyContext`

| Avantage                                        | Détail                                                          |
| ----------------------------------------------- | --------------------------------------------------------------- |
| ✅ Moins de code boilerplate                     | Pas besoin de `make(chan os.Signal, 1)` ni de goroutine dédiée. |
| ✅ Compatible avec `context.Context`             | Peut être passé à des fonctions, workers, serveurs, etc.        |
| ✅ Fonctionne bien avec `select` et `ctx.Done()` | Idéal dans des systèmes concurrentiels.                         |
| ✅ Fournit `stop()`                              | Permet de libérer les ressources liées au signal listener.      |

---

Souhaites-tu un exemple avec des workers ou goroutines qui écoutent ce contexte pour s’arrêter proprement aussi ?

Parfait ! Voici un **exemple complet** avec un **serveur HTTP + un worker en goroutine**, tous deux écoutant un `context.Context` généré par `signal.NotifyContext` pour **s'arrêter proprement** quand un signal système est reçu (comme `Ctrl+C`).

---

### 🧵 Serveur + worker avec arrêt propre via `context`

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
)

func main() {
    // Crée un contexte qui s'annule sur SIGINT/SIGTERM
    ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
    defer stop()

    // Lance un worker simulant une tâche périodique
    go startWorker(ctx)

    // Démarre un serveur HTTP
    mux := http.NewServeMux()
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "Bonjour depuis le serveur HTTP !")
    })

    srv := &http.Server{
        Addr:    ":8080",
        Handler: mux,
    }

    go func() {
        log.Println("Serveur HTTP lancé sur :8080")
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Erreur serveur: %v", err)
        }
    }()

    // Attente du signal système
    <-ctx.Done()
    log.Println("Signal reçu, arrêt des composants...")

    // Timeout pour shutdown propre du serveur HTTP
    shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if err := srv.Shutdown(shutdownCtx); err != nil {
        log.Printf("Erreur arrêt serveur HTTP : %v", err)
    }

    log.Println("Fin propre du programme")
}

func startWorker(ctx context.Context) {
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()

    log.Println("Worker démarré")

    for {
        select {
        case <-ctx.Done():
            log.Println("Arrêt du worker demandé")
            return
        case t := <-ticker.C:
            log.Println("Worker actif à", t.Format(time.RFC3339))
        }
    }
}
```

---

### 🧠 Ce que fait ce programme :

* Le **serveur HTTP** accepte les connexions sur `:8080`.
* Un **worker** exécute une tâche répétée toutes les secondes.
* Quand `Ctrl+C` est pressé :

  * Le `context` est annulé.
  * Le worker détecte `ctx.Done()` et s’arrête.
  * Le serveur est arrêté proprement avec un timeout.
  * Le programme se termine correctement.

---

Parfait, allons un cran plus loin ! Voici un **exemple avec plusieurs workers**, chacun exécutant une **tâche simulée plus complexe** (ex. traitement de fichiers, appel API, etc.), **interrompable proprement via `context.Context`**.

---

### 🏗️ Structure du programme

* Plusieurs **jobs** sont lancés en goroutines.
* Chaque job prend entre 2 et 5 secondes à s’exécuter.
* Si un signal (`SIGINT`, `SIGTERM`) est reçu, tous les jobs en cours sont interrompus proprement.

---

### 🧪 Exemple : plusieurs jobs simulés avec annulation

```go
package main

import (
    "context"
    "fmt"
    "log"
    "math/rand"
    "os"
    "os/signal"
    "sync"
    "syscall"
    "time"
)

func main() {
    rand.Seed(time.Now().UnixNano())

    // Crée un contexte qui s'annule sur interruption
    ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
    defer stop()

    log.Println("Démarrage des workers...")

    var wg sync.WaitGroup
    numWorkers := 3

    for i := 1; i <= numWorkers; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            worker(ctx, id)
        }(i)
    }

    // Attente du signal
    <-ctx.Done()
    log.Println("Signal reçu, arrêt en cours...")

    // Attente que tous les workers terminent
    wg.Wait()
    log.Println("Tous les workers sont arrêtés proprement")
}

func worker(ctx context.Context, id int) {
    log.Printf("[Worker %d] Démarré", id)

    for task := 1; task <= 5; task++ {
        select {
        case <-ctx.Done():
            log.Printf("[Worker %d] Arrêt demandé, interruption du job %d", id, task)
            return
        default:
            // Simule une tâche complexe
            duration := time.Duration(rand.Intn(3)+2) * time.Second
            log.Printf("[Worker %d] Lancement du job %d (durée : %v)", id, task, duration)

            select {
            case <-time.After(duration):
                log.Printf("[Worker %d] Job %d terminé", id, task)
            case <-ctx.Done():
                log.Printf("[Worker %d] Interruption pendant le job %d", id, task)
                return
            }
        }
    }

    log.Printf("[Worker %d] Tous les jobs sont terminés", id)
}
```

---

### 🧠 Ce que fait ce programme :

| Élément                | Fonction                                                            |
| ---------------------- | ------------------------------------------------------------------- |
| `signal.NotifyContext` | Crée un `context.Context` annulé à la réception d’un signal système |
| `sync.WaitGroup`       | Attend que tous les workers aient terminé                           |
| `time.After(duration)` | Simule une tâche qui prend du temps                                 |
| `ctx.Done()`           | Permet une interruption propre à tout moment                        |

---

### 📌 Bonus

Ce modèle est **parfait pour des traitements batch**, des **pipelines**, ou encore des **services worker** en Go qui doivent pouvoir être interrompus de manière fiable et propre.

Souhaites-tu que je t’aide à encapsuler cette logique dans une architecture plus modulaire (ex : package `jobrunner`, injection de logger, etc.) ?


Parfait ! Structurons ce code pour qu’il soit **modulaire, maintenable et réutilisable**. On va créer un **petit framework de gestion de jobs** qu’on peut facilement intégrer dans n’importe quel projet Go.

---

## 🧱 Objectif de l’architecture

* Séparer la logique des workers dans un **package `jobrunner`**.
* Permettre de :

  * Créer des workers avec identifiant.
  * Gérer proprement le shutdown via `context`.
  * Loguer les actions avec un logger injecté.

---

## 📁 Structure proposée

```
/project-root
├── main.go
└── jobrunner/
    └── runner.go
```

---

### 📄 `jobrunner/runner.go`

```go
package jobrunner

import (
    "context"
    "log"
    "math/rand"
    "sync"
    "time"
)

type Worker struct {
    ID     int
    Logger *log.Logger
}

// Simule une série de jobs complexes
func (w *Worker) Run(ctx context.Context, wg *sync.WaitGroup) {
    defer wg.Done()
    w.Logger.Printf("[Worker %d] Démarré", w.ID)

    for task := 1; task <= 5; task++ {
        select {
        case <-ctx.Done():
            w.Logger.Printf("[Worker %d] Arrêt demandé, job %d interrompu", w.ID, task)
            return
        default:
            duration := time.Duration(rand.Intn(3)+2) * time.Second
            w.Logger.Printf("[Worker %d] Démarrage du job %d (durée : %v)", w.ID, task, duration)

            select {
            case <-time.After(duration):
                w.Logger.Printf("[Worker %d] Job %d terminé", w.ID, task)
            case <-ctx.Done():
                w.Logger.Printf("[Worker %d] Interruption pendant le job %d", w.ID, task)
                return
            }
        }
    }

    w.Logger.Printf("[Worker %d] Tous les jobs sont terminés", w.ID)
}
```

---

### 📄 `main.go`

```go
package main

import (
    "context"
    "log"
    "math/rand"
    "os"
    "os/signal"
    "sync"
    "syscall"
    "time"

    "project/jobrunner" // Remplace "project" par ton module Go
)

func main() {
    rand.Seed(time.Now().UnixNano())

    ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
    defer stop()

    logger := log.New(os.Stdout, "", log.LstdFlags)
    var wg sync.WaitGroup

    numWorkers := 3
    logger.Println("Lancement de", numWorkers, "workers...")

    for i := 1; i <= numWorkers; i++ {
        worker := jobrunner.Worker{
            ID:     i,
            Logger: logger,
        }
        wg.Add(1)
        go worker.Run(ctx, &wg)
    }

    <-ctx.Done()
    logger.Println("Signal système reçu. Attente de l'arrêt des workers...")

    wg.Wait()
    logger.Println("Tous les workers sont terminés proprement")
}
```

---

### ✅ Avantages

| Caractéristique     | Bénéfice                           |
| ------------------- | ---------------------------------- |
| Package modulaire   | Réutilisable dans d’autres projets |
| Injection du logger | Testable et configurable           |
| Context partagé     | Contrôle uniforme de l’arrêt       |
| Sync.WaitGroup      | Attente sûre de la fin des jobs    |

---

Souhaites-tu que je t’aide à ajouter :

* un **timeout par job individuel** ?
* une interface `Job` pour supporter différents types de tâches ?
* ou bien une version concurrente avec **file d’attente des jobs** ?

