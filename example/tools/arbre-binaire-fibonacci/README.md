Voici deux scripts en Go : un pour le **parcours d’un arbre binaire** et un autre pour la **suite de Fibonacci (récursif et itératif)**. Tu peux choisir selon ce qui t'intéresse le plus.

---

## 🌳 1. **Parcours récursif d’un arbre binaire**

```go
package main

import "fmt"

// Définition du noeud de l'arbre
type Node struct {
	Value int
	Left  *Node
	Right *Node
}

// Parcours in-order (gauche - racine - droite)
func InOrderTraversal(n *Node) {
	if n == nil {
		return
	}
	InOrderTraversal(n.Left)
	fmt.Print(n.Value, " ")
	InOrderTraversal(n.Right)
}

func main() {
	// Construction de l'arbre :
	//       4
	//     /   \
	//    2     6
	//   / \   / \
	//  1   3 5   7

	root := &Node{
		Value: 4,
		Left: &Node{
			Value: 2,
			Left:  &Node{Value: 1},
			Right: &Node{Value: 3},
		},
		Right: &Node{
			Value: 6,
			Left:  &Node{Value: 5},
			Right: &Node{Value: 7},
		},
	}

	fmt.Print("In-Order Traversal: ")
	InOrderTraversal(root)
	fmt.Println()
}
```

---

## 🧮 2. **Suite de Fibonacci (récursif et itératif)**

```go
package main

import "fmt"

// Version récursive
func FibonacciRecursive(n int) int {
	if n <= 1 {
		return n
	}
	return FibonacciRecursive(n-1) + FibonacciRecursive(n-2)
}

// Version itérative (optimisée)
func FibonacciIterative(n int) int {
	if n <= 1 {
		return n
	}
	a, b := 0, 1
	for i := 2; i <= n; i++ {
		a, b = b, a+b
	}
	return b
}

func main() {
	n := 10
	fmt.Printf("Fibonacci (récursif) de %d = %d\n", n, FibonacciRecursive(n))
	fmt.Printf("Fibonacci (itératif) de %d = %d\n", n, FibonacciIterative(n))
}
```

---

Souhaites-tu que je t’ajoute un test unitaire ou benchmark pour l’un de ces deux scripts ?

Très bien ! Voici les **tests unitaires et benchmarks** pour les deux scripts :

---

## ✅ 1. Tests pour la **suite de Fibonacci**

### 📁 `fibonacci.go`

```go
package fibonacci

func FibonacciRecursive(n int) int {
	if n <= 1 {
		return n
	}
	return FibonacciRecursive(n-1) + FibonacciRecursive(n-2)
}

func FibonacciIterative(n int) int {
	if n <= 1 {
		return n
	}
	a, b := 0, 1
	for i := 2; i <= n; i++ {
		a, b = b, a+b
	}
	return b
}
```

### 🧪 `fibonacci_test.go`

```go
package fibonacci

import "testing"

func TestFibonacci(t *testing.T) {
	cases := []struct {
		n        int
		expected int
	}{
		{0, 0},
		{1, 1},
		{2, 1},
		{3, 2},
		{4, 3},
		{5, 5},
		{10, 55},
	}

	for _, c := range cases {
		if res := FibonacciRecursive(c.n); res != c.expected {
			t.Errorf("FibonacciRecursive(%d) = %d; want %d", c.n, res, c.expected)
		}
		if res := FibonacciIterative(c.n); res != c.expected {
			t.Errorf("FibonacciIterative(%d) = %d; want %d", c.n, res, c.expected)
		}
	}
}

func BenchmarkFibonacciRecursive(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FibonacciRecursive(20) // Attention : coûteux pour n > 30
	}
}

func BenchmarkFibonacciIterative(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FibonacciIterative(20)
	}
}
```

---

## ✅ 2. Test pour le **parcours d’arbre**

### 📁 `tree.go`

```go
package tree

type Node struct {
	Value int
	Left  *Node
	Right *Node
}

func InOrderTraversal(n *Node, visit func(int)) {
	if n == nil {
		return
	}
	InOrderTraversal(n.Left, visit)
	visit(n.Value)
	InOrderTraversal(n.Right, visit)
}
```

### 🧪 `tree_test.go`

```go
package tree

import (
	"reflect"
	"testing"
)

func TestInOrderTraversal(t *testing.T) {
	// Arbre :
	//       4
	//     /   \
	//    2     6
	//   / \   / \
	//  1   3 5   7
	root := &Node{
		Value: 4,
		Left: &Node{
			Value: 2,
			Left:  &Node{Value: 1},
			Right: &Node{Value: 3},
		},
		Right: &Node{
			Value: 6,
			Left:  &Node{Value: 5},
			Right: &Node{Value: 7},
		},
	}

	var result []int
	InOrderTraversal(root, func(v int) {
		result = append(result, v)
	})

	expected := []int{1, 2, 3, 4, 5, 6, 7}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("InOrderTraversal() = %v; want %v", result, expected)
	}
}
```

---

Tu peux lancer tous les tests avec :

```bash
go test ./...
```

Et les benchmarks avec :

```bash
go test -bench=.
```

Souhaites-tu aussi des **représentations graphiques** du benchmark ou une **visualisation de l’arbre** ?

Parfait ! Voici deux choses pour enrichir ton projet :

---

## 📊 1. **Représentation graphique des benchmarks**

Go ne génère pas de graphiques directement, mais tu peux utiliser `benchstat` pour comparer les performances entre deux versions, ou exporter les données et les tracer avec un outil externe (comme Excel ou Python).

