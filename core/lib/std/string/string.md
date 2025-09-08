En Go, un **`string`** représente une **séquence immuable d'octets**, souvent utilisée pour contenir du texte. Voici une explication claire et complète de ce type :

---

### 🔹 1. **Définition**

```go
var s string = "Bonjour"
```

* Une `string` en Go est **immuable** (elle ne peut pas être modifiée après sa création).
* Elle contient une séquence de **bytes** (octets), pas directement des caractères Unicode.

---

### 🔹 2. **Immuabilité**

Tu ne peux pas faire ceci :

```go
s[0] = 'b' // Erreur de compilation
```

Car les `string` sont **en lecture seule**. Si tu veux modifier une chaîne, il faut en créer une nouvelle.

---

### 🔹 3. **Longueur**

```go
len("école") // retourne 6, car "é" prend 2 octets
```

⚠️ `len()` retourne le **nombre d’octets**, pas de caractères.

---

### 🔹 4. **Parcourir une string**

Pour itérer correctement sur les **runes** (caractères Unicode), utilise `range` :

```go
for i, r := range "école" {
    fmt.Printf("Index: %d Rune: %c\n", i, r)
}
```

Ici, `r` est une `rune`, qui représente un caractère Unicode (int32).

---

### 🔹 5. **Conversions**

* `string([]byte)` : conversion d’un tableau d’octets en string.
* `[]byte("bonjour")` : conversion d’un string en tableau d’octets.
* `[]rune("école")` : pour obtenir les caractères Unicode séparément.

---

### 🔹 6. **Fonctions utiles**

