Le module `os/signal` en Go permet à un programme de **recevoir et de gérer des signaux du système d'exploitation** (comme `SIGINT`, `SIGTERM`, etc.). Ces signaux sont souvent utilisés pour **interrompre, terminer ou redémarrer proprement un programme**, comme lorsqu'on appuie sur `Ctrl+C` dans un terminal.

---

### 📦 Importation

```go
import (
    "os"
    "os/signal"
    "syscall"
)
```

---

### 🔧 Utilisation typique

Voici un exemple simple pour **intercepter un signal SIGINT** (`Ctrl+C`) et **faire un arrêt propre** du programme :

```go
package main

import (
    "fmt"
    "os"
    "os/signal"
    "syscall"
)

func main() {
    // Canal pour recevoir les signaux
    sigs := make(chan os.Signal, 1)

    // Notifier le canal lorsqu'on reçoit SIGINT ou SIGTERM
    signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

    fmt.Println("En attente du signal (Ctrl+C pour quitter)")

    // Bloque ici jusqu'à réception d'un signal
    sig := <-sigs

    fmt.Println("Signal reçu :", sig)
    fmt.Println("Arrêt du programme proprement.")
}
```

---

### 🔄 Détails des fonctions importantes

| Fonction                   | Description                                                  |
| -------------------------- | ------------------------------------------------------------ |
| `signal.Notify(c, sig...)` | Enregistre un canal `c` pour recevoir les signaux spécifiés. |
| `signal.Stop(c)`           | Arrête d'envoyer des signaux au canal `c`.                   |
| `signal.Reset(sig...)`     | Restaure le comportement par défaut des signaux.             |

---

### 🧠 Exemple avancé : gestion d’un serveur HTTP

Si tu as un serveur HTTP ou un processus long, tu peux utiliser ce module pour arrêter proprement :

```go
go func() {
    sig := <-sigs
    log.Printf("Signal reçu : %s", sig)
    // Appelle ici la fonction d'arrêt (ex : server.Shutdown)
}()
```

---

Souhaites-tu un exemple complet avec un serveur HTTP ou une application CLI ?

Parfait ! Voici un exemple **complet d'application CLI** en Go qui utilise le module `os/signal` pour gérer une interruption (`Ctrl+C`) proprement, par exemple pour nettoyer des ressources ou sauvegarder avant de quitter.

---

### 🧪 Exemple : CLI avec gestion de `Ctrl+C`

```go
package main

import (
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "time"
)

func main() {
    // Création d’un canal pour les signaux
    sigs := make(chan os.Signal, 1)

    // Enregistre les signaux d'interruption (Ctrl+C) et de terminaison
    signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

    // Canal pour détecter la fin du programme
    done := make(chan bool, 1)

    // Goroutine qui attend un signal
    go func() {
        sig := <-sigs
        fmt.Println("\n🔴 Signal reçu :", sig)
        fmt.Println("💡 Nettoyage avant de quitter...")
        time.Sleep(1 * time.Second) // simulation du nettoyage
        done <- true
    }()

    fmt.Println("👋 CLI en cours d'exécution. Appuyez sur Ctrl+C pour arrêter.")
    fmt.Println("⌛ Simulation de travail...")

    // Boucle simulant un travail en cours
    go func() {
        for {
            fmt.Print(".")
            time.Sleep(500 * time.Millisecond)
        }
    }()

    // Attend que `done` soit déclenché
    <-done
    fmt.Println("\n✅ Arrêt propre effectué.")
}
```

---

### ✅ Ce que ce programme fait :

* Lance une tâche de fond (boucle de simulation).
* Attend qu’un signal soit reçu (comme `Ctrl+C`).
* Lorsqu’il est reçu, il affiche un message, attend un peu pour simuler un nettoyage, puis quitte proprement.


--------------------------------------------------------------------------------------------------------------------------------------------

Tu veux que je t’aide à transformer ça en binaire CLI installable (`go install`) ou à y ajouter des commandes avec `cobra` ou `urfave/cli` ?

