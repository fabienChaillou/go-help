Les **pointeurs en Go** (Golang) sont un concept fondamental pour gérer efficacement la **mémoire**, **les performances** et **le partage de données** entre fonctions.

---

## 🧠 Qu'est-ce qu’un pointeur ?

Un **pointeur** est une **variable qui contient l’adresse mémoire** d’une autre variable.

### 🧾 Syntaxe de base :

```go
var a int = 42
var p *int = &a // p est un pointeur vers a
```

* `a` est une variable entière.
* `&a` est l’adresse mémoire de `a`.
* `*int` signifie "pointeur vers un int".

---

## 🔁 Accès et modification via pointeur

```go
fmt.Println(*p) // Affiche la valeur pointée par p → 42
*p = 100        // Modifie la valeur de a via le pointeur
fmt.Println(a)  // Affiche 100
```

---

## ✅ Pourquoi utiliser des pointeurs ?

| Avantage                  | Explication                                                |
| ------------------------- | ---------------------------------------------------------- |
| ✅ **Modifier une valeur** | Permet de modifier une variable passée à une fonction      |
| ✅ **Éviter la copie**     | Plus performant avec des structures ou tableaux volumineux |
| ✅ **Partage de données**  | Deux variables peuvent accéder à la même mémoire           |

---

## 📦 Exemple complet

```go
package main

import "fmt"

func increment(x *int) {
    *x = *x + 1
}

func main() {
    a := 10
    increment(&a)     // Passe l’adresse de a
    fmt.Println(a)    // Affiche 11
}
```

* `increment` reçoit un `*int` → pointeur vers un entier.
* `*x = *x + 1` → lit/modifie la valeur pointée.

---

## 🧱 Pointeurs et `struct`

```go
type Person struct {
    Name string
    Age  int
}

func birthday(p *Person) {
    p.Age++
}

func main() {
    p := Person{"Alice", 30}
    birthday(&p)
    fmt.Println(p.Age) // 31
}
```

---

## ❌ Ce que Go ne permet pas avec les pointeurs (contrairement à C/C++) :

* 🚫 Pas d’arithmétique de pointeur (`p++`, `p+1`, etc.)
* 🚫 Pas de pointeur NULL pur (on utilise `nil`)
* 🚫 Pas de malloc/free manuel (Go a un **ramasse-miettes**)

---

Souhaites-tu que je t’explique la différence entre passer une valeur **par valeur** vs **par pointeur**, ou comment utiliser les pointeurs avec des slices ou des maps ?

---------------------------------------------------------------------------------------------------------------------------

Parfait ! Voyons maintenant la **différence entre passer par valeur et par pointeur** en Go, et pourquoi c’est important.

---

## 📦 1. **Passage par valeur** (comportement par défaut en Go)

En Go, quand tu passes une **variable à une fonction**, elle est **copiée**.

```go
func setToZero(x int) {
    x = 0
}

func main() {
    a := 5
    setToZero(a)
    fmt.Println(a) // ➜ 5 → a n’a pas changé
}
```

➡️ Ici, `x` est une **copie de a**. Modifier `x` ne change pas `a`.

---

## 🧷 2. **Passage par pointeur** (modifie la valeur d'origine)

Pour **modifier la valeur originale**, il faut passer **l'adresse (`&`)** et **déférencer (`*`)** dans la fonction.

```go
func setToZero(x *int) {
    *x = 0
}

func main() {
    a := 5
    setToZero(&a)
    fmt.Println(a) // ➜ 0 → a a bien été modifié
}
```

---

## 🧠 Résumé visuel

| Type         | Ce que la fonction reçoit | Effet sur la variable d’origine |
| ------------ | ------------------------- | ------------------------------- |
| Par valeur   | Une copie                 | ❌ Pas modifié                   |
| Par pointeur | Une adresse (`*Type`)     | ✅ Modifié                       |

---

## 🧱 Et avec `struct` ?

```go
type User struct {
    Name string
}

func changeName(u User) {
    u.Name = "Bob"
}

func changeNamePtr(u *User) {
    u.Name = "Bob"
}

func main() {
    u := User{"Alice"}
    changeName(u)
    fmt.Println(u.Name) // ➜ "Alice" (copie)

    changeNamePtr(&u)
    fmt.Println(u.Name) // ➜ "Bob" (pointeur)
}
```

---

## 🎯 Et avec slices, maps, channels ?

Bonne nouvelle : **slices, maps, channels sont déjà des références** en Go !
Tu **n'as pas besoin de pointeur supplémentaire**.

```go
func add(s []int) {
    s[0] = 999
}

func main() {
    nums := []int{1, 2, 3}
    add(nums)
    fmt.Println(nums) // ➜ [999 2 3]
}
```

Mais si tu changes la **référence entière**, tu auras besoin d’un pointeur si tu veux que ça se reflète dehors :

```go
func reset(s *[]int) {
    *s = []int{0, 0, 0}
}

func main() {
    nums := []int{1, 2, 3}
    reset(&nums)
    fmt.Println(nums) // ➜ [0 0 0]
}
```

---

Souhaites-tu un petit tableau de comparaison "Quand utiliser un pointeur" ou un exercice à faire pour t'entraîner ?
