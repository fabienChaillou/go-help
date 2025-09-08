La fonction `**signal.NotifyContext**` en Go est une **version améliorée** de `signal.Notify` introduite dans Go 1.16. Elle permet de **lier la réception d'un signal système à un contexte `context.Context`**. Cela facilite la **gestion propre de l'annulation** d'opérations longues (serveur, traitement, etc.) lorsqu’un signal comme `SIGINT` ou `SIGTERM` est reçu.

---

### 📌 Prototype

```go
func NotifyContext(parent context.Context, signals ...os.Signal) (ctx context.Context, stop context.CancelFunc)
```

* `parent`: le contexte parent (souvent `context.Background()`).
* `signals`: les signaux à écouter (`syscall.SIGINT`, etc.).
* Retourne :

  * `ctx`: un contexte annulé automatiquement quand un de ces signaux est reçu.
  * `stop`: une fonction à appeler pour arrêter manuellement l'écoute.

---

### ✅ Avantages

* Pas besoin de gérer manuellement les canaux de signaux.
* Compatible avec tous les outils et APIs utilisant `context.Context`.
* Plus lisible et idiomatique pour du code moderne en Go.

---

### 🧪 Exemple simple

```go
package main

import (
    "context"
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "time"
)

func main() {
    // Crée un contexte lié à SIGINT ou SIGTERM
    ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
    defer stop()

    fmt.Println("En attente du signal... (Ctrl+C)")

    select {
    case <-ctx.Done():
        fmt.Println("🛑 Signal reçu, contexte annulé.")
    case <-time.After(10 * time.Second):
        fmt.Println("⏰ Temps écoulé, aucun signal reçu.")
    }

    fmt.Println("Fin propre.")
}
```

---

### 🔁 Utilisation typique avec un serveur ou traitement

```go
func run(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            fmt.Println("Annulation demandée.")
            return
        default:
            fmt.Println("Traitement en cours...")
            time.Sleep(1 * time.Second)
        }
    }
}
```

---

--------------------------------------------------------------------------------------------------------------------------------

Souhaites-tu un exemple avec un serveur HTTP ou une tâche CLI intégrée à `cobra` avec `NotifyContext` ?

Parfait ! Voici un exemple **complet** d’un **serveur HTTP en Go** utilisant `http.Server` et `signal.NotifyContext` pour permettre **un arrêt propre avec `Ctrl+C` ou `SIGTERM`**.

---

## ✅ Exemple : Serveur HTTP + `signal.NotifyContext`

```go
package main

import (
    "context"
    "fmt"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
)

func main() {
    // Crée un contexte qui sera annulé quand un signal est reçu
    ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
    defer stop()

    mux := http.NewServeMux()
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "Bienvenue sur le serveur 👋")
    })

    srv := &http.Server{
        Addr:    ":8080",
        Handler: mux,
    }

    // Démarre le serveur dans une goroutine
    go func() {
        fmt.Println("🚀 Serveur démarré sur http://localhost:8080")
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            fmt.Printf("❌ Erreur serveur : %v\n", err)
        }
    }()

    // Bloque ici jusqu'à ce que le signal soit reçu
    <-ctx.Done()
    fmt.Println("\n🛑 Signal reçu. Arrêt du serveur en cours...")

    // Contexte de timeout pour l'arrêt propre
    shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := srv.Shutdown(shutdownCtx); err != nil {
        fmt.Printf("❌ Erreur à l'arrêt du serveur : %v\n", err)
    } else {
        fmt.Println("✅ Serveur arrêté proprement.")
    }
}
```

---

### 📌 Ce que ce code fait :

1. Crée un contexte `ctx` avec `signal.NotifyContext` lié à `SIGINT`/`SIGTERM`.
2. Lance le serveur HTTP dans une goroutine.
3. Attend l’annulation du contexte (signal reçu).
4. Utilise `Shutdown()` avec un `context.WithTimeout` pour permettre aux connexions de se fermer proprement.

---

Souhaites-tu que ce serveur soit intégré dans une commande `cobra` ou qu’il expose une API REST simple ?

--------------------------------------------------------------------------------------------------------------------------------

Parfait ! Voici comment intégrer ce **serveur HTTP avec gestion de `SIGINT` / `SIGTERM` via `signal.NotifyContext`** dans une **commande `cobra`**.

---

## 🧱 Structure de projet recommandée

```
mycli/
├── cmd/
│   ├── root.go
│   └── serve.go       ← ✅ Commande `serve`
├── main.go
```

---

## 1. ✨ `main.go`

```go
package main

import "mycli/cmd"

func main() {
    cmd.Execute()
}
```

---

## 2. 🧩 `cmd/root.go`

```go
package cmd

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "mycli",
    Short: "Une CLI avec un serveur HTTP intégré",
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
```

---

## 3. 🚀 `cmd/serve.go`

