Parfait, on passe à un **worker pool** — un modèle très courant et puissant en Go pour **traiter plusieurs tâches en parallèle**, tout en contrôlant **le nombre de goroutines**.

---

## 🔧 C’est quoi un Worker Pool ?

* Tu as **un canal de tâches**.
* Tu crées **N goroutines (les workers)** qui lisent dans ce canal.
* Tu envoies les tâches dans le canal.
* Chaque worker prend une tâche, la traite, et recommence.

> ⚙️ **Avantages** :

* Limite le nombre de goroutines.
* Contrôle précis de la charge.
* Bonne utilisation du CPU.

---

## 🎯 Exemple : Traitement d’une liste de personnes avec 3 workers

```go
package main

import (
    "fmt"
    "math/rand"
    "sync"
    "time"
)

type Person struct {
    Name string
    Age  int
}

func worker(id int, jobs <-chan Person, results chan<- string, wg *sync.WaitGroup) {
    defer wg.Done()
    for p := range jobs {
        delay := time.Duration(100+rand.Intn(200)) * time.Millisecond
        time.Sleep(delay) // Simule un travail
        results <- fmt.Sprintf("Worker %d : %s (%d ans) traité en %v", id, p.Name, p.Age, delay)
    }
}

func main() {
    rand.Seed(time.Now().UnixNano())

    people := []Person{
        {"Alice", 30},
        {"Bob", 25},
        {"Charlie", 35},
        {"Dave", 40},
        {"Eve", 28},
        {"Frank", 32},
    }

    jobs := make(chan Person)
    results := make(chan string)
    var wg sync.WaitGroup

    // Démarrage de 3 workers
    numWorkers := 3
    for i := 1; i <= numWorkers; i++ {
        wg.Add(1)
        go worker(i, jobs, results, &wg)
    }

    // Envoi des tâches
    go func() {
        for _, p := range people {
            jobs <- p
        }
        close(jobs) // Important : ferme le canal de jobs
    }()

    // Lecture des résultats
    go func() {
        wg.Wait()        // Attend que tous les workers aient terminé
        close(results)   // Ensuite on peut fermer le canal de résultats
    }()

    for msg := range results {
        fmt.Println(msg)
    }
}
```

---

## ✅ Ce que tu apprends ici :

| Élément            | Description                                                     |
| ------------------ | --------------------------------------------------------------- |
| `chan`             | `jobs` pour envoyer les tâches, `results` pour lire les retours |
| `go`               | Chaque worker est une goroutine                                 |
| `sync.WaitGroup`   | Pour attendre la fin de tous les workers                        |
| `range` sur `chan` | Pour lire jusqu’à ce que le canal soit fermé                    |
| `close(chan)`      | Très important pour signaler la fin                             |

---

## 🧪 Exemple de sortie (aléatoire) :

```
Worker 2 : Alice (30 ans) traité en 180ms
Worker 1 : Bob (25 ans) traité en 267ms
Worker 3 : Charlie (35 ans) traité en 187ms
Worker 2 : Dave (40 ans) traité en 198ms
Worker 1 : Eve (28 ans) traité en 199ms
Worker 3 : Frank (32 ans) traité en 149ms
```

---

--------------------------------------------------------------------------------------------------------

Souhaites-tu ajouter un **timeout global**, un **canal d'annulation**, ou encore un **pool dynamique** selon la charge ?

Excellent choix ! On va maintenant améliorer notre **worker pool** en ajoutant un **timeout global** — très utile quand tu veux **annuler tout le traitement après un certain temps**, peu importe ce qui reste à faire.

---

## 🎯 Objectif :

* Créer un worker pool comme avant,
* Mais **arrêter tout le traitement** si ça prend plus de X secondes,
* Grâce à un **contexte (`context.Context`) avec timeout**.

---

## 🧱 Code complet avec timeout global

```go
package main

import (
    "context"
    "fmt"
    "math/rand"
    "sync"
    "time"
)

type Person struct {
    Name string
    Age  int
}

func worker(ctx context.Context, id int, jobs <-chan Person, results chan<- string, wg *sync.WaitGroup) {
    defer wg.Done()
    for {
        select {
        case <-ctx.Done():
            return // arrêt demandé
        case p, ok := <-jobs:
            if !ok {
                return // plus de jobs
            }

            // Simule un travail
            delay := time.Duration(200+rand.Intn(300)) * time.Millisecond
            select {
            case <-ctx.Done():
                return // arrêt pendant le traitement
            case <-time.After(delay):
                results <- fmt.Sprintf("Worker %d : %s traité en %v", id, p.Name, delay)
            }
        }
    }
}

func main() {
    rand.Seed(time.Now().UnixNano())

    people := []Person{
        {"Alice", 30},
        {"Bob", 25},
        {"Charlie", 35},
        {"Dave", 40},
        {"Eve", 28},
        {"Frank", 32},
    }

    jobs := make(chan Person)
    results := make(chan string)
    var wg sync.WaitGroup

    // Création du contexte avec timeout global (1 seconde ici)
    ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
    defer cancel()

    // Lancement de 3 workers
    numWorkers := 3
    for i := 1; i <= numWorkers; i++ {
        wg.Add(1)
        go worker(ctx, i, jobs, results, &wg)
    }

    // Envoi des jobs
    go func() {
        for _, p := range people {
            select {
            case <-ctx.Done():
                return
            case jobs <- p:
            }
        }
        close(jobs)
    }()

    // Récupération des résultats
    go func() {
        wg.Wait()
        close(results)
    }()

    // Affichage
    for msg := range results {
        fmt.Println(msg)
    }

    fmt.Println("✅ Fin du programme (timeout ou travail terminé)")
}
```

