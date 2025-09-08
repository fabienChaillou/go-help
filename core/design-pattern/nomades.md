Le **pattern monade** (ou *monad pattern*) est un concept issu de la **programmation fonctionnelle** (notamment Haskell), qui peut paraître un peu abstrait au départ, surtout en Go qui est **impératif et non-fonctionnel** par nature. Pourtant, on peut **imiter certains comportements monadiques** en Go pour **enchaîner des opérations**, **gérer les erreurs**, ou **éviter le code répétitif**.

---

## 🧠 Définition intuitive

Une **monade** est un "conteneur d'effet" qui permet de **composer des opérations** tout en **propageant automatiquement un contexte** (erreur, valeur absente, effet, etc.).

Une monade :

* **Wrappe** une valeur (ex: `Just(x)`, `Some(x)`)
* Fournit un moyen de **mapper** une fonction sur cette valeur (`map`, `flatMap` ou `bind`)
* Gère **implicitement le contexte** (ex: erreurs, absences, logs, etc.)

---

## 📦 En Go : simulation de monades

Même si Go n’a **pas de monades natives**, on peut **imiter** des comportements monadiques dans trois cas courants :

---

### 1. **Option/Maybe monade** (valeur présente ou absente)

```go
type Maybe[T any] struct {
    value *T
}

func Some[T any](v T) Maybe[T] {
    return Maybe[T]{value: &v}
}

func None[T any]() Maybe[T] {
    return Maybe[T]{value: nil}
}

func (m Maybe[T]) Map(f func(T) T) Maybe[T] {
    if m.value == nil {
        return None[T]()
    }
    res := f(*m.value)
    return Some(res)
}

func (m Maybe[T]) GetOrElse(defaultVal T) T {
    if m.value == nil {
        return defaultVal
    }
    return *m.value
}
```

### 🔍 Utilisation

```go
func main() {
    some := Some(10)
    none := None[int]()

    double := func(x int) int { return x * 2 }

    fmt.Println(some.Map(double).GetOrElse(0)) // 20
    fmt.Println(none.Map(double).GetOrElse(0)) // 0
}
```

---

### 2. **Either monade** (valeur ou erreur)

```go
type Either[T any] struct {
    value T
    err   error
}

func Right[T any](v T) Either[T] {
    return Either[T]{value: v, err: nil}
}

func Left[T any](err error) Either[T] {
    var zero T
    return Either[T]{value: zero, err: err}
}

func (e Either[T]) Map(f func(T) T) Either[T] {
    if e.err != nil {
        return e
    }
    return Right(f(e.value))
}

func (e Either[T]) GetOrElse(defaultVal T) T {
    if e.err != nil {
        return defaultVal
    }
    return e.value
}
```

---

### 3. **Chaining (façon monadique)**

Go gère souvent les erreurs manuellement, mais tu peux les chaîner de façon plus fluide avec une structure :

```go
type Chain[T any] struct {
    value T
    err   error
}

func Wrap[T any](v T, err error) Chain[T] {
    return Chain[T]{v, err}
}

func (c Chain[T]) Then(f func(T) (T, error)) Chain[T] {
    if c.err != nil {
        return c
    }
    val, err := f(c.value)
    return Wrap(val, err)
}

func (c Chain[T]) Unwrap() (T, error) {
    return c.value, c.err
}
```

### 🔧 Exemple d’usage

```go
func ajouter1(x int) (int, error) {
    return x + 1, nil
}

func diviserPar2(x int) (int, error) {
    if x%2 != 0 {
        return 0, fmt.Errorf("non divisible par 2")
    }
    return x / 2, nil
}

func main() {
    result := Wrap(4, nil).
        Then(ajouter1).    // 5
        Then(diviserPar2). // erreur ici
        Unwrap()

    fmt.Println(result) // 0, erreur
}
```

---

## ✅ Résumé