```go
package cmd

import (
    "context"
    "fmt"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
    Use:   "serve",
    Short: "Lance un serveur HTTP avec arrêt propre",
    Run: func(cmd *cobra.Command, args []string) {
        ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
        defer stop()

        mux := http.NewServeMux()
        mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
            fmt.Fprintln(w, "Bienvenue sur le serveur CLI 👋")
        })

        srv := &http.Server{
            Addr:    ":8080",
            Handler: mux,
        }

        go func() {
            fmt.Println("🚀 Serveur disponible sur http://localhost:8080")
            if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
                fmt.Printf("❌ Erreur serveur : %v\n", err)
            }
        }()

        <-ctx.Done()
        fmt.Println("\n🛑 Signal reçu. Arrêt du serveur...")

        shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()

        if err := srv.Shutdown(shutdownCtx); err != nil {
            fmt.Printf("❌ Erreur à l'arrêt : %v\n", err)
        } else {
            fmt.Println("✅ Serveur arrêté proprement.")
        }
    },
}

func init() {
    rootCmd.AddCommand(serveCmd)
}
```

---

## ✅ Test

Dans le terminal :

```bash
go run main.go serve
```

Puis :

* Accède à [http://localhost:8080](http://localhost:8080)
* Appuie sur `Ctrl+C` → le serveur se ferme proprement

---

Souhaites-tu :

* Ajouter un **flag `--port`** ?
* Ajouter un **point d’API REST (ex: `/health`)** ?
* Générer un **binaire installable** (`go install`) ?


--------------------------------------------------------------------------------------------------------------------------------

### Ajout de Unite test

Oui, ajoutons des **tests unitaires** pour la commande `serve` de ta CLI `cobra`. Comme le serveur HTTP tourne de manière asynchrone et interagit avec le système via `os/signal`, on va tester ce qu'on peut :

---

## ✅ Ce qu'on va tester :

1. Que la **commande `serve` existe et s’exécute** sans erreur.
2. Isoler la **logique serveur dans une fonction testable**.
3. Simuler un arrêt via **contexte annulé** pour tester le shutdown.

---

## 🧱 Fichier de test : `cmd/serve_test.go`

```go
package cmd

import (
    "context"
    "net/http"
    "testing"
    "time"
)

// Test que la fonction de démarrage et d'arrêt du serveur fonctionne sans erreur
func TestServeLogic(t *testing.T) {
    // Contexte annulé après 1 seconde
    ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
    defer cancel()

    // Fonction qui retourne un handler basique (pour test)
    handler := http.NewServeMux()
    handler.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
    })

    // Appelle la fonction serveur dans une goroutine
    go func() {
        if err := runHTTPServer(ctx, ":8081", handler); err != nil && err != http.ErrServerClosed {
            t.Errorf("Erreur inattendue du serveur : %v", err)
        }
    }()

    // Donne un peu de temps au serveur pour démarrer
    time.Sleep(200 * time.Millisecond)

    // Requête simple vers le serveur
    resp, err := http.Get("http://localhost:8081/test")
    if err != nil {
        t.Fatalf("Erreur lors de la requête : %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        t.Errorf("Code HTTP inattendu : %d", resp.StatusCode)
    }
}
```

---

## 🔁 Refactor : isoler la logique serveur

Modifie `serve.go` pour y ajouter une fonction testable :

```go
// Fonction réutilisable pour lancer le serveur HTTP avec contexte
func runHTTPServer(ctx context.Context, addr string, handler http.Handler) error {
    srv := &http.Server{
        Addr:    addr,
        Handler: handler,
    }

    // Lancement en arrière-plan
    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            fmt.Printf("❌ Serveur erreur: %v\n", err)
        }
    }()

    <-ctx.Done()
    fmt.Println("\n🛑 Contexte annulé. Fermeture serveur...")

    shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    return srv.Shutdown(shutdownCtx)
}
```

Puis, dans `serveCmd`, utilise-la :

```go
Run: func(cmd *cobra.Command, args []string) {
    ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
    defer stop()

    mux := http.NewServeMux()
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "Bienvenue sur le serveur CLI 👋")
    })

    if err := runHTTPServer(ctx, ":8080", mux); err != nil {
        fmt.Printf("❌ Erreur arrêt du serveur : %v\n", err)
    }
},
```

---

## 🚀 Pour exécuter le test :

```bash
go test ./cmd -v
```

---

Souhaites-tu ajouter un mock pour tester d’autres routes ou un test de performance ?

-------------------------------


### test avec un mock

Parfait ! Ajoutons un **mock de handler HTTP** et un test unitaire plus poussé pour :

* Vérifier qu’une route personnalisée retourne le bon résultat.
* Utiliser un **serveur de test (`httptest.Server`)** si besoin.
* Simuler différents comportements (200 OK, erreurs, délais, etc.).

---

## 🧪 1. Ajout d’un handler HTTP mocké

On va créer un handler avec plusieurs routes à tester :

```go
func createTestHandler() http.Handler {
    mux := http.NewServeMux()

    mux.HandleFunc("/ok", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        fmt.Fprint(w, "OK")
    })

    mux.HandleFunc("/slow", func(w http.ResponseWriter, r *http.Request) {
        time.Sleep(2 * time.Second)
        w.WriteHeader(http.StatusGatewayTimeout)
    })

    mux.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    })

    return mux
}
```

---

## 🧪 2. Ajout des tests unitaires dans `serve_test.go`

```go
func TestServeWithMockHandlers(t *testing.T) {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    handler := createTestHandler()

    // Démarre le serveur sur un port différent
    go func() {
        _ = runHTTPServer(ctx, ":8082", handler)
    }()

    time.Sleep(300 * time.Millisecond) // Laisse le serveur démarrer

    tests := []struct {
        path       string
        wantCode   int
        wantBody   string
    }{
        {"/ok", http.StatusOK, "OK"},
        {"/error", http.StatusInternalServerError, "Internal Server Error"},
    }

    client := &http.Client{Timeout: 1 * time.Second}

    for _, tc := range tests {
        t.Run(tc.path, func(t *testing.T) {
            resp, err := client.Get("http://localhost:8082" + tc.path)
            if err != nil {
                t.Fatalf("Erreur requête : %v", err)
            }
            defer resp.Body.Close()

            if resp.StatusCode != tc.wantCode {
                t.Errorf("Code HTTP attendu %d, obtenu %d", tc.wantCode, resp.StatusCode)
            }

            buf := new(bytes.Buffer)
            _, _ = buf.ReadFrom(resp.Body)
            if !strings.Contains(buf.String(), tc.wantBody) {
                t.Errorf("Réponse attendue contenant %q, obtenue : %q", tc.wantBody, buf.String())
            }
        })
    }
}
```

---

## ✅ Résultat : ce test vérifie que…

* `/ok` retourne bien `200 OK` avec `"OK"`.
* `/error` retourne bien `500 Internal Server Error`.
* Le serveur gère correctement des routes différentes via un `ServeMux` personnalisé.

---

Souhaites-tu aussi :

* Tester `/slow` avec un timeout forcé ?
* Couvrir une route avec un paramètre (ex: `/user?id=42`) ?
* Générer un rapport de couverture (`go test -cover`) ?

---------------------------------------------

Parfait ! Voici comment compléter notre test pour :

### ✅ 1. Tester une **route lente** `/slow` avec timeout

### ✅ 2. Ajouter une **route avec paramètre de requête** `/user?id=42`

### ✅ 3. Générer un **rapport de couverture de test**

---

## 🧩 1. Ajout de la route `/user?id=42` dans le handler mock

Modifie ou complète la fonction `createTestHandler` :

```go
mux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Query().Get("id")
    if id == "" {
        http.Error(w, "Missing ID", http.StatusBadRequest)
        return
    }
    fmt.Fprintf(w, "User ID: %s", id)
})
```

---

## 🧪 2. Mise à jour des cas de test

Ajoute dans la liste `tests` du test `TestServeWithMockHandlers` :

```go
{"/user?id=42", http.StatusOK, "User ID: 42"},
{"/user", http.StatusBadRequest, "Missing ID"},
```

Et pour `/slow`, on le teste à part car on veut déclencher un timeout :

---

## ⏳ 3. Test spécial pour `/slow` avec timeout forcé

Ajoute ce test à part dans `serve_test.go` :

```go
func TestSlowRouteTimeout(t *testing.T) {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    handler := createTestHandler()

    go func() {
        _ = runHTTPServer(ctx, ":8083", handler)
    }()

    time.Sleep(300 * time.Millisecond)

    client := &http.Client{
        Timeout: 1 * time.Second, // timeout forcé
    }

    _, err := client.Get("http://localhost:8083/slow")
    if err == nil {
        t.Fatal("On attendait une erreur de timeout, mais aucune n’a été reçue")
    }
    if !strings.Contains(err.Error(), "Client.Timeout") {
        t.Errorf("Erreur inattendue : %v", err)
    }
}
```

---

## 📊 4. Générer un **rapport de couverture**

Dans le terminal, exécute :

```bash
go test ./cmd -cover -coverprofile=coverage.out
go tool cover -func=coverage.out
```

Et pour **voir une vue HTML** :

```bash
go tool cover -html=coverage.out
```

Cela ouvrira un rapport visuel montrant quelles lignes sont couvertes ou non.

---

Souhaites-tu qu’on ajoute :

* Des **tests pour erreurs serveur** (ex : port déjà utilisé) ?
* Une **intégration continue** (ex: GitHub Actions pour lancer les tests) ?

