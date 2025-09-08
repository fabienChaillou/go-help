Voici les **grands principes de la programmation en Go (Golang)**, conçus pour la simplicité, la performance et la lisibilité :

---

### 1. **Simplicité et minimalisme**

Go est volontairement minimaliste :

* Peu de concepts : pas d’héritage, pas de génériques (introduits seulement à partir de Go 1.18).
* Syntaxe claire et concise.
* Peu de magie : ce que tu vois, c’est ce que tu obtiens (WYSIWYG).

---

### 2. **Types statiques et compilation rapide**

* Go est **statiquement typé** : les erreurs sont détectées à la compilation.
* La compilation est **rapide**, ce qui améliore le cycle de développement.

---

### 3. **Goroutines et concurrence native**

* Go propose une gestion de la **concurrence légère** avec :

  * `goroutine` pour exécuter des fonctions de manière concurrente.
  * `channel` pour communiquer entre goroutines.
* Inspiré du modèle CSP (Communicating Sequential Processes).

---

### 4. **Gestion explicite des erreurs**

* Pas d'exceptions classiques.
* Les erreurs sont **traitées comme des valeurs** retournées (`error`).
* Encourage à vérifier chaque erreur, ce qui rend le code plus robuste.

```go
res, err := doSomething()
if err != nil {
    // gérer l'erreur
}
```

---

### 5. **Composition via les interfaces**

* Pas d’héritage : Go favorise la **composition** (has-a) plutôt que l’héritage (is-a).
* Interfaces **implicites** : une struct implémente une interface si elle possède les méthodes requises.

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}
```

---

### 6. **Packages et organisation claire**

* Tout est organisé en **packages**, chaque fichier commence par `package nom`.
* La convention : `main` pour les programmes exécutables, les autres pour les bibliothèques.

---

### 7. **Formatage automatique avec `gofmt`**

* Code toujours formaté de manière cohérente avec `gofmt` ou `go fmt`.
* Supprime les débats sur le style de code : tout le monde écrit du Go "de la même façon".

---

### 8. **Outils intégrés puissants**

Go fournit des outils natifs :

* `go build`, `go test`, `go fmt`, `go vet`, `go run`, etc.
* Un système de modules (`go mod`) pour la gestion des dépendances.

---

### 9. **Garbage collection**

* Go est **garbage-collected** : il gère automatiquement la mémoire, mais avec un très bon niveau de performance.

---

### 10. **Interopérabilité C et faible empreinte**

* Peut interagir avec du C via `cgo` si nécessaire.
* Parfait pour des microservices ou des outils CLI performants.

---

Un **langage impératif** est un type de langage de programmation dans lequel les instructions sont données sous forme de **séquences d'ordres** que l'ordinateur doit exécuter pas à pas. Il s'agit d'une approche où l'on **décrit *comment* faire** quelque chose, en manipulant directement l'état du programme à l'aide de **variables, boucles, conditions, affectations**, etc.

### Caractéristiques principales :

* Utilisation d’**instructions** pour modifier l’état du programme.
* **Contrôle explicite du flux d’exécution** (avec des `if`, `for`, `while`, etc.).
* Programmation centrée sur **le "comment"** (par opposition à la programmation déclarative qui est centrée sur le "quoi").

### Exemples de langages impératifs :

* **C**
* **Java**
* **Python** (majoritairement impératif, même s’il supporte d’autres paradigmes)
* **Go**
* **Pascal**

### Exemple simple (en C) :

```c
int somme = 0;
for (int i = 1; i <= 10; i++) {
    somme += i;
}
```

Ce code donne à l’ordinateur des instructions précises : initialiser une variable, faire une boucle, ajouter des valeurs.

Voici un **exemple simple en Go** qui illustre le style **impératif** : on calcule la somme des nombres de 1 à 10.

### Exemple impératif en Go :

```go
package main

import "fmt"

func main() {
    sum := 0
    for i := 1; i <= 10; i++ {
        sum += i
    }
    fmt.Println("Somme de 1 à 10 :", sum)
}
```

### Explication :

* `sum := 0` : on initialise une variable (changement d’état).
* `for i := 1; i <= 10; i++` : une boucle contrôle le flux.
* `sum += i` : on modifie explicitement l’état (la valeur de `sum`).
* `fmt.Println(...)` : on affiche le résultat.

On dit que c’est **impératif** parce que le programmeur contrôle précisément **comment** le résultat est obtenu, étape par étape.

## example version fonctionnel

Très bien ! Go n'est pas un langage purement fonctionnel, mais on peut tout de même adopter un **style fonctionnel** en minimisant les effets de bord, en **évitant les variables mutables** et en **utilisant des fonctions pures** (qui ne modifient pas d’état extérieur).

Voici une **version plus fonctionnelle** du même calcul (somme de 1 à 10) :

```go
package main

import (
    "fmt"
)

func sumRange(start, end int) int {
    if start > end {
        return 0
    }
    return start + sumRange(start+1, end)
}

func main() {
    result := sumRange(1, 10)
    fmt.Println("Somme de 1 à 10 :", result)
}
```
