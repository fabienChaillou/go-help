La **programmation fonctionnelle** est un **paradigme de programmation** qui met l'accent sur l’utilisation de **fonctions pures**, l’**immuabilité**, et l’**absence d’effets de bord**.

Bien que Go ne soit **pas un langage fonctionnel** à la base (il est impérativement structuré), **il permet d’adopter plusieurs concepts fonctionnels** grâce à son support des fonctions de premier ordre, des closures et de l’anonymat fonctionnel.

---

## 🔹 Concepts fonctionnels en Go

### 1. **Fonctions de premier ordre**

> Les fonctions sont des **valeurs** : elles peuvent être stockées dans des variables, passées en argument, ou retournées.

```go
func doubler(x int) int {
    return x * 2
}

func appliquer(f func(int) int, v int) int {
    return f(v)
}

func main() {
    fmt.Println(appliquer(doubler, 5)) // 10
}
```

---

### 2. **Fonctions anonymes et closures**

```go
carre := func(x int) int {
    return x * x
}
fmt.Println(carre(4)) // 16
```

Les **closures** permettent de capturer un état tout en restant dans un style fonctionnel.

---

### 3. **Fonctions retournant des fonctions**

```go
func multiplier(facteur int) func(int) int {
    return func(x int) int {
        return x * facteur
    }
}

fois3 := multiplier(3)
fmt.Println(fois3(4)) // 12
```

---

### 4. **Immuabilité**

> Go n’impose pas l’immuabilité, mais tu peux l'appliquer par discipline.

Exemple :

```go
func ajouter(x int, liste []int) []int {
    nouvelleListe := append([]int{}, liste...) // copie
    nouvelleListe = append(nouvelleListe, x)
    return nouvelleListe
}
```

Ici, on **ne modifie pas l'entrée** `liste`, on crée une **nouvelle version**.

---

### 5. **Pas d’effets de bord**

Une fonction pure :

* Ne modifie pas l’état externe.
* Ne dépend pas de l’état externe.
* Retourne toujours le même résultat pour les mêmes entrées.

```go
func somme(a, b int) int {
    return a + b // fonction pure
}
```

---

## 🔹 Limites de la programmation fonctionnelle en Go

* Pas de **map/filter/reduce** intégrés comme en Python ou JavaScript (mais on peut les écrire).
* Pas de **currying automatique**, **pattern matching**, ou **monades**.
* Go préfère la simplicité et la lisibilité sur les abstractions complexes.

---

## 🔹 Exemple : `filter` fonctionnel

```go
func filter(liste []int, predicat func(int) bool) []int {
    resultat := []int{}
    for _, v := range liste {
        if predicat(v) {
            resultat = append(resultat, v)
        }
    }
    return resultat
}

func main() {
    nombres := []int{1, 2, 3, 4, 5}
    pairs := filter(nombres, func(x int) bool { return x%2 == 0 })
    fmt.Println(pairs) // [2 4]
}
```

---

## 🔹 En résumé

| Concept                | Supporté en Go ?    | Exemple                      |
| ---------------------- | ------------------- | ---------------------------- |
| Fonctions de 1er ordre | ✅                   | `appliquer(f, x)`            |
| Fonctions anonymes     | ✅                   | `func(x int) int { ... }`    |
| Closures               | ✅                   | `capturer une variable`      |
| Immuabilité            | 🟡 (manuelle)       | `append([]T{}, original...)` |
| Effets de bord évités  | 🟡 (par discipline) | fonctions pures              |
| map/filter/reduce      | 🟡 (implémentables) | `filter(slice, func)`        |

---


### Caractéristiques fonctionnelles ici :

* La fonction `sumRange` est **récursive** (pas de boucle).
* Elle est **pure** : même entrée → même sortie, pas de modification d'état global.
* Il n'y a **pas de variable mutable** comme `sum`.

Ce style est plus **déclaratif** : on décrit *ce qu’on veut* (`sumRange`) plutôt que comment le calcul se fait étape par étape.

Souhaitez-vous un exemple encore plus poussé avec des fonctions d'ordre supérieur ?

