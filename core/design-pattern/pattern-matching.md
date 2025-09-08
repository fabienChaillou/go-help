En Go, le **pattern matching** (correspondance de motifs) n’est **pas un concept natif** comme dans d'autres langages tels que Rust, Haskell ou encore Scala. Mais on peut **imiter** ce comportement avec certains outils et techniques du langage.

---

## 🔍 Définition (générique)

Le **pattern matching** est une manière d’exprimer des conditions **structurées** sur des valeurs :

> *"si une valeur a telle forme, alors fais ceci..."*

Exemples typiques dans d'autres langages :

```rust
match x {
    Some(v) => println!("Valeur : {}", v),
    None => println!("Aucune valeur"),
}
```

---

## 🛠 En Go : Approche équivalente

Go ne fournit **pas de `match` natif**, mais tu peux le simuler avec :

### 1. **Type switches** (`switch x := val.(type)`)

Utile pour faire du matching sur le type d’une interface.

```go
func MatchType(x interface{}) {
    switch v := x.(type) {
    case int:
        fmt.Println("int :", v)
    case string:
        fmt.Println("string :", v)
    case bool:
        fmt.Println("bool :", v)
    default:
        fmt.Println("type inconnu")
    }
}
```

### 2. **Switch classique avec valeurs**

Utile pour matcher sur des constantes ou conditions simples.

```go
func MatchValue(x int) {
    switch {
    case x < 0:
        fmt.Println("Négatif")
    case x == 0:
        fmt.Println("Zéro")
    case x > 0 && x <= 10:
        fmt.Println("Petit positif")
    default:
        fmt.Println("Grand nombre")
    }
}
```

---

### 3. **Pattern matching structuré (avec struct)**

Tu peux aussi combiner des `switch`, `if`, et des déstructurations pour simuler un matching structuré.

```go
type Event struct {
    Type string
    Data interface{}
}

func HandleEvent(e Event) {
    switch e.Type {
    case "LOGIN":
        if data, ok := e.Data.(string); ok {
            fmt.Println("User login:", data)
        }
    case "ERROR":
        if err, ok := e.Data.(error); ok {
            fmt.Println("Erreur:", err)
        }
    default:
        fmt.Println("Événement inconnu")
    }
}
```

---

### 📦 Bonus : Bibliothèques tierces

Il existe des **libs de pattern matching** pour Go, comme :

* [`match`](https://github.com/alexpantyukhin/go-pattern-matching) (recherche de motifs structurés)
* [`gopatterns`](https://github.com/phelmkamp/gopatterns) (pour les fonctions inspirées de Scala)

Mais dans la plupart des cas, un `switch` bien conçu suffit.

---

## ✅ Résumé

| Langage        | Mot-clé natif                                     | Go supporte ?   |
| -------------- | ------------------------------------------------- | --------------- |
| Rust / Haskell | `match`                                           | ❌               |
| Go             | `switch`, `type switch`                           | ✅ partiellement |
| Solution Go    | `switch`, `if`, `type assertion`, struct matching | ✅               |

---

Souhaites-tu un exemple de pattern matching simulé sur une structure métier (par ex. `Commande`, `Utilisateur`, etc.) ?

## sur une struct commande

Très bien ! Voici un exemple de **pattern matching simulé en Go** sur une structure `Commande`. L’idée est de traiter différents **types d’états** d’une commande (ex : `EnAttente`, `EnCours`, `Terminée`, `Annulée`) avec des actions différentes, comme tu pourrais le faire avec un `match` dans un autre langage.

---

## 🧱 Modèle de base : `Commande`

```go
package main

import (
    "fmt"
)

// Définition d’un état de commande
type EtatCommande string

const (
    EnAttente  EtatCommande = "en_attente"
    EnCours    EtatCommande = "en_cours"
    Terminee   EtatCommande = "terminee"
    Annulee    EtatCommande = "annulee"
)

type Commande struct {
    ID    int
    Etat  EtatCommande
    Total float64
}
```

---

## 🎯 Fonction simulant le pattern matching

```go
func TraiterCommande(c Commande) {
    switch c.Etat {
    case EnAttente:
        fmt.Printf("Commande #%d est en attente. Préparation en cours...\n", c.ID)
    case EnCours:
        fmt.Printf("Commande #%d est en cours de livraison.\n", c.ID)
    case Terminee:
        fmt.Printf("Commande #%d a été terminée. Montant payé : %.2f€\n", c.ID, c.Total)
    case Annulee:
        fmt.Printf("Commande #%d a été annulée. Aucun montant dû.\n", c.ID)
    default:
        fmt.Printf("Commande #%d dans un état inconnu : %s\n", c.ID, c.Etat)
    }
}
```

---

## 🧪 Exemple d’utilisation

```go
func main() {
    commandes := []Commande{
        {ID: 101, Etat: EnAttente, Total: 45.50},
        {ID: 102, Etat: EnCours, Total: 23.00},
        {ID: 103, Etat: Terminee, Total: 99.99},
        {ID: 104, Etat: Annulee, Total: 0.0},
    }

    for _, c := range commandes {
        TraiterCommande(c)
    }
}
```

---

### ✅ Résultat simulé

```
Commande #101 est en attente. Préparation en cours...
Commande #102 est en cours de livraison.
Commande #103 a été terminée. Montant payé : 99.99€
Commande #104 a été annulée. Aucun montant dû.
```

---

### 🧠 Remarques

* Le `switch` ici agit comme un **pattern matcher sur le champ `Etat`**.
* On pourrait étendre cette logique avec des `struct` différentes si les commandes étaient polymorphiques (ex: interface `CommandeTraitement` avec plusieurs types concrets).
