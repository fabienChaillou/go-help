En Go, une **closure** (ou fermeture) est une fonction **définie à l'intérieur d'une autre fonction**, qui **capture et peut accéder aux variables** de la fonction englobante, **même après que celle-ci ait terminé son exécution**.

---

### 📌 Définition simple

Une **closure** est une fonction **anonyme** (ou pas) qui **enferme** dans son contexte les variables qu’elle utilise et qui viennent de l'extérieur de son propre corps.

---

### ✅ Exemple simple

```go
package main

import "fmt"

func adder() func(int) int {
    sum := 0
    return func(x int) int {
        sum += x
        return sum
    }
}

func main() {
    f := adder()
    fmt.Println(f(10)) // 10
    fmt.Println(f(5))  // 15
    fmt.Println(f(3))  // 18
}
```

---

### 🔍 Ce qu’il se passe :

* La fonction `adder` retourne **une fonction** qui **accède** à la variable `sum`.
* Même après que `adder()` ait terminé son exécution, la fonction retournée continue de **se souvenir de la variable `sum`**.
* Chaque appel à `f()` modifie **la même instance** de `sum`.

---

### 💡 Pourquoi utiliser des closures ?

* **Éviter les variables globales**.
* **Créer des générateurs ou des compteurs**.
* Passer des fonctions qui gardent un état dans des boucles ou des goroutines.
* Encapsuler une logique avec un **état privé**.

---

### 🧠 Attention

* Chaque appel à `adder()` retourne **une nouvelle closure indépendante**, avec sa propre copie de `sum`.
* Cela permet d’avoir **plusieurs accumulateurs** indépendants :

```go
a := adder()
b := adder()

fmt.Println(a(1)) // 1
fmt.Println(a(2)) // 3
fmt.Println(b(1)) // 1
```

Souhaites-tu un exemple plus pratique, par exemple dans un contexte de traitement de données ?

Parfait ! Voici un exemple de **closure dans un traitement de données** en Go. Imaginons qu'on veuille **filtrer** une liste de nombres selon un critère **paramétrable dynamiquement** (comme un `min` ou `max`) — au lieu de coder chaque filtre séparément, on utilise une closure.

---

### 🎯 Objectif : générer dynamiquement des filtres

```go
package main

import "fmt"

// filter générique utilisant une fonction (closure) comme critère
func filter(data []int, test func(int) bool) []int {
    var result []int
    for _, val := range data {
        if test(val) {
            result = append(result, val)
        }
    }
    return result
}

// fabrique une closure qui filtre les nombres > min
func greaterThan(min int) func(int) bool {
    return func(n int) bool {
        return n > min
    }
}

// fabrique une closure qui filtre les nombres pairs
func isEven() func(int) bool {
    return func(n int) bool {
        return n%2 == 0
    }
}

func main() {
    data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

    fmt.Println("Original :", data)

    gt5 := greaterThan(5)
    even := isEven()

    fmt.Println("Supérieur à 5 :", filter(data, gt5))
    fmt.Println("Pairs :", filter(data, even))
}
```

---

### 🔍 Ce que fait ce code :

* `filter` est une fonction générique qui applique un test (une fonction) à chaque élément.
* `greaterThan` et `isEven` retournent des **closures**, c’est-à-dire des fonctions **personnalisées** qui capturent leurs paramètres (`min` ici).
* Ces closures sont utilisées comme filtres **dynamiques**.

---

### 🧠 Avantages des closures ici :

* Code réutilisable.
* Comportement paramétrable sans recopier le même code.
* On peut créer à la volée des fonctions qui encapsulent une logique métier spécifique.

Souhaites-tu un exemple plus poussé, comme une closure utilisée dans une pipeline de transformation de données ?

Parfait ! Voici un exemple **plus poussé**, avec :

* Une **pipeline** de traitement de données numériques.
* Des **closures** pour encapsuler la logique métier.
* Des **tests unitaires** avec le package `testing`.

---

## 📦 Structure proposée

```
data_pipeline/
├── pipeline.go
├── main.go
└── pipeline_test.go
```

---

## ✅ `pipeline.go` – Traitement avec closures

```go
package datapipeline

// Transformer représente une fonction qui transforme une donnée
type Transformer func(int) int

// Filter représente une fonction qui décide si une donnée passe
type Filter func(int) bool

// ProcessData applique un filtre et une transformation sur une liste de données
func ProcessData(data []int, filter Filter, transformer Transformer) []int {
    var result []int
    for _, val := range data {
        if filter(val) {
            result = append(result, transformer(val))
        }
    }
    return result
}

// Closures :

// GreaterThan retourne une closure filtrant les nombres > min
func GreaterThan(min int) Filter {
    return func(n int) bool {
        return n > min
    }
}

// MultiplyBy retourne une closure de transformation
func MultiplyBy(factor int) Transformer {
    return func(n int) int {
        return n * factor
    }
}
```

---

## 🚀 `main.go` – Exemple d'utilisation

```go
package main

import (
    "fmt"

    "data_pipeline/datapipeline"
)

func main() {
    data := []int{1, 5, 8, 10, 3}

    // Filtrer > 5, puis multiplier par 2
    result := datapipeline.ProcessData(
        data,
        datapipeline.GreaterThan(5),
        datapipeline.MultiplyBy(2),
    )

    fmt.Println("Résultat :", result) // [16, 20]
}
```

---

## 🧪 `pipeline_test.go` – Tests unitaires

```go
package datapipeline

import (
    "reflect"
    "testing"
)

func TestProcessData_GreaterThanAndMultiply(t *testing.T) {
    data := []int{2, 4, 6, 8}

    filter := GreaterThan(5)
    transformer := MultiplyBy(3)

    expected := []int{18, 24}
    result := ProcessData(data, filter, transformer)

    if !reflect.DeepEqual(result, expected) {
        t.Errorf("expected %v, got %v", expected, result)
    }
}

func TestGreaterThan(t *testing.T) {
    f := GreaterThan(10)
    if f(9) != false || f(11) != true {
        t.Errorf("GreaterThan(10) closure not working correctly")
    }
}

func TestMultiplyBy(t *testing.T) {
    f := MultiplyBy(4)
    if f(3) != 12 {
        t.Errorf("MultiplyBy(4)(3) = %d; want 12", f(3))
    }
}
```

---

### ✅ Pour lancer les tests :

Dans le dossier du projet :

```bash
go test ./...
```

---
