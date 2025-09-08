Le système de **goroutines** en **Go (Golang)** est un mécanisme léger et puissant pour gérer la **concurrence**. Voici une explication claire et progressive :

---

### 🔹 Qu'est-ce qu'une goroutine ?

Une **goroutine** est une fonction qui s’exécute **concurrentement** avec d’autres goroutines dans le même programme.
Elles sont gérées par le **runtime de Go**, et non directement par le système d’exploitation.

---

### 🔹 Syntaxe de base

Pour lancer une goroutine :

```go
go maFonction()
```

Cela lance `maFonction()` **en parallèle** du reste du programme.

---

### 🔹 Caractéristiques principales

* **Très légères** : chaque goroutine utilise quelques kilo-octets de mémoire (beaucoup moins qu’un thread classique).
* **Multiplexées** : le runtime de Go répartit automatiquement les goroutines sur les threads disponibles du système.
* **Pas de gestion explicite des threads** : tout est abstrait, ce qui simplifie énormément la programmation concurrente.

---

### 🔹 Exemple simple

```go
package main

import (
    "fmt"
    "time"
)

func direBonjour() {
    for i := 0; i < 5; i++ {
        fmt.Println("Bonjour")
        time.Sleep(100 * time.Millisecond)
    }
}

func main() {
    go direBonjour() // lance direBonjour en tant que goroutine
    for i := 0; i < 5; i++ {
        fmt.Println("Salut")
        time.Sleep(100 * time.Millisecond)
    }
}
```

**Sortie possible :**

```
Salut
Bonjour
Salut
Bonjour
...
```

Les lignes "Salut" et "Bonjour" s'entremêlent, car les deux fonctions tournent **en parallèle**.

---

### 🔹 Synchronisation

Comme les goroutines peuvent s’exécuter en même temps, tu dois parfois **synchroniser** leur exécution.

Go propose plusieurs outils :

* **Channels** : pour communiquer entre goroutines.
* **WaitGroups** : pour attendre la fin de plusieurs goroutines.

**Exemple avec `sync.WaitGroup`** :

```go
import (
    "fmt"
    "sync"
)

func travail(wg *sync.WaitGroup) {
    fmt.Println("Travail en cours...")
    wg.Done() // signale que la goroutine est terminée
}

func main() {
    var wg sync.WaitGroup
    wg.Add(1) // on attend une goroutine
    go travail(&wg)
    wg.Wait() // on attend la fin
    fmt.Println("Fini")
}
```

---

### 🔹 Avantages des goroutines

* Haute performance : peut en lancer des **milliers** sans gros coût mémoire.
* Simplicité de syntaxe.
* Intégration fluide avec les channels pour la communication.

---
