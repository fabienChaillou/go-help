Le billet de blog ["Go Concurrency Patterns: Pipelines and Cancellation"](https://go.dev/blog/pipelines) introduit deux concepts essentiels de la programmation concurrente en Go : les **pipelines** et la **cancellation**. Voici une explication claire et structurée de ces concepts.

---

## 🔄 1. Pipelines en Go

### ✅ Définition

Un **pipeline** est une chaîne de fonctions connectées entre elles via des **canaux (`chan`)**, où la sortie d'une fonction est l'entrée de la suivante.

### 📦 Exemple simple

Imaginons trois étapes :

* **Générateur** : produit une séquence de nombres.
* **Carré** : calcule le carré de chaque nombre.
* **Affichage** : affiche les résultats.

```go
func gen(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        for _, n := range nums {
            out <- n
        }
        close(out)
    }()
    return out
}

func square(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        for n := range in {
            out <- n * n
        }
        close(out)
    }()
    return out
}

func main() {
    for n := range square(gen(2, 3, 4)) {
        fmt.Println(n)
    }
}
```

Chaque étape tourne dans sa propre goroutine, et les données circulent via les canaux. On a ainsi un flux asynchrone, non bloquant, mais contrôlable.

---

## ❌ 2. Cancellation (annulation)

### ✅ Pourquoi ?

Sans cancellation, les goroutines peuvent rester bloquées pour toujours, en attente d’écrire ou de lire sur un canal fermé ou ignoré. Cela entraîne des **goroutines orphelines** et des **fuites mémoire**.

### 🛑 Mécanisme

On utilise un **canal `done <-chan struct{}`** pour signaler l’annulation à tous les composants du pipeline.

### 🧱 Exemple modifié avec cancellation

```go
func gen(done <-chan struct{}, nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for _, n := range nums {
            select {
            case out <- n:
            case <-done:
                return
            }
        }
    }()
    return out
}

func square(done <-chan struct{}, in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for {
            select {
            case n, ok := <-in:
                if !ok {
                    return
                }
                out <- n * n
            case <-done:
                return
            }
        }
    }()
    return out
}
```

Et dans le `main` :

```go
func main() {
    done := make(chan struct{})
    defer close(done)

    in := gen(done, 2, 3, 4)
    out := square(done, in)

    for n := range out {
        fmt.Println(n)
        // on pourrait décider de fermer done ici si une condition est remplie
    }
}
```

---

## 🧠 Résumé

| Concept          | Description                                                                              |
| ---------------- | ---------------------------------------------------------------------------------------- |
| **Pipeline**     | Chaîne de traitements où chaque étape a sa propre goroutine et communique via des canaux |
| **Cancellation** | Permet d'arrêter proprement les goroutines en utilisant un canal `done`                  |

---

### Squaring number
Bien sûr ! Le concept de **"Squaring numbers"** signifie simplement **élever chaque nombre au carré**, c’est-à-dire multiplier un nombre par lui-même.

---

## 📐 Définition mathématique

Pour un nombre `n`, son carré est :

$$
n^2 = n \times n
$$

Par exemple :

* $2^2 = 4$
* $3^2 = 9$
* $4^2 = 16$

---

## 🧑‍💻 En Go (Golang)

Dans un pipeline de Go, **squaring numbers** est souvent une étape de transformation. Voici un exemple de fonction qui prend un canal d'entrée (`in`) et renvoie un canal de sortie (`out`) contenant les carrés :

```go
func square(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        for n := range in {
            out <- n * n
        }
        close(out)
    }()
    return out
}
```

### 🔄 Comment ça marche

* `in <-chan int` : canal en lecture seule contenant des entiers.
* `out chan int` : canal en écriture pour les résultats.
* La goroutine lit chaque nombre depuis `in`, calcule le carré (`n * n`), puis l'envoie dans `out`.

---

## 🧪 Exemple complet

```go
func gen(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        for _, n := range nums {
            out <- n
        }
        close(out)
    }()
    return out
}

func square(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        for n := range in {
            out <- n * n
        }
        close(out)
    }()
    return out
}

func main() {
    for n := range square(gen(2, 3, 4)) {
        fmt.Println(n)
    }
}
```

### 🖨️ Résultat affiché :

```
4
9
16
```

---

## Fan-in Fan-out
Bien sûr ! En Go, les **patterns de concurrence** appelés **fan-out** et **fan-in** sont très utilisés pour paralléliser et agréger des tâches. Voici une explication claire de chacun.

---

## 🔀 1. Fan-Out

### ✅ Définition

**Fan-out** signifie **distribuer un travail sur plusieurs goroutines** pour exécuter des tâches en parallèle.

### 🎯 Objectif

* **Accélérer le traitement** en exécutant plusieurs opérations simultanément.
* Répartir les données entrantes vers plusieurs workers.

### 🧠 Illustration

Imaginons une fonction `square` qui reçoit un flux de nombres, et on veut paralléliser ce calcul :

```go
func square(done <-chan struct{}, in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for {
            select {
            case n, ok := <-in:
                if !ok {
                    return
                }
                out <- n * n
            case <-done:
                return
            }
        }
    }()
    return out
}
```

Dans le `main`, on peut lancer plusieurs workers :

```go
done := make(chan struct{})
defer close(done)

in := gen(2, 3, 4, 5, 6, 7)

// Fan-out : 2 workers lisent en parallèle depuis le même canal `in`
c1 := square(done, in)
c2 := square(done, in)
```

---

## 🔁 2. Fan-In

### ✅ Définition

**Fan-in** signifie **rassembler les résultats de plusieurs goroutines** dans un seul canal de sortie.

### 🎯 Objectif

* **Aggréger les résultats** venant de plusieurs sources concurrentes.
* Centraliser l'écoute des résultats.

### 🧠 Illustration

Voici une fonction `merge` qui fait du fan-in :

```go
func merge(done <-chan struct{}, cs ...<-chan int) <-chan int {
    var wg sync.WaitGroup
    out := make(chan int)

    output := func(c <-chan int) {
        defer wg.Done()
        for {
            select {
            case n, ok := <-c:
                if !ok {
                    return
                }
                out <- n
            case <-done:
                return
            }
        }
    }

    wg.Add(len(cs))
    for _, c := range cs {
        go output(c)
    }

    // Ferme le canal de sortie une fois que tous les workers sont terminés
    go func() {
        wg.Wait()
        close(out)
    }()

    return out
}
```

Dans le `main` :

```go
for n := range merge(done, c1, c2) {
    fmt.Println(n)
}
```

---

## 🧭 Résumé

| Pattern     | Description                                                       | Exemple               |
| ----------- | ----------------------------------------------------------------- | --------------------- |
| **Fan-out** | Plusieurs workers lisent à partir du même canal d'entrée          | `square` en parallèle |
| **Fan-in**  | Un seul canal de sortie agrège les résultats de plusieurs workers | `merge`               |

---

## Stopping short

Bien sûr ! Dans le contexte de la **concurrence en Go**, **"Stopping short"** signifie **arrêter prématurément un pipeline** dès qu’on n’a plus besoin de tous les résultats.

Cela évite :

* de gaspiller des ressources CPU
* d’avoir des goroutines orphelines
* des fuites mémoire (goroutines qui bloquent à attendre des données)

---

## 📍 Cas concret

Imaginons que tu veux juste **le premier résultat valide** d’un pipeline, puis **tout arrêter**.

---

## 🧪 Exemple de problème sans "stopping short"

```go
func gen(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        for _, n := range nums {
            out <- n
        }
        close(out)
    }()
    return out
}

func square(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        for n := range in {
            out <- n * n
        }
        close(out)
    }()
    return out
}

func main() {
    in := gen(2, 3, 4, 5, 6, 7)
    out := square(in)

    fmt.Println(<-out) // on lit juste le premier résultat

    // ❌ Mais les goroutines continuent à tourner !
}
```

---

## ✅ Solution : `done` channel (Stopping short)

Ajoutons un **canal `done`** pour permettre l’annulation :

```go
func gen(done <-chan struct{}, nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for _, n := range nums {
            select {
            case out <- n:
            case <-done:
                return
            }
        }
    }()
    return out
}

func square(done <-chan struct{}, in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for {
            select {
            case n, ok := <-in:
                if !ok {
                    return
                }
                select {
                case out <- n * n:
                case <-done:
                    return
                }
            case <-done:
                return
            }
        }
    }()
    return out
}
```

Et dans le `main` :

```go
func main() {
    done := make(chan struct{})
    defer close(done)

    in := gen(done, 2, 3, 4, 5, 6)
    out := square(done, in)

    fmt.Println(<-out) // ✅ on lit un seul résultat
    // Le defer close(done) arrête tout le pipeline proprement
}
```

---

## 🧠 Pourquoi c’est important ?

* Sans cancellation, des goroutines peuvent rester bloquées.
* Go ne **garbage collect** pas les goroutines actives.
* Bien gérer "stopping short" est crucial pour des pipelines robustes, réactifs et sans fuite.

---

### Examples

Parfait ! Voici comment gérer **"stopping short"** en Go en utilisant le **package `context`**, qui est la **méthode recommandée** dans les applications modernes.

---

## ✅ Pourquoi `context.Context` ?

Le package `context` permet de :

* définir un **timeout** ou une **deadline**
* propager une **annulation** à toutes les goroutines concernées
* éviter d’avoir à gérer manuellement un canal `done`

---

## 🔧 Exemple avec `context.WithCancel`

### 🔁 Pipeline avec contexte

```go
func gen(ctx context.Context, nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for _, n := range nums {
            select {
            case <-ctx.Done():
                return
            case out <- n:
            }
        }
    }()
    return out
}

func square(ctx context.Context, in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for {
            select {
            case <-ctx.Done():
                return
            case n, ok := <-in:
                if !ok {
                    return
                }
                select {
                case out <- n * n:
                case <-ctx.Done():
                    return
                }
            }
        }
    }()
    return out
}
```

### 🧪 Utilisation dans `main`

```go
func main() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel() // 🛑 déclenche l'annulation à la fin

    in := gen(ctx, 2, 3, 4, 5)
    out := square(ctx, in)

    // Lire seulement le premier résultat
    fmt.Println(<-out)

    // ✅ Les goroutines s'arrêtent automatiquement grâce au `context`
}
```

---

## ⏰ Variante avec `context.WithTimeout`

On peut aussi **arrêter automatiquement après un délai** :

```go
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
defer cancel()
```

Cela fonctionne avec les mêmes fonctions `gen` et `square`.

---

## 🧠 Avantages de `context`

| Avantage     | Description                                                        |
| ------------ | ------------------------------------------------------------------ |
| ✅ Centralisé | Toutes les annulations passent par un seul objet `ctx`             |
| ✅ Intégré    | Beaucoup de bibliothèques Go le supportent nativement              |
| ✅ Lisible    | Moins de code personnalisé à gérer (pas de canal `done` à la main) |

---

Bien sûr ! En Go, **explicit cancellation** (annulation explicite) signifie que **tu contrôles manuellement le moment où un pipeline ou une série de goroutines doit s’arrêter**, au lieu d’attendre qu’elles terminent d’elles-mêmes ou qu’un timeout survienne.

---

## ✅ Quand l’utiliser ?

Tu veux arrêter **dès que :**

* tu as obtenu un résultat satisfaisant,
* une erreur est survenue,
* l’utilisateur a annulé l’opération,
* tu ne veux plus traiter le reste des données.

---

## 🧱 Moyens d’implémenter une cancellation explicite

### 1. Avec un **canal `done`** (approche manuelle)

C’est la forme "brute" d'annulation explicite. Tu fermes un canal pour signaler que tout doit s’arrêter.

#### Exemple :

```go
done := make(chan struct{})
defer close(done) // ✋ Annulation explicite

in := gen(done, 1, 2, 3, 4)
out := square(done, in)

fmt.Println(<-out) // on lit juste le premier résultat
```

Dans les fonctions `gen` et `square`, tu vérifies régulièrement :

```go
select {
case <-done:
    return
default:
    // continuer normalement
}
```

---

### 2. Avec `context.WithCancel()` (approche recommandée)

C’est la méthode moderne et préférée en Go. Elle fournit une fonction `cancel()` que tu peux appeler **explicitement**.

#### Exemple :

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel() // 👈 annulation explicite quand la fonction se termine

in := gen(ctx, 1, 2, 3)
out := square(ctx, in)

fmt.Println(<-out) // on lit un seul résultat

cancel() // 🛑 explicitement, on arrête toutes les goroutines qui écoutent ce contexte
```

Et dans `gen` / `square`, tu écoutes :

```go
select {
case <-ctx.Done():
    return
}
```

---

## 🧠 Pourquoi faire une cancellation explicite ?

| Problème évité                       | Comment ?                                                                  |
| ------------------------------------ | -------------------------------------------------------------------------- |
| Goroutines qui tournent dans le vide | On les arrête via `cancel()` ou `close(done)`                              |
| Fuites mémoire                       | Pas de goroutines bloquées inutilement                                     |
| Réactivité                           | Le programme réagit immédiatement dès que tu n’as plus besoin de continuer |

---

## 🔁 Résumé

| Type              | Mécanisme     | Utilisation                       |
| ----------------- | ------------- | --------------------------------- |
| Canal `done`      | `close(done)` | Simple mais bas niveau            |
| `context.Context` | `cancel()`    | Recommandé, flexible, standardisé |

---

### Digesting a tree

Bien sûr ! Le pattern **"Digesting a tree"** en Go, dans le contexte du blog [Go Concurrency Patterns: Pipelines and Cancellation](https://go.dev/blog/pipelines), signifie **parcourir une structure arborescente (comme un système de fichiers) de manière concurrente**, tout en **contrôlant les ressources** (comme le nombre de goroutines) et en gérant l'annulation.

---

## 🌳 Qu’est-ce qu’un arbre ici ?

Un **arbre** est une structure récursive : un **nœud** peut avoir plusieurs **sous-nœuds**. Exemple typique : les **dossiers et fichiers** dans un système de fichiers.

---

## 🧠 Objectif : digérer un arbre

Cela veut dire :

1. Parcourir l’arbre (ex : tous les fichiers d’un répertoire),
2. Lister les fichiers (ou effectuer une action sur chaque nœud),
3. Gérer la **concurrence** (pour aller plus vite),
4. **Limiter les ressources** (éviter d’ouvrir 1000 goroutines),
5. Supporter l’**annulation** (si on veut arrêter à mi-chemin).

---

## 📦 Exemple tiré du blog : `du` (disk usage)

L’outil `du` parcourt un dossier et calcule la taille totale des fichiers. Voici l’idée :

```go
func walkDir(dir string, n *sync.WaitGroup, fileSizes chan<- int64) {
    defer n.Done()
    for _, entry := range entries(dir) {
        if entry.IsDir() {
            n.Add(1)
            subdir := filepath.Join(dir, entry.Name())
            go walkDir(subdir, n, fileSizes)
        } else {
            fileSizes <- entry.Size()
        }
    }
}
```

---

## ✋ Problème : trop de goroutines

Si tu as des milliers de sous-dossiers, tu risques de lancer trop de goroutines. Solution :

### ✅ Limiter le parallélisme avec un **semaphore** (canal tamponné)

```go
var sema = make(chan struct{}, 20) // max 20 goroutines concurrentes

func entries(dir string) []os.FileInfo {
    sema <- struct{}{}           // acquire
    defer func() { <-sema }()    // release
    // lire les fichiers dans dir
}
```

---

## 🚨 Et si on veut arrêter au milieu ? (Digesting + cancellation)

Tu ajoutes un `done` ou un `context.Context` :

```go
func walkDir(ctx context.Context, dir string, n *sync.WaitGroup, fileSizes chan<- int64) {
    defer n.Done()
    for _, entry := range entries(dir) {
        select {
        case <-ctx.Done():
            return
        default:
        }

        if entry.IsDir() {
            n.Add(1)
            subdir := filepath.Join(dir, entry.Name())
            go walkDir(ctx, subdir, n, fileSizes)
        } else {
            select {
            case fileSizes <- entry.Size():
            case <-ctx.Done():
                return
            }
        }
    }
}
```

---

## 🧭 Résumé : Digesting a tree

| Élément                     | Rôle                                                        |
| --------------------------- | ----------------------------------------------------------- |
| Arbre                       | Structure récursive à parcourir (ex : système de fichiers)  |
| Goroutines                  | Chaque sous-dossier peut être exploré en parallèle          |
| Sémaphore (canal tamponné)  | Limite le nombre de goroutines concurrentes                 |
| `done` ou `context.Context` | Permet d’annuler proprement le parcours                     |
| `sync.WaitGroup`            | Attend que tout le traitement soit terminé avant de quitter |

---

## Parallel digestion

Bien sûr ! En Go, **Parallel Digestion** (ou "digestion parallèle") est une extension du pattern **"Digesting a tree"**, où tu parcours une **structure arborescente de façon concurrente** pour **accélérer le traitement** de chaque branche ou nœud.

---

## 🧠 Objectif de Parallel Digestion

* Traiter **chaque sous-partie de l’arbre en parallèle**.
* Utiliser des **goroutines pour explorer plusieurs branches simultanément**.
* Gérer les **ressources**, les **erreurs**, et la **cancellation** proprement.

---

## 📦 Exemple concret : exploration d'un système de fichiers

Imaginons que tu veux calculer la taille totale de tous les fichiers dans un répertoire (comme l’outil `du`), **mais en parallèle**.

---

### 🧱 Étapes typiques de parallel digestion

1. **Lister les sous-dossiers** d’un répertoire.
2. **Lancer une goroutine** pour chaque sous-dossier.
3. **Envoyer les résultats** sur un canal.
4. **Coordonner** avec `sync.WaitGroup`.
5. **Limiter le nombre de goroutines** avec un **sémaphore** (canal tamponné).
6. **Annuler proprement** via `context.Context`.

---

## 🧑‍💻 Code simplifié : Parallel Digestion

```go
func walkDir(ctx context.Context, dir string, wg *sync.WaitGroup, sizes chan<- int64, sema chan struct{}) {
    defer wg.Done()

    select {
    case sema <- struct{}{}: // acquire token
    case <-ctx.Done():
        return
    }
    defer func() { <-sema }() // release token

    entries, err := os.ReadDir(dir)
    if err != nil {
        return
    }

    for _, entry := range entries {
        if entry.IsDir() {
            subdir := filepath.Join(dir, entry.Name())
            wg.Add(1)
            go walkDir(ctx, subdir, wg, sizes, sema)
        } else {
            info, err := entry.Info()
            if err == nil {
                select {
                case sizes <- info.Size():
                case <-ctx.Done():
                    return
                }
            }
        }
    }
}
```

### 🔁 Utilisation :

```go
func main() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    sizes := make(chan int64)
    var wg sync.WaitGroup
    sema := make(chan struct{}, 20) // max 20 concurrent goroutines

    root := "/some/dir"

    wg.Add(1)
    go walkDir(ctx, root, &wg, sizes, sema)

    go func() {
        wg.Wait()
        close(sizes)
    }()

    var total int64
    for size := range sizes {
        total += size
    }

    fmt.Printf("Total size: %.1f MB\n", float64(total)/1e6)
}
```

---

## ⚙️ Résumé de Parallel Digestion

| Composant         | Rôle                                                            |
| ----------------- | --------------------------------------------------------------- |
| `goroutines`      | Parcours parallèle de l’arbre                                   |
| `sync.WaitGroup`  | Synchronise la fin de toutes les explorations                   |
| `chan struct{}`   | Sémaphore pour limiter le nombre de goroutines                  |
| `context.Context` | Permet une annulation propre du traitement                      |
| `chan int64`      | Canal pour transmettre les résultats (ex : tailles de fichiers) |

---

## 📌 Pourquoi c’est utile ?

* C’est **scalable** : les dossiers/projets énormes sont traités efficacement.
* C’est **propre** : les erreurs et annulations sont bien gérées.
* C’est **rapide** : tu exploites tous les cœurs CPU disponibles.

---

## Bounded parallelism

Bien sûr ! Le concept de **Bounded Parallelism** (ou **parallélisme borné**) en Go consiste à **limiter le nombre de goroutines exécutées en parallèle** pour :

* éviter de saturer le CPU ou la mémoire,
* respecter des quotas (ex : nombre maximum de connexions à une API, nombre de fichiers ouverts, etc.),
* garder un **contrôle fin sur la consommation de ressources**.

---

## 📍 Pourquoi pas un parallélisme illimité ?

Même si Go permet de lancer facilement des milliers de goroutines, cela peut poser problème :

* Trop de goroutines = trop de **RAM** utilisée.
* Certaines opérations sont limitées (ex : nombre max de fichiers ouverts).
* Tu risques de créer une **tempête de traitement** difficile à gérer.

---

## ✅ Solution : canal tamponné = **sémaphore**

Tu utilises un **canal de type `chan struct{}` avec capacité fixe** comme **sémaphore** pour limiter le nombre de goroutines **actives en même temps**.

---

## 🧱 Exemple minimal : 5 goroutines max

```go
var sem = make(chan struct{}, 5) // maximum 5 goroutines en même temps

