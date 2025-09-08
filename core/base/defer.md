En Go, le mot-clé `defer` permet de différer l'exécution d'une fonction jusqu'à la fin de la fonction englobante, c’est-à-dire juste avant qu’elle ne retourne.

### Syntaxe

```go
defer fonction()
```

### Comportement

* Les appels `defer` sont empilés : le dernier `defer` déclaré est exécuté en premier (LIFO – Last In, First Out).
* Utilisé souvent pour libérer des ressources : fichiers, connexions, verrous, etc.

### Exemple simple

```go
package main

import "fmt"

func main() {
    fmt.Println("Début")

    defer fmt.Println("Déféré : 1")
    defer fmt.Println("Déféré : 2")

    fmt.Println("Fin")
}
```

**Sortie :**

```
Début
Fin
Déféré : 2
Déféré : 1
```

### Exemple classique : fermeture de fichier

```go
f, err := os.Open("monfichier.txt")
if err != nil {
    log.Fatal(err)
}
defer f.Close() // garantit que le fichier sera fermé à la fin de la fonction
```

### À savoir

* Les arguments passés à la fonction `defer` sont évalués **immédiatement**, mais l’appel à la fonction est différé.

```go
func test() {
    x := 10
    defer fmt.Println("x =", x)
    x = 20
}
```

**Sortie :** `x = 10` (car `x` est évalué au moment du `defer`, pas au moment de l’exécution)

Souhaites-tu un exemple avec gestion d'erreur ou un contexte plus spécifique ?

Parfait, voici un exemple complet avec gestion d'erreur et utilisation de `defer` dans un contexte typique : lecture de fichier.

---

### 🔧 **Exemple : lecture d’un fichier avec gestion d’erreur et `defer`**

```go
package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
)

func lireFichier(nom string) {
    // Ouverture du fichier
    f, err := os.Open(nom)
    if err != nil {
        log.Fatalf("Erreur d'ouverture du fichier : %v", err)
    }

    // Le fichier sera fermé quoi qu’il arrive (erreur ou non)
    defer func() {
        if err := f.Close(); err != nil {
            log.Printf("Erreur lors de la fermeture du fichier : %v", err)
        }
    }()

    // Lecture ligne par ligne
    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        fmt.Println(scanner.Text())
    }

    if err := scanner.Err(); err != nil {
        log.Printf("Erreur de lecture : %v", err)
    }
}

func main() {
    lireFichier("exemple.txt")
}
```

---

### ✅ Ce que fait `defer` ici :

1. **`defer f.Close()`** est appelé **après** que le fichier est ouvert avec succès.
2. Même si la fonction rencontre une erreur plus loin (ex : `scanner.Err()`), le fichier sera **fermé automatiquement**.
3. On encapsule `f.Close()` dans une fonction anonyme avec un `if err != nil` pour capturer aussi les erreurs de fermeture.

---

Souhaites-tu un exemple similaire pour les connexions à une base de données ou avec `panic/recover` ?

Très bien ! Voici un exemple avec `panic`, `recover` et `defer`, qui montre comment gérer proprement les erreurs inattendues dans une fonction Go — typiquement utile dans les serveurs, traitements critiques, etc.

---

### ⚠️ Exemple : `panic`, `recover` et `defer` pour sécuriser une fonction

```go
package main

import "fmt"

func operationRisquee() {
    // Ce defer va intercepter toute panique
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("⚠️ Récupération après un panic :", r)
        }
    }()

    fmt.Println("Début de l'opération")
    
    // Quelque chose de critique provoque un panic
    panic("Quelque chose s'est mal passé !")

    fmt.Println("Fin de l'opération (jamais atteinte)")
}

func main() {
    operationRisquee()
    fmt.Println("Programme terminé proprement")
}
```

---

### 🧠 Ce qui se passe :

* `panic(...)` interrompt immédiatement l'exécution normale.
* Grâce à `defer`, la fonction anonyme avec `recover()` est appelée juste avant la sortie de `operationRisquee`.
* `recover()` intercepte le panic, permettant au programme de **continuer sans crasher**.
* Sans `recover()`, le programme s'arrêterait brutalement.

---

### 🧪 Résultat attendu :

```
Début de l'opération
⚠️ Récupération après un panic : Quelque chose s'est mal passé !
Programme terminé proprement
```

---