### Exemple : Générer les benchmarks dans un fichier

```bash
go test -bench=. -benchmem > bench.txt
```

Puis, tu peux utiliser [`benchstat`](https://pkg.go.dev/golang.org/x/perf/cmd/benchstat) :

```bash
benchstat old.txt new.txt
```

> Ou bien, tu peux extraire les temps et les tracer avec Python / Excel pour obtenir une courbe de temps d’exécution en fonction de `n`.

Souhaites-tu que je t’écrive un petit script Python pour tracer un graphe comparatif ?

---

## 🌳 2. **Visualisation de l’arbre binaire**

On peut afficher l’arbre en texte (simple) ou générer un fichier `.dot` compatible avec [Graphviz](https://graphviz.org/) pour une vraie visualisation graphique.

### ✅ Génération d’un fichier `.dot`

Voici un code Go pour exporter ton arbre en format Graphviz :

```go
package tree

import (
	"fmt"
	"os"
)

func ExportDOT(root *Node, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintln(f, "digraph G {")
	exportDOTNode(f, root)
	fmt.Fprintln(f, "}")
	return nil
}

func exportDOTNode(f *os.File, n *Node) {
	if n == nil {
		return
	}
	if n.Left != nil {
		fmt.Fprintf(f, "  %d -> %d;\n", n.Value, n.Left.Value)
		exportDOTNode(f, n.Left)
	}
	if n.Right != nil {
		fmt.Fprintf(f, "  %d -> %d;\n", n.Value, n.Right.Value)
		exportDOTNode(f, n.Right)
	}
}
```

### Utilisation :

```go
func main() {
	root := &Node{ /* comme avant */ }
	ExportDOT(root, "tree.dot")
}
```

### Générer l'image avec Graphviz :

```bash
dot -Tpng tree.dot -o tree.png
```

> Cela te donne une belle image `tree.png` de ton arbre !

---

Voici un **Makefile** simple et pratique pour ton projet Go, incluant les commandes les plus utiles : build, test, benchmark, nettoyage, génération de visuel d’arbre, etc.

---

### 📄 `Makefile`

```makefile
# Variables
BINARY_NAME=main
PKG=./...
DOT_FILE=tree.dot
PNG_FILE=tree.png

# Compilation
build:
	go build -o $(BINARY_NAME) .

# Exécution
run: build
	./$(BINARY_NAME)

# Tests
test:
	go test $(PKG)

# Tests avec couverture
cover:
	go test -cover $(PKG)

# Benchmarks
bench:
	go test -bench=. -benchmem $(PKG)

# Générer et convertir le fichier .dot (nécessite Graphviz)
treepng:
	go run .  # Assure-toi que ExportDOT() est dans main
	dot -Tpng $(DOT_FILE) -o $(PNG_FILE)
	@echo "Image générée : $(PNG_FILE)"

# Nettoyage
clean:
	go clean
	rm -f $(BINARY_NAME) $(DOT_FILE) $(PNG_FILE)
```

---

### 🔧 Utilisation

* `make build` : compile le projet.
* `make run` : compile et exécute.
* `make test` : lance les tests.
* `make cover` : lance les tests avec couverture.
* `make bench` : lance les benchmarks.
* `make treepng` : génère un visuel `.png` de ton arbre.
* `make clean` : nettoie les fichiers générés.

---

Souhaites-tu que je t’ajoute une **commande pour générer les stats de benchmark sous forme de graphe (via Python ou autre)** directement depuis le Makefile ?

Parfait ! Voici comment **ajouter une commande dans le Makefile** pour générer un **graphe de performance (Fibonacci récursif vs itératif)** avec un script Python.

---

## 🧮 Étape 1 : Créer le script Python (`plot_benchmark.py`)

Crée ce fichier dans ton projet :

```python
# plot_benchmark.py
import matplotlib.pyplot as plt
import time

def fibonacci_recursive(n):
    if n <= 1:
        return n
    return fibonacci_recursive(n - 1) + fibonacci_recursive(n - 2)

def fibonacci_iterative(n):
    if n <= 1:
        return n
    a, b = 0, 1
    for _ in range(2, n + 1):
        a, b = b, a + b
    return b

max_n = 30
ns = list(range(1, max_n + 1))
times_recursive = []
times_iterative = []

for n in ns:
    # Récursif (plus lent)
    start = time.time()
    if n <= 30:  # au-delà, c’est très lent
        fibonacci_recursive(n)
    times_recursive.append(time.time() - start)

    # Itératif
    start = time.time()
    fibonacci_iterative(n)
    times_iterative.append(time.time() - start)

plt.plot(ns, times_recursive, label='Récursif')
plt.plot(ns, times_iterative, label='Itératif')
plt.xlabel('n (position dans la suite de Fibonacci)')
plt.ylabel('Temps (secondes)')
plt.title('Comparaison performances Fibonacci')
plt.legend()
plt.grid(True)
plt.tight_layout()
plt.savefig("fibonacci_bench.png")
print("✅ Graphe généré : fibonacci_bench.png")
```

---

## 🧪 Étape 2 : Ajouter au **Makefile**

Ajoute cette commande à ton Makefile :

```makefile
# Génère un graphe de benchmark avec Python
fibograph:
	python3 plot_benchmark.py
```

---

## ✅ Étape 3 : Lancer le graphe

Dans le terminal, exécute :

```bash
make fibograph
```

Tu obtiendras une image `fibonacci_bench.png` avec la comparaison visuelle.

---