| Monade simulée   | Go équivalent                 |
| ---------------- | ----------------------------- |
| `Maybe`          | struct avec `*T` ou `nil`     |
| `Either`         | struct avec `val` ou `err`    |
| `Bind / flatMap` | méthode `Then(f)` ou `Map(f)` |
| Composition      | Chaînage d’appels sécurisés   |

---

Souhaites-tu un exemple concret dans un contexte métier (ex: traitement d'une commande, validation utilisateur, etc.) en utilisant cette approche monadique ?

## avec un traitement d'une commande

Super, prenons un **exemple métier** où l'on traite une `Commande` avec plusieurs étapes :

* Validation de la commande (produits valides, stock suffisant)
* Calcul du total avec remises
* Paiement
* Mise à jour du statut

On va utiliser une **approche monadique simulée en Go**, via un enchaînement fluide avec une struct `Chain`.

---

## 📦 Structures de base

```go
package main

import (
    "errors"
    "fmt"
)

type Commande struct {
    ID        int
    Produits  []Produit
    Statut    string
    Total     float64
    Payée     bool
}

type Produit struct {
    Nom   string
    Prix  float64
    Stock int
}
```

---

## 🔄 Simuler une monade de traitement

```go
type Chain[T any] struct {
    value T
    err   error
}

func Wrap[T any](v T, err error) Chain[T] {
    return Chain[T]{v, err}
}

func (c Chain[T]) Then(f func(T) (T, error)) Chain[T] {
    if c.err != nil {
        return c
    }
    val, err := f(c.value)
    return Wrap(val, err)
}

func (c Chain[T]) Unwrap() (T, error) {
    return c.value, c.err
}
```

---

## 🧠 Étapes métier

```go
func ValiderCommande(cmd Commande) (Commande, error) {
    if len(cmd.Produits) == 0 {
        return cmd, errors.New("aucun produit dans la commande")
    }
    for _, p := range cmd.Produits {
        if p.Stock <= 0 {
            return cmd, fmt.Errorf("produit '%s' en rupture", p.Nom)
        }
    }
    return cmd, nil
}

func CalculerTotal(cmd Commande) (Commande, error) {
    var total float64
    for _, p := range cmd.Produits {
        total += p.Prix
    }

    // Exemple de remise si +3 articles
    if len(cmd.Produits) >= 3 {
        total *= 0.9
    }

    cmd.Total = total
    return cmd, nil
}

func PayerCommande(cmd Commande) (Commande, error) {
    if cmd.Total <= 0 {
        return cmd, errors.New("total invalide")
    }

    cmd.Payée = true
    return cmd, nil
}

func FinaliserCommande(cmd Commande) (Commande, error) {
    if !cmd.Payée {
        return cmd, errors.New("paiement non effectué")
    }

    cmd.Statut = "terminée"
    return cmd, nil
}
```

---

## 🚀 Exemple d’exécution

```go
func main() {
    commande := Commande{
        ID: 1,
        Produits: []Produit{
            {"Livre", 20, 10},
            {"Stylo", 5, 5},
            {"Cahier", 7, 3},
        },
    }

    result, err := Wrap(commande, nil).
        Then(ValiderCommande).
        Then(CalculerTotal).
        Then(PayerCommande).
        Then(FinaliserCommande).
        Unwrap()

    if err != nil {
        fmt.Println("Erreur :", err)
    } else {
        fmt.Printf("Commande traitée avec succès : %+v\n", result)
    }
}
```

---

## ✅ Résultat simulé

```
Commande traitée avec succès : {ID:1 Produits:[{Livre 20 10} {Stylo 5 5} {Cahier 7 3}] Statut:terminée Total:28.8 Payée:true}
```

---

## 🧠 Avantages

* Évite le code spaghetti avec `if err != nil` partout
* Enchaînement clair et fluide des étapes métier
* Reproduit une **forme de monade de résultat (Either / Result)**
