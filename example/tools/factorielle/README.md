La **récursivité** en Go (ou Golang), comme dans d'autres langages, est une technique où une fonction s'appelle elle-même pour résoudre un problème. C'est souvent utilisé pour résoudre des problèmes décomposables en sous-problèmes similaires (ex. : calcul du facteur d'une valeur, parcours d'arbre, etc.).

### 1. Exemple simple : factorielle

```go
package main

import "fmt"

func factorielle(n int) int {
	if n == 0 {
		return 1
	}
	return n * factorielle(n-1)
}

func main() {
	fmt.Println(factorielle(5)) // Affiche 120
}
```

### Explication

* **Base case** (cas d'arrêt) : `if n == 0` → pour éviter une boucle infinie.
* **Recursive call** : `factorielle(n-1)` → la fonction s'appelle elle-même.

Chaque appel attend le résultat de l’appel suivant jusqu’à ce que le cas de base soit atteint. Ensuite, les appels "empilés" se résolvent dans l’ordre inverse.

---

### 2. Exemple visuel d'exécution pour `factorielle(3)` :

```
factorielle(3)
=> 3 * factorielle(2)
          => 2 * factorielle(1)
                   => 1 * factorielle(0)
                            => 1 (cas de base)
                   => 1 * 1 = 1
          => 2 * 1 = 2
=> 3 * 2 = 6
```

---

### Bonnes pratiques :

* Toujours **définir un cas d'arrêt**.
* Attention à **ne pas dépasser la pile d'appels** (`stack overflow`) si la récursion est trop profonde.
* Pour certains problèmes, la **récursivité peut être remplacée par une boucle** pour de meilleures performances.

Bien sûr ! Voici comment tu peux ajouter un **test unitaire** pour la fonction récursive `factorielle` en Go, en utilisant le package standard `testing`.

---

### 🔧 1. Fichier source (`factorielle.go`)

```go
package factorielle

func Factorielle(n int) int {
	if n == 0 {
		return 1
	}
	return n * Factorielle(n-1)
}
```

> Remarque : j'ai mis la fonction dans un package `factorielle` pour respecter la bonne pratique de structuration.

---

### 🧪 2. Fichier de test (`factorielle_test.go`)

```go
package factorielle

import "testing"

func TestFactorielle(t *testing.T) {
	tests := []struct {
		input    int
		expected int
	}{
		{0, 1},
		{1, 1},
		{2, 2},
		{3, 6},
		{4, 24},
		{5, 120},
	}

	for _, tt := range tests {
		result := Factorielle(tt.input)
		if result != tt.expected {
			t.Errorf("Factorielle(%d) = %d; want %d", tt.input, result, tt.expected)
		}
	}
}
```

---

### ▶️ 3. Exécution des tests

Tu peux lancer les tests avec cette commande :

```bash
go test
```

Cela exécutera le fichier `factorielle_test.go` et te dira si tous les cas sont OK.

Souhaites-tu que je t’ajoute un test pour une version itérative aussi, ou un test de performance (benchmark) ?

Parfait ! Voici :

---

### 🔁 Version **itérative** de la fonction `Factorielle` (dans `factorielle.go`)

```go
func FactorielleIterative(n int) int {
	if n == 0 {
		return 1
	}
	result := 1
	for i := 2; i <= n; i++ {
		result *= i
	}
	return result
}
```

---

### ✅ Test pour la version itérative (`factorielle_test.go`)

```go
func TestFactorielleIterative(t *testing.T) {
	tests := []struct {
		input    int
		expected int
	}{
		{0, 1},
		{1, 1},
		{2, 2},
		{3, 6},
		{4, 24},
		{5, 120},
	}

	for _, tt := range tests {
		result := FactorielleIterative(tt.input)
		if result != tt.expected {
			t.Errorf("FactorielleIterative(%d) = %d; want %d", tt.input, result, tt.expected)
		}
	}
}
```

---

### 🚀 Benchmark des deux versions (dans `factorielle_test.go`)

```go
func BenchmarkFactorielleRecursive(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Factorielle(10)
	}
}

func BenchmarkFactorielleIterative(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FactorielleIterative(10)
	}
}
```

---

### 🔍 Lancer les benchmarks

```bash
go test -bench=.
```

Cela affichera les performances des deux versions, typiquement en ns/op (nanosecondes par opération).

Souhaites-tu que je te montre une comparaison entre les deux avec des graphes ou des recommandations sur quand utiliser l’une ou l’autre ?


Souhaites-tu un autre exemple plus concret (comme un parcours d'arbre ou la suite de Fibonacci) ?