func process(i int) {
    sem <- struct{}{}        // 🛑 bloque si 5 goroutines sont déjà en cours
    go func() {
        defer func() { <-sem }() // ✅ libère une place à la fin

        fmt.Printf("Processing item %d\n", i)
        time.Sleep(1 * time.Second)
    }()
}
```

Et dans ta boucle principale :

```go
for i := 0; i < 20; i++ {
    process(i)
}

// Attendre la fin (sinon le programme se termine avant)
time.Sleep(5 * time.Second)
```

---

## 🧑‍🔧 Avec `sync.WaitGroup` pour attendre proprement

```go
var wg sync.WaitGroup
sem := make(chan struct{}, 5)

for i := 0; i < 20; i++ {
    sem <- struct{}{}
    wg.Add(1)

    go func(i int) {
        defer func() {
            <-sem
            wg.Done()
        }()
        fmt.Printf("Job %d started\n", i)
        time.Sleep(1 * time.Second)
    }(i)
}

wg.Wait()
fmt.Println("All done")
```

---

## 🧠 Résumé : Bounded Parallelism

| Élément                  | Rôle                                   |
| ------------------------ | -------------------------------------- |
| `chan struct{}` tamponné | Sert de **sémaphore**                  |
| `len(sem)`               | Donne le nombre de goroutines en cours |
| `WaitGroup`              | Synchronise la fin des tâches          |

---

## 📌 Cas d’usage typiques

* Explorer un système de fichiers (comme dans `walkDir`)
* Téléchargements concurrents mais limités
* Traitement de requêtes API sans dépasser les quotas
* Compression/conversion de fichiers en parallèle

---

Souhaites-tu un exemple complet avec **context**, **bounded parallelism** et **fan-in/out** combinés ?

### Examples

Parfait ! Voici un **exemple complet en Go** combinant :

* ✅ **Context** pour la cancellation,
* ✅ **Bounded parallelism** pour limiter les ressources,
* ✅ **Fan-out** pour distribuer le travail,
* ✅ **Fan-in** pour collecter les résultats.

---

## 🎯 Objectif du code

> Traiter une liste de fichiers (fictifs ici), en parallèle, mais avec au plus 5 goroutines actives, et s’arrêter proprement si l'utilisateur annule.

---

## 🧑‍💻 Code Go Complet

```go
package main

