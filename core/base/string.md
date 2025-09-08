Voici un exemple simple d'utilisation de **`string`** en Go (Golang) :

```go
package main

import "fmt"

func main() {
    // Déclaration d'une variable de type string
    var message string = "Bonjour, Go !"

    // Ou bien avec l'inférence de type
    autreMessage := "Salut tout le monde"

    // Affichage des chaînes
    fmt.Println(message)
    fmt.Println(autreMessage)

    // Quelques opérations courantes
    fmt.Println("Longueur du message :", len(message))
    fmt.Println("Sous-chaîne :", message[0:7]) // "Bonjour"
    fmt.Println("Concaténation :", message + " " + autreMessage)
}
```

### Explications :

* `string` est un type de données natif représentant une suite d’octets UTF-8.
* On peut faire de la **concaténation** avec `+`, extraire une **sous-chaîne** avec des indices (`[start:end]`), et obtenir la **longueur** avec `len()`.

Souhaites-tu voir des opérations plus avancées comme la recherche, le remplacement ou la comparaison de chaînes ?

Parfait ! Voici un exemple plus avancé montrant différentes opérations sur les chaînes de caractères en Go, à l’aide du package `strings` :

```go
package main

import (
    "fmt"
    "strings"
)

func main() {
    texte := "Go est un langage puissant et simple."

    // 1. Recherche (contains, has prefix/suffix)
    fmt.Println("Contient 'puissant' ?", strings.Contains(texte, "puissant"))
    fmt.Println("Commence par 'Go' ?", strings.HasPrefix(texte, "Go"))
    fmt.Println("Se termine par 'simple.' ?", strings.HasSuffix(texte, "simple."))

    // 2. Remplacement
    nouveauTexte := strings.Replace(texte, "puissant", "rapide", 1)
    fmt.Println("Texte modifié :", nouveauTexte)

    // 3. Division (Split)
    mots := strings.Split(texte, " ")
    fmt.Println("Mots :", mots)

    // 4. Jointure (Join)
    jointure := strings.Join(mots, "-")
    fmt.Println("Texte avec tirets :", jointure)

    // 5. Mise en majuscules / minuscules
    fmt.Println("Majuscules :", strings.ToUpper(texte))
    fmt.Println("Minuscules :", strings.ToLower(texte))

    // 6. Comparaison (insensible à la casse)
    fmt.Println("Égalité (insensible à la casse) :", strings.EqualFold("Go", "gO"))
}
```

### Résumé des fonctions utilisées :

* `strings.Contains`, `HasPrefix`, `HasSuffix` → pour la recherche.
* `strings.Replace` → pour remplacer un mot ou une portion de texte.
* `strings.Split` / `Join` → pour découper ou regrouper des chaînes.
* `ToUpper`, `ToLower` → pour changer la casse.
* `EqualFold` → pour comparer sans tenir compte des majuscules/minuscules.

Souhaites-tu voir un exemple avec des expressions régulières (`regexp`) ou des conversions entre chaînes et autres types (int, float, etc.) ?

### avec des regex

Très bien ! Voici deux types d'exemples supplémentaires :

---

## 🔹 1. **Expressions régulières avec `regexp`**

```go
package main

import (
    "fmt"
    "regexp"
)

func main() {
    texte := "Mon email est exemple@test.com"

    // Détection d'un email avec une expression régulière
    re := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
    
    email := re.FindString(texte)
    fmt.Println("Email trouvé :", email)

    // Remplacer un email par une version masquée
    masque := re.ReplaceAllString(texte, "[email protégé]")
    fmt.Println("Texte masqué :", masque)
}
```

> `regexp.MustCompile()` compile l'expression régulière et panique en cas d’erreur (utilisé quand on est sûr que l'expression est valide).

---

## 🔹 2. **Conversions entre string et int/float avec `strconv`**

```go
package main

import (
    "fmt"
    "strconv"
)

func main() {
    // Conversion d'un int en string
    age := 30
    ageStr := strconv.Itoa(age)
    fmt.Println("Age en string :", ageStr)

    // Conversion d'un float en string
    pi := 3.14159
    piStr := strconv.FormatFloat(pi, 'f', 2, 64)
    fmt.Println("Pi en string :", piStr)

    // Conversion string -> int
    str := "42"
    n, err := strconv.Atoi(str)
    if err == nil {
        fmt.Println("String en int :", n)
    }

    // Conversion string -> float
    strFloat := "3.14"
    f, err := strconv.ParseFloat(strFloat, 64)
    if err == nil {
        fmt.Println("String en float :", f)
    }
}
```

---