Le package [`strings`](https://pkg.go.dev/strings) fournit de nombreuses fonctions :

```go
import "strings"

strings.ToUpper("go")       // "GO"
strings.Contains("go", "o") // true
strings.Split("a,b,c", ",") // ["a" "b" "c"]
strings.TrimSpace(" go ")   // "go"
```

---

### Résumé

| Élément     | Description                        |
| ----------- | ---------------------------------- |
| `string`    | Séquence d’octets immuable         |
| `len(s)`    | Taille en octets                   |
| `range`     | Itération rune par rune (Unicode)  |
| `[]byte(s)` | Conversion vers tableau d’octets   |
| `strings`   | Fonctions de manipulation de texte |

---

Voici un exemple simple de fonction en Go qui manipule des chaînes : elle prend une phrase, retire les espaces autour, met tout en majuscules et compte les mots.

---

### 🔧 Exemple : Traitement de texte

```go
package main

import (
    "fmt"
    "strings"
)

// Fonction qui formate une phrase
func NettoyerEtCompter(phrase string) (string, int) {
    // Enlever les espaces avant/après
    nettoyee := strings.TrimSpace(phrase)

    // Convertir en majuscules
    majuscules := strings.ToUpper(nettoyee)

    // Compter les mots
    mots := strings.Fields(majuscules) // découpe par espace
    nbMots := len(mots)

    return majuscules, nbMots
}

func main() {
    phrase := "  Bonjour le monde de Go  "
    nettoyee, nbMots := NettoyerEtCompter(phrase)

    fmt.Println("Texte nettoyé :", nettoyee)
    fmt.Println("Nombre de mots :", nbMots)
}
```

---

### 📌 Sortie :

```
Texte nettoyé : BONJOUR LE MONDE DE GO
Nombre de mots : 5
```

Cette fonction montre comment utiliser plusieurs fonctions du package `strings` pour travailler efficacement avec des chaînes de caractères.

## manupuler des rune
Parfait ! Voici un exemple clair d’utilisation des **`rune`** pour manipuler des caractères Unicode dans une chaîne (comme les lettres accentuées, les emojis, etc.).

---

### 🧪 Exemple : Inverser une chaîne Unicode (avec `rune`)

```go
package main

import (
    "fmt"
)

// Fonction qui inverse les caractères (runes) d'une chaîne
func Inverser(s string) string {
    // Convertir en slice de runes
    runes := []rune(s)

    // Inverser les runes
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]
    }

    // Retourner la chaîne reconstruite
    return string(runes)
}

func main() {
    original := "Éléphant 🐘"
    inversée := Inverser(original)

    fmt.Println("Original :", original)
    fmt.Println("Inversée :", inversée)
}
```

---

### ✅ Résultat :

```
Original : Éléphant 🐘
Inversée : 🐘 tnahpélÉ
```

---

### 📌 Pourquoi utiliser `[]rune` ?

Si tu faisais `[]byte("Éléphant")`, le caractère `É` serait mal découpé, car il occupe **plus d’un octet**. `[]rune` résout ce problème car chaque `rune` représente un **caractère Unicode complet** (`int32`).

Souhaites-tu voir un autre exemple avec des accents ou des caractères spéciaux (comme le filtrage ou le remplacement) ?

Très bien ! Voici un exemple de **filtrage des caractères accentués ou spéciaux** dans une chaîne en Go, en utilisant les `rune`.

---

### 🎯 Objectif :

Créer une fonction qui enlève tous les caractères **non alphanumériques** (lettres et chiffres), en conservant les **lettres accentuées**.

---

### 🧪 Exemple : Nettoyer une chaîne

```go
package main

import (
    "fmt"
    "unicode"
)

// Fonction pour filtrer les caractères alphanumériques et lettres accentuées
func NettoyerTexte(s string) string {
    var resultat []rune

    for _, r := range s {
        if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) {
            resultat = append(resultat, r)
        }
        // Sinon on ignore (ponctuation, symboles, etc.)
    }

    return string(resultat)
}

func main() {
    original := "Salut! Ça va ? Très bien 😊 #GoLang123"
    nettoye := NettoyerTexte(original)

    fmt.Println("Original :", original)
    fmt.Println("Nettoyé  :", nettoye)
}
```

---

### ✅ Résultat :

```
Original : Salut! Ça va ? Très bien 😊 #GoLang123
Nettoyé  : Salut Ça va Très bien GoLang123
```

---

### 📌 Détail :

* `unicode.IsLetter(r)` garde les lettres, y compris accentuées (`é`, `ç`, etc.)
* `unicode.IsDigit(r)` garde les chiffres
* `unicode.IsSpace(r)` garde les espaces
* Les symboles, emojis et ponctuations sont ignorés

Souhaites-tu ajouter une fonction qui **remplace les accents** par leur version sans accent (ex: `é → e`) ?

Parfait ! Voici comment **remplacer les lettres accentuées** par leur version **non accentuée** (ex: `é → e`, `ç → c`, etc.) en Go.

Go ne fournit pas cette conversion directement, mais on peut utiliser le package `golang.org/x/text/unicode/norm` combiné à `unicode.IsMn` pour enlever les marques diacritiques (accents).

---

### 📦 Prérequis

D'abord, installe le module externe (si ce n'est pas déjà fait) :

```bash
go get golang.org/x/text
```

---

### 🧪 Exemple : Supprimer les accents

```go
package main

import (
    "fmt"
    "strings"
    "unicode"

    "golang.org/x/text/unicode/norm"
)

// Supprimer les accents d'une chaîne (normalisation NFD + filtre des marques)
func SupprimerAccents(s string) string {
    var b strings.Builder
    // Décomposer les caractères accentués (NFD)
    for _, r := range norm.NFD.String(s) {
        // Ne pas inclure les marques diacritiques (accents)
        if !unicode.Is(unicode.Mn, r) {
            b.WriteRune(r)
        }
    }
    return b.String()
}

func main() {
    original := "Élève très motivé, ça c’est sûr !"
    sansAccents := SupprimerAccents(original)

    fmt.Println("Original     :", original)
    fmt.Println("Sans accents :", sansAccents)
}
```

---

### ✅ Résultat :

```
Original     : Élève très motivé, ça c’est sûr !
Sans accents : Eleve tres motive, ca c'est sur !
```

---

### 📌 Ce que fait cette fonction :

1. **Normalisation NFD** décompose les lettres accentuées (ex: `é → e + ́`)
2. **Filtrage avec `unicode.IsMn`** enlève les marques diacritiques (les accents)