import (
    "context"
    "fmt"
    "math/rand"
    "sync"
    "time"
)

func main() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // Simulons une liste de fichiers à traiter
    files := []string{"a.txt", "b.txt", "c.txt", "d.txt", "e.txt", "f.txt", "g.txt"}

    results := make(chan string)
    sem := make(chan struct{}, 5) // max 5 goroutines
    var wg sync.WaitGroup

    // Fan-out: une goroutine par fichier
    for _, file := range files {
        select {
        case <-ctx.Done():
            break
        case sem <- struct{}{}:
            wg.Add(1)
            go func(f string) {
                defer func() {
                    <-sem
                    wg.Done()
                }()

                result, err := processFile(ctx, f)
                if err != nil {
                    fmt.Printf("Error processing %s: %v\n", f, err)
                    return
                }

                select {
                case results <- result:
                case <-ctx.Done():
                    return
                }
            }(file)
        }
    }

    // Fan-in: collecteur des résultats
    go func() {
        wg.Wait()
        close(results)
    }()

    // Lire les résultats, ou arrêter à une condition
    for r := range results {
        fmt.Println("Résultat:", r)

        if r == "b.txt traité" {
            fmt.Println("Condition atteinte, annulation.")
            cancel() // 🛑 annulation explicite
        }
    }

    fmt.Println("✅ Tous les traitements terminés ou annulés.")
}

// Simule un traitement de fichier
func processFile(ctx context.Context, name string) (string, error) {
    delay := time.Duration(rand.Intn(1000)+500) * time.Millisecond
    select {
    case <-time.After(delay): // simule le travail
        return name + " traité", nil
    case <-ctx.Done():
        return "", ctx.Err()
    }
}
```

---

## 📌 Ce que ce code montre

| Pattern                 | Où il est utilisé                                |
| ----------------------- | ------------------------------------------------ |
| **Context**             | `context.WithCancel()` et `ctx.Done()`           |
| **Fan-out**             | Une goroutine par fichier                        |
| **Fan-in**              | Canal `results` collectant tous les résultats    |
| **Bounded parallelism** | Canal `sem` limitant à 5 goroutines concurrentes |

---

## 🧪 Résultat attendu (aléatoire)

```
Résultat: c.txt traité
Résultat: b.txt traité
Condition atteinte, annulation.
Résultat: e.txt traité
✅ Tous les traitements terminés ou annulés.
```

---

