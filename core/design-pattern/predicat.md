Le **pattern prédicat** (ou *predicate pattern*) en Go est un **patron de conception fonctionnel** qui permet de **filtrer ou sélectionner des données** à partir d’une condition logique exprimée sous forme de fonction. Ce pattern repose sur l’idée de passer une fonction (le prédicat) qui retourne un booléen pour décider si un élément correspond à un critère donné.

---

### 🔍 Définition

Un **prédicat** est simplement une fonction qui prend un élément en entrée et retourne un `bool` :

```go
type Predicate[T any] func(T) bool
```

---

### 📌 Exemple simple : filtrer une liste

```go
package main

import (
    "fmt"
)

// Définition du type prédicat générique
type Predicate[T any] func(T) bool

// Fonction Filter qui applique le prédicat
func Filter[T any](items []T, pred Predicate[T]) []T {
    var result []T
    for _, item := range items {
        if pred(item) {
            result = append(result, item)
        }
    }
    return result
}

func main() {
    numbers := []int{1, 2, 3, 4, 5, 6}

    // Définir un prédicat pour les nombres pairs
    isEven := func(n int) bool {
        return n%2 == 0
    }

    // Utiliser Filter avec le prédicat
    evenNumbers := Filter(numbers, isEven)

    fmt.Println("Nombres pairs :", evenNumbers)
}
```

### ✅ Résultat

```
Nombres pairs : [2 4 6]
```

---

### 🧠 Avantages du pattern prédicat

* **Lisibilité** : Les intentions de filtrage sont claires et réutilisables.
* **Réutilisabilité** : Les fonctions prédicat peuvent être combinées ou utilisées dans différents contextes.
* **Flexibilité** : On peut changer la logique métier en injectant d'autres fonctions.

---

### 📚 Utilisation courante

* Filtrage en mémoire (collections, slices)
* Conditions pour des requêtes dynamiques (par exemple, en combinaison avec Squirrel ou des ORM)
* Mocks de règles métiers dans des tests unitaires

---

