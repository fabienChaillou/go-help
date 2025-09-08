En Go, pour « déconstruire » un slice (c'est-à-dire en extraire les éléments ou le répartir dans plusieurs variables), tu peux utiliser **l’indexation**, **le range**, ou **l'opérateur de décomposition `...`** selon le contexte.

Voici plusieurs cas d’usage selon ce que tu veux faire :

---

### 🔹 1. **Extraire manuellement des éléments d’un slice**

```go
mySlice := []int{1, 2, 3}
a, b, c := mySlice[0], mySlice[1], mySlice[2]
fmt.Println(a, b, c) // Affiche: 1 2 3
```

---

### 🔹 2. **Extraire les éléments avec `range`**

```go
for i, v := range mySlice {
    fmt.Printf("Index %d: Valeur %d\n", i, v)
}
```

---

### 🔹 3. **Décomposer un slice en arguments de fonction avec `...`**

Si tu as une fonction qui prend plusieurs arguments, et que tu veux passer un slice :

```go
func addition(a, b, c int) int {
    return a + b + c
}

nums := []int{1, 2, 3}
result := addition(nums[0], nums[1], nums[2]) // Classique
// ou bien :
result := addition(nums...) // Décomposition avec ...
```

⚠️ Le `...` ne fonctionne que si le nombre d’éléments du slice correspond exactement aux paramètres attendus.

---

### 🔹 4. **Diviser un slice (par exemple : tête + reste)**

```go
mySlice := []int{10, 20, 30, 40}
head := mySlice[0]
tail := mySlice[1:]

fmt.Println("Head:", head) // 10
fmt.Println("Tail:", tail) // [20 30 40]
```

---

### example avec une function de la lib std

Bien sûr ! Voici un exemple simple de **décomposition d’un slice** avec une fonction de la **librairie standard** de Go qui accepte des arguments variadiques (c'est-à-dire un nombre variable d’arguments), comme `fmt.Println`.

---

### ✅ Exemple avec `fmt.Println` (fonction variadique)

```go
package main

import (
    "fmt"
)

func main() {
    words := []string{"Go", "est", "super"}

    // Décomposition du slice avec ...
    fmt.Println(words...) // Affiche : Go est super
}
```

👉 Ici, `fmt.Println` attend une suite d’arguments `...interface{}`, et on lui donne un `[]string` en le "décomposant" avec `...`.

---

### ✅ Autre exemple avec `append`

`append` est aussi une fonction variadique : tu peux y passer un slice d’éléments.

```go
package main

import (
    "fmt"
)

func main() {
    base := []int{1, 2}
    extras := []int{3, 4, 5}

    // Décomposition de extras avec ...
    result := append(base, extras...)

    fmt.Println(result) // Affiche : [1 2 3 4 5]
}
```

---

Ces deux cas (`fmt.Println` et `append`) sont parmi les plus fréquents. Souhaites-tu un exemple avec une fonction personnalisée aussi ?


Parfait ! Voici un ensemble d'exemples utilisant **l’opérateur de décomposition `...`** en Go, **avec des tests unitaires**. Cela te donnera une base claire pour utiliser cet opérateur dans différents contextes.

---

## 🔧 Fichier `decompose.go`

```go
package decompose

import "fmt"

// PrintStrings affiche une liste de chaînes avec fmt.Println
func PrintStrings(words ...string) string {
	output := fmt.Sprintln(words...)
	return output
}

// Sum additionne une liste d'entiers
func Sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

// MergeSlices fusionne deux slices d'entiers
func MergeSlices(a []int, b []int) []int {
	return append(a, b...) // décomposition de b
}
```

---

## ✅ Fichier `decompose_test.go`

```go
package decompose

import (
	"strings"
	"testing"
)

func TestPrintStrings(t *testing.T) {
	input := []string{"Go", "est", "génial"}
	expected := "Go est génial\n"

	// Décomposition du slice en arguments
	result := PrintStrings(input...)

	if result != expected {
		t.Errorf("PrintStrings() = %q, want %q", result, expected)
	}
}

func TestSum(t *testing.T) {
	input := []int{1, 2, 3, 4}
	expected := 10

	// Décomposition du slice
	result := Sum(input...)

	if result != expected {
		t.Errorf("Sum() = %d, want %d", result, expected)
	}
}

func TestMergeSlices(t *testing.T) {
	a := []int{1, 2}
	b := []int{3, 4}

	expected := []int{1, 2, 3, 4}
	result := MergeSlices(a, b)

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("MergeSlices()[%d] = %d, want %d", i, v, expected[i])
		}
	}
}
```

---

## ▶️ Exécution des tests

Pour lancer les tests :

```bash
go test
```

---

Souhaites-tu un exemple encore plus avancé, par exemple avec une fonction personnalisée variadique qui utilise un type struct ou interface ?