Parfait ! Voici comment intégrer la **gestion de signaux (`os/signal`) dans une CLI basée sur le framework [`cobra`](https://github.com/spf13/cobra)** – très utilisé pour créer des outils en ligne de commande en Go.

---

### 🧱 Structure finale :

```
mycli/
├── cmd/
│   └── root.go
├── main.go
```

---

### 1. 📦 Initialisation du projet

Assure-toi d’avoir `cobra` installé :

```bash
go install github.com/spf13/cobra-cli@latest
```

Crée une base de projet :

```bash
mkdir mycli && cd mycli
go mod init mycli
cobra-cli init
```

---

### 2. ✏️ `main.go`

```go
package main

import "mycli/cmd"

func main() {
    cmd.Execute()
}
```

---

### 3. 🧠 `cmd/root.go` avec gestion du signal

```go
package cmd

import (
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "mycli",
    Short: "Une CLI avec gestion de Ctrl+C",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("👋 Bienvenue dans la CLI. Appuyez sur Ctrl+C pour arrêter.")

        // Canal pour les signaux
        sigs := make(chan os.Signal, 1)
        signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

        done := make(chan bool, 1)

        go func() {
            sig := <-sigs
            fmt.Println("\n🔴 Signal reçu :", sig)
            fmt.Println("🧹 Nettoyage...")
            time.Sleep(1 * time.Second)
            done <- true
        }()

        // Simulation de tâche longue
        go func() {
            for {
                fmt.Print(".")
                time.Sleep(500 * time.Millisecond)
            }
        }()

        <-done
        fmt.Println("\n✅ Fin propre du programme.")
    },
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
```

---

### 4. 🚀 Test de la CLI

```bash
go run main.go
```

Tu verras des `.` défiler dans le terminal. Appuie sur `Ctrl+C` et la CLI affichera un message de nettoyage avant de s'arrêter proprement.

---

Souhaites-tu que je t’aide à :

* Ajouter une sous-commande ?
* Générer un binaire installable (`go install`) ?
* Ajouter un flag (ex: `--verbose`) ?

--------------------------------------------------------------------------------------------------------------------------------------------

Absolument ! On peut tester certaines parties de la CLI, même si la **gestion des signaux** elle-même est difficile à tester directement en test unitaire (car elle dépend du système). En revanche, on peut isoler et tester **les fonctions de logique métier** (comme la logique de nettoyage, les flags, la structure des commandes, etc.).

---

### ✅ Objectif :

Ajouter des **tests unitaires pour la commande `root`** et simuler une partie du comportement de la CLI sans attendre un vrai `Ctrl+C`.

---

### 1. 📁 Structure avec fichier de test

Crée un fichier `cmd/root_test.go` :

```
mycli/
├── cmd/
│   ├── root.go
│   └── root_test.go  ← ✅
```

---

### 2. 🧪 `cmd/root_test.go` — test de la commande `root`

```go
package cmd

import (
    "bytes"
    "strings"
    "testing"

    "github.com/spf13/cobra"
)

func TestRootCommandHelp(t *testing.T) {
    // Capture la sortie de la commande
    buf := new(bytes.Buffer)
    rootCmd.SetOut(buf)
    rootCmd.SetArgs([]string{"--help"})

    err := rootCmd.Execute()
    if err != nil {
        t.Fatalf("Erreur à l'exécution de la commande : %v", err)
    }

    output := buf.String()
    if !strings.Contains(output, "Une CLI avec gestion de Ctrl+C") {
        t.Errorf("L'aide ne contient pas la description attendue. Sortie : %s", output)
    }
}
```

---

### 3. 🚫 Pourquoi on ne teste pas directement `os/signal`

Les tests unitaires **ne peuvent pas intercepter ou déclencher directement `SIGINT` ou `SIGTERM`** sans faire du test d’intégration avec processus séparé, ce qui sort du cadre unitaire. Pour cela, on **déplacerait la logique métier dans une fonction testable**.

---

### 4. 🧼 Refactor : fonction testable `runApp()`

Dans `root.go`, on isole la logique principale dans une fonction :

```go
func runApp(stop <-chan struct{}) {
    fmt.Println("👋 Bienvenue dans la CLI. Appuyez sur Ctrl+C pour arrêter.")

    go func() {
        for {
            fmt.Print(".")
            time.Sleep(500 * time.Millisecond)
        }
    }()

    <-stop
    fmt.Println("\n✅ Fin propre du programme.")
}
```

Puis, dans le `Run`, on l’utilise :

```go
Run: func(cmd *cobra.Command, args []string) {
    sigs := make(chan os.Signal, 1)
    done := make(chan struct{}, 1)
    signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        <-sigs
        fmt.Println("\n🔴 Signal reçu. Nettoyage...")
        time.Sleep(1 * time.Second)
        close(done)
    }()

    runApp(done)
},
```

---

### 5. 🔁 Test de `runApp` avec stop simulé

Ajoute dans `root_test.go` :

```go
func TestRunAppSimulé(t *testing.T) {
    stop := make(chan struct{})
    go func() {
        time.Sleep(1 * time.Second)
        close(stop) // simule un signal d'arrêt
    }()

    runApp(stop)
}
```

📝 Ce test **vérifie que `runApp` réagit correctement à un signal d’arrêt simulé**, sans avoir besoin de signaux réels du système.

---

Souhaites-tu que je t’aide à :

* Ajouter des sous-commandes avec tests ?
* Convertir ça en test d’intégration ?
* Utiliser `cobra` avec `viper` pour ajouter des options/flags configurables ?