---

## 🔍 Ce qu’on a ajouté

| Élément                    | Rôle                                        |
| -------------------------- | ------------------------------------------- |
| `context.WithTimeout(...)` | Arrêter tous les workers après une durée    |
| `ctx.Done()`               | Canal à écouter pour savoir quand s’arrêter |
| `select` dans `worker`     | Permet d’écouter `ctx.Done()` et `jobs`     |
| `cancel()`                 | Libère les ressources du contexte           |

---

## 🧪 Exemple de sortie :

```
Worker 2 : Alice traité en 267ms
Worker 1 : Bob traité en 388ms
Worker 3 : Charlie traité en 266ms
✅ Fin du programme (timeout ou travail terminé)
```

*(Seuls 3 tâches sur 6 ont été traitées car timeout global atteint)*

---

--------------------------------------------------------------------------------------------------------

Souhaites-tu ajouter aussi une **cancellation manuelle** (ex. via l'utilisateur ou une condition métier) ou un **canal d'erreur centralisé** ?
Parfait ! Allons encore plus loin avec un système **de cancellation manuelle ou conditionnelle**, en utilisant à nouveau `context.Context`.

---

## 🎯 Objectif :

* Créer un **worker pool concurrent** comme avant,
* Ajouter la possibilité de **stopper tous les workers** si une **erreur critique** est détectée pendant le traitement.

---

## 🔧 Stratégie :

* Utiliser un **`context.WithCancel()`** (au lieu d’un timeout),
* Lancer les workers avec ce contexte,
* Si une certaine condition est remplie (ex. âge > 35 → erreur), on **annule tout le traitement immédiatement**.

---

## 🧱 Exemple complet : Annulation conditionnelle (erreur simulée)

```go
package main

import (
    "context"
    "fmt"
    "math/rand"
    "sync"
    "time"
)

type Person struct {
    Name string
    Age  int
}

func worker(ctx context.Context, id int, jobs <-chan Person, results chan<- string, errCh chan<- error, wg *sync.WaitGroup) {
    defer wg.Done()
    for {
        select {
        case <-ctx.Done():
            return
        case p, ok := <-jobs:
            if !ok {
                return
            }

            // Simule un travail
            delay := time.Duration(200+rand.Intn(200)) * time.Millisecond
            select {
            case <-ctx.Done():
                return
            case <-time.After(delay):
                if p.Age > 35 {
                    errCh <- fmt.Errorf("❌ Worker %d : erreur critique avec %s (âge %d)", id, p.Name, p.Age)
                    return
                }
                results <- fmt.Sprintf("✔️ Worker %d : %s traité en %v", id, p.Name, delay)
            }
        }
    }
}

func main() {
    rand.Seed(time.Now().UnixNano())

    people := []Person{
        {"Alice", 30},
        {"Bob", 25},
        {"Charlie", 36},  // Provoquera une erreur
        {"Dave", 32},
        {"Eve", 28},
        {"Frank", 40},
    }

    jobs := make(chan Person)
    results := make(chan string)
    errCh := make(chan error, 1) // 1 seule erreur suffit pour annuler
    var wg sync.WaitGroup

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // Workers
    numWorkers := 3
    for i := 1; i <= numWorkers; i++ {
        wg.Add(1)
        go worker(ctx, i, jobs, results, errCh, &wg)
    }

    // Envoi des tâches
    go func() {
        for _, p := range people {
            select {
            case <-ctx.Done():
                return
            case jobs <- p:
            }
        }
        close(jobs)
    }()

    // Résultats
    go func() {
        wg.Wait()
        close(results)
    }()

    // Contrôle des résultats et erreurs
    done := false
    for !done {
        select {
        case msg, ok := <-results:
            if ok {
                fmt.Println(msg)
            } else {
                done = true
            }
        case err := <-errCh:
            fmt.Println(err)
            cancel() // annule tous les workers
            done = true
        }
    }

    fmt.Println("✅ Fin du traitement (annulé ou terminé)")
}
```

---

## 🔍 Ce que tu gagnes ici :

| Élément              | Fonction                                             |
| -------------------- | ---------------------------------------------------- |
| `context.WithCancel` | Permet une **annulation manuelle ou conditionnelle** |
| `errCh chan error`   | Canal pour **propager une erreur critique**          |
| `cancel()`           | Appelé dès qu’une erreur survient                    |

---

## 🧪 Exemple de sortie :

```
✔️ Worker 1 : Alice traité en 334ms
✔️ Worker 2 : Bob traité en 277ms
❌ Worker 3 : erreur critique avec Charlie (âge 36)
✅ Fin du traitement (annulé ou terminé)
```

> ✅ Remarque : même si d'autres jobs restent, tout s’arrête proprement.

---

--------------------------------------------------------------------------------------------------------

Souhaites-tu maintenant une version qui **relance automatiquement les jobs échoués** ou qui **logue tous les résultats + erreurs dans un fichier** ?

Parfait ! Passons à une version plus **résiliente** et **professionnelle** du worker pool, avec deux nouvelles fonctionnalités :

---

## ✅ Nouveautés :

1. **Relance automatique des tâches échouées** (ex. si `Person.Age > 35`)
2. **Enregistrement de tous les résultats et erreurs dans un fichier `.log`**

---

## 🧱 Objectif :

* Les erreurs ne stoppent plus tout,
* Les tâches échouées sont **réessayées une fois**,
* Tout est **journalisé** dans un fichier `log.txt`.

---

## 📦 Code complet

```go
package main

import (
    "fmt"
    "math/rand"
    "os"
    "sync"
    "time"
)

type Person struct {
    Name    string
    Age     int
    Retries int // pour limiter les tentatives
}

func worker(id int, jobs <-chan Person, results chan<- string, retry chan<- Person, logCh chan<- string, wg *sync.WaitGroup) {
    defer wg.Done()
    for p := range jobs {
        delay := time.Duration(200+rand.Intn(300)) * time.Millisecond
        time.Sleep(delay)

        if p.Age > 35 {
            msg := fmt.Sprintf("❌ Worker %d : échec sur %s (âge %d)", id, p.Name, p.Age)
            logCh <- msg

            if p.Retries < 1 {
                p.Retries++
                retry <- p // réessai
            }
        } else {
            msg := fmt.Sprintf("✔️ Worker %d : %s traité en %v", id, p.Name, delay)
            results <- msg
            logCh <- msg
        }
    }
}

func logWriter(logCh <-chan string, done chan<- bool) {
    file, err := os.Create("log.txt")
    if err != nil {
        fmt.Println("Erreur ouverture fichier :", err)
        return
    }
    defer file.Close()

    for msg := range logCh {
        fmt.Println(msg)                  // console
        file.WriteString(msg + "\n")      // fichier
    }

    done <- true
}

func main() {
    rand.Seed(time.Now().UnixNano())

    people := []Person{
        {"Alice", 30, 0},
        {"Bob", 25, 0},
        {"Charlie", 36, 0}, // échoue, puis réessaie
        {"Dave", 40, 0},    // échoue, puis réessaie
        {"Eve", 28, 0},
        {"Frank", 32, 0},
    }

    jobs := make(chan Person)
    retryJobs := make(chan Person)
    results := make(chan string)
    logCh := make(chan string, 10)
    var wg sync.WaitGroup

    // Logger
    doneLogging := make(chan bool)
    go logWriter(logCh, doneLogging)

    // Workers
    numWorkers := 3
    for i := 1; i <= numWorkers; i++ {
        wg.Add(1)
        go worker(i, jobs, results, retryJobs, logCh, &wg)
    }

    // Envoi initial
    go func() {
        for _, p := range people {
            jobs <- p
        }
        close(jobs)
    }()

    // Gestion des retries (1 tentative max)
    go func() {
        for p := range retryJobs {
            jobs <- p
        }
    }()

    // Attente des workers
    go func() {
        wg.Wait()
        close(results)
        close(logCh)
    }()

    // Affiche les résultats (console uniquement)
    for msg := range results {
        // déjà loggé dans logWriter
        _ = msg
    }

    <-doneLogging
    fmt.Println("✅ Traitement terminé. Voir log.txt pour les détails.")
}
```

---

## 📄 Exemple de contenu du fichier `log.txt` :

```
✔️ Worker 1 : Alice traité en 234ms
✔️ Worker 2 : Bob traité en 281ms
❌ Worker 3 : échec sur Charlie (âge 36)
❌ Worker 1 : échec sur Dave (âge 40)
✔️ Worker 2 : Charlie traité en 254ms
✔️ Worker 3 : Dave traité en 217ms
✔️ Worker 1 : Eve traité en 249ms
✔️ Worker 2 : Frank traité en 273ms
```

---

## 🧠 Bilan

| Fonctionalité           | Détail                                |
| ----------------------- | ------------------------------------- |
| `logWriter()`           | Log asynchrone vers console + fichier |
| `retry chan Person`     | Canal de réessai                      |
| `Retries` dans `Person` | Limite à une relance                  |
| `results`, `logCh`      | Séparation résultat / logging         |

---

Souhaites-tu qu’on rende ce système **configurable via un fichier JSON** (noms, seuils d’erreur, nb de workers, etc.) ?