Parfait ! Voici un exemple en Go utilisant des **fonctions d’ordre supérieur** (fonctions qui prennent d'autres fonctions en paramètre ou en retournent), dans un style plus **fonctionnel**, bien que Go ne soit pas optimisé pour ça.

Nous allons :

1. Créer une fonction `map` pour transformer des valeurs.
2. Créer une fonction `reduce` pour les agréger.
3. Utiliser ces fonctions pour calculer la somme des nombres de 1 à 10.

---

### 💡 Exemple fonctionnel en Go avec `map` et `reduce` :

```go
package main

import (
    "fmt"
)

// map applique une fonction à chaque élément d'une slice
func mapSlice[T any, R any](input []T, f func(T) R) []R {
    result := make([]R, len(input))
    for i, v := range input {
        result[i] = f(v)
    }
    return result
}

// reduce applique une fonction cumulatrice à chaque élément
func reduce[T any, R any](input []T, initial R, f func(R, T) R) R {
    acc := initial
    for _, v := range input {
        acc = f(acc, v)
    }
    return acc
}

func main() {
    numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

    // Map inutile ici, mais on aurait pu doubler chaque nombre par exemple
    doubled := mapSlice(numbers, func(x int) int {
        return x // ou: return x * 2
    })

    sum := reduce(doubled, 0, func(acc, x int) int {
        return acc + x
    })

    fmt.Println("Somme de 1 à 10 :", sum)
}
```

---

### 🧠 Points clés :

* `mapSlice` et `reduce` sont inspirées de la programmation fonctionnelle.
* Pas de boucles ou mutations visibles dans `main`.
* On décompose le calcul en **transformations pures** sur les données.

Parfait ! Structurer ce code en modules/pakkets (fichiers Go séparés) améliore la **réutilisabilité** et **lisibilité**. Voici une organisation simple en deux fichiers :

---

### 📁 Arborescence du projet

```
myproject/
├── main.go
└── functional/
    └── functional.go
```

---

### ✅ `functional/functional.go`

Ce fichier contient les fonctions génériques `Map` et `Reduce`.

```go
package functional

// Map applique une fonction à chaque élément d'une slice.
func Map[T any, R any](input []T, f func(T) R) []R {
	result := make([]R, len(input))
	for i, v := range input {
		result[i] = f(v)
	}
	return result
}

// Reduce applique une fonction d'accumulation à une slice.
func Reduce[T any, R any](input []T, initial R, f func(R, T) R) R {
	acc := initial
	for _, v := range input {
		acc = f(acc, v)
	}
	return acc
}
```

---

### ✅ `main.go`

Ce fichier importe et utilise les fonctions de `functional`.

```go
package main

import (
	"fmt"
	"myproject/functional" // Remplacez "myproject" par le bon nom de module si besoin
)

func main() {
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	doubled := functional.Map(numbers, func(x int) int {
		return x // ou x * 2 pour transformer
	})

	sum := functional.Reduce(doubled, 0, func(acc, x int) int {
		return acc + x
	})

	fmt.Println("Somme de 1 à 10 :", sum)
}
```

---

### 🛠️ Initialisation du module

Si ce n’est pas encore fait, dans le dossier racine (`myproject/`), exécute :

```bash
go mod init myproject
go mod tidy
```

Cela permettra à `main.go` de reconnaître le package `functional`.

Souhaitez-vous que je vous fournisse un exemple avec un test unitaire de ces fonctions aussi ?

Très bon choix ! Ajouter des **tests unitaires** renforce la fiabilité de votre code. Voici comment tester les fonctions `Map` et `Reduce` du package `functional`.

---

### 📁 Nouvelle arborescence

```
myproject/
├── go.mod
├── main.go
└── functional/
    ├── functional.go
    └── functional_test.go
```

---

### ✅ `functional/functional_test.go`

```go
package functional

import (
	"testing"
)

func TestMap(t *testing.T) {
	input := []int{1, 2, 3}
	expected := []int{2, 4, 6}
	result := Map(input, func(x int) int {
		return x * 2
	})

	for i := range result {
		if result[i] != expected[i] {
			t.Errorf("Map failed at index %d: got %d, want %d", i, result[i], expected[i])
		}
	}
}

func TestReduce(t *testing.T) {
	input := []int{1, 2, 3, 4}
	expected := 10
	result := Reduce(input, 0, func(acc, x int) int {
		return acc + x
	})

	if result != expected {
		t.Errorf("Reduce failed: got %d, want %d", result, expected)
	}
}
```

---

### ▶️ Pour lancer les tests

Dans le terminal, placez-vous à la racine (`myproject/`) et exécutez :

```bash
go test ./functional
```

Vous verrez :

```
ok  	myproject/functional	0.XXXs
```

---

### 🧪 Bonus : testez aussi avec d’autres types (ex. string)

Si vous le souhaitez, je peux vous montrer comment tester `Map[string, int]` ou d'autres combinaisons avec des types génériques.

Souhaitez-vous un exemple avec un type personnalisé (struct, string, etc.) ?

Excellent ! Utiliser les **génériques** avec des **types personnalisés** montre toute la puissance de `Map` et `Reduce`. Voici un exemple avec une `struct` personnalisée et des tests.

---

### ✅ Exemple avec une struct : `Person`

```go
type Person struct {
	Name string
	Age  int
}
```

On va :

* Mapper une slice de `Person` pour extraire les noms (`[]string`).
* Reduire pour calculer l'âge total (`int`).

---

### ✅ Ajoutez ce test à `functional/functional_test.go`

```go
func TestMapStructToString(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	people := []Person{
		{"Alice", 30},
		{"Bob", 25},
		{"Charlie", 40},
	}

	expectedNames := []string{"Alice", "Bob", "Charlie"}
	names := Map(people, func(p Person) string {
		return p.Name
	})

	for i := range names {
		if names[i] != expectedNames[i] {
			t.Errorf("Map struct failed at index %d: got %s, want %s", i, names[i], expectedNames[i])
		}
	}
}

func TestReduceStructToInt(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	people := []Person{
		{"Alice", 30},
		{"Bob", 25},
		{"Charlie", 40},
	}

	expectedSum := 95
	totalAge := Reduce(people, 0, func(acc int, p Person) int {
		return acc + p.Age
	})

	if totalAge != expectedSum {
		t.Errorf("Reduce struct failed: got %d, want %d", totalAge, expectedSum)
	}
}
```

---

### ▶️ Lancer les tests :

```bash
go test ./functional
```

---

### 🧠 Ce que ça montre :

* `Map` peut convertir des `[]Person` en `[]string`.
* `Reduce` peut agréger des infos comme l’âge.
* Le code est générique, réutilisable et très lisible.

Souhaitez-vous que je vous propose un package `collections` complet avec tri, filtre (`Filter`) et autres utilitaires fonctionnels ?

Parfait ! Créons un petit **package `collections`** inspiré des langages fonctionnels comme JavaScript (lodash), Rust (iterators) ou Haskell, avec les fonctions suivantes :

### ✅ Objectifs du package `collections` :

* `Map`
* `Reduce`
* `Filter`
* `Find`
* `Some` / `Every`
* (optionnel) `SortBy`

---

### 📁 Arborescence proposée

```
myproject/
├── go.mod
├── main.go
└── collections/
    ├── collections.go
    └── collections_test.go
```

---

### ✅ `collections/collections.go`

```go
package collections

// Map transforme chaque élément d'une slice.
func Map[T any, R any](input []T, f func(T) R) []R {
	result := make([]R, len(input))
	for i, v := range input {
		result[i] = f(v)
	}
	return result
}

// Reduce accumule les éléments d'une slice.
func Reduce[T any, R any](input []T, initial R, f func(R, T) R) R {
	acc := initial
	for _, v := range input {
		acc = f(acc, v)
	}
	return acc
}

// Filter garde les éléments qui satisfont un prédicat.
func Filter[T any](input []T, predicate func(T) bool) []T {
	result := []T{}
	for _, v := range input {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}

// Find retourne le premier élément qui satisfait un prédicat.
func Find[T any](input []T, predicate func(T) bool) (T, bool) {
	var zero T
	for _, v := range input {
		if predicate(v) {
			return v, true
		}
	}
	return zero, false
}

// Some retourne vrai si au moins un élément satisfait le prédicat.
func Some[T any](input []T, predicate func(T) bool) bool {
	for _, v := range input {
		if predicate(v) {
			return true
		}
	}
	return false
}

// Every retourne vrai si tous les éléments satisfont le prédicat.
func Every[T any](input []T, predicate func(T) bool) bool {
	for _, v := range input {
		if !predicate(v) {
			return false
		}
	}
	return true
}
```

---

### ✅ `collections/collections_test.go` (extraits de tests)

```go
package collections

import (
	"testing"
)

func TestFilter(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	expected := []int{2, 4}

	result := Filter(input, func(x int) bool {
		return x%2 == 0
	})

	if len(result) != len(expected) {
		t.Errorf("Filter length mismatch: got %v, want %v", result, expected)
	}
	for i := range result {
		if result[i] != expected[i] {
			t.Errorf("Filter failed at index %d: got %d, want %d", i, result[i], expected[i])
		}
	}
}

func TestFind(t *testing.T) {
	input := []int{1, 3, 5, 8}
	val, found := Find(input, func(x int) bool {
		return x%2 == 0
	})

	if !found || val != 8 {
		t.Errorf("Find failed: got %d, found=%v", val, found)
	}
}

func TestSomeEvery(t *testing.T) {
	values := []int{2, 4, 6}

	if !Some(values, func(x int) bool { return x == 4 }) {
		t.Error("Some failed: expected true")
	}

	if !Every(values, func(x int) bool { return x%2 == 0 }) {
		t.Error("Every failed: expected true")
	}
}
```

---

### 🧪 Tester

```bash
go test ./collections
```

---
