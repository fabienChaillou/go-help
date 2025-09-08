En Go, si tu as une séquence de `byte` représentant une chaîne encodée en UTF-8 (ce qui est souvent le cas), tu peux simplement la convertir en `string` pour l'afficher correctement, y compris les caractères Unicode.

### Exemple :

```go
package main

import "fmt"

func main() {
    // Une chaîne encodée en UTF-8 sous forme de slice de bytes
    unicodeBytes := []byte{0xe2, 0x9c, 0x94} // ✅ (check mark)

    // Conversion en string
    unicodeStr := string(unicodeBytes)

    // Impression
    fmt.Println("Unicode character:", unicodeStr)
}
```

### Explication :

* `[]byte{0xe2, 0x9c, 0x94}` correspond à `U+2705` (✅) en UTF-8.
* En Go, `string([]byte)` fait la conversion correctement si les bytes sont en UTF-8.

Tu peux aussi tester avec d'autres caractères Unicode comme les emojis ou les lettres accentuées.

Souhaites-tu un exemple où le byte slice est mal encodé ou partiellement corrompu pour voir comment gérer les erreurs ?

Ah, parfait — si tu veux **imprimer une chaîne à partir d'une slice de `rune`** (et non pas de `[]byte`), c'est encore plus simple, car une `rune` représente déjà un **caractère Unicode** (valeur `int32` correspondant à un point de code Unicode).

### Exemple :

```go
package main

import "fmt"

func main() {
    // Une slice de runes représentant des caractères Unicode
    runes := []rune{0x1F600, 0x2764, 0x20AC} // 😀 ❤ €

    // Conversion en string
    str := string(runes)

    // Impression
    fmt.Println("Chaîne depuis des runes :", str)
}
```

### Détails :

* `0x1F600` = 😀 (Emoji visage souriant)
* `0x2764` = ❤ (Cœur)
* `0x20AC` = € (Euro)

### Résultat :

```
Chaîne depuis des runes : 😀❤€
```

➡️ En Go, faire `string([]rune)` construit correctement une chaîne UTF-8.

Souhaites-tu aussi convertir une `string` en `[]rune` pour faire l’inverse ?

Parfait ! Pour **convertir une `string` en `[]rune`** en Go, tu peux simplement faire un cast direct. Cela permet de découper la chaîne en **caractères Unicode** (et non en bytes), ce qui est très utile pour manipuler les caractères multibyte comme les emojis ou les lettres accentuées.

### Exemple :

```go
package main

import "fmt"

func main() {
    // Une chaîne contenant des caractères Unicode
    str := "Hello 🌍 éàç"

    // Conversion en slice de runes
    runes := []rune(str)

    // Affichage des runes et de leurs codes
    for i, r := range runes {
        fmt.Printf("Rune %d: %c (U+%04X)\n", i, r, r)
    }
}
```

### Résultat attendu :

```
Rune 0: H (U+0048)
Rune 1: e (U+0065)
Rune 2: l (U+006C)
Rune 3: l (U+006C)
Rune 4: o (U+006F)
Rune 5:   (U+0020)
Rune 6: 🌍 (U+1F30D)
Rune 7:   (U+0020)
Rune 8: é (U+00E9)
Rune 9: à (U+00E0)
Rune 10: ç (U+00E7)
```

### Pourquoi utiliser `[]rune` ?

* `string` est une suite de **bytes UTF-8**.
* `[]rune` est une suite de **caractères Unicode** (points de code).
* Utile pour parcourir des chaînes contenant des caractères multibytes (emojis, idéogrammes, etc.).

## Add unit test

Bien sûr ! Voici un exemple complet avec une **fonction** qui convertit une `string` en `[]rune`, et un **test unitaire** associé.

---

### 🧩 Fichier : `runes.go`

```go
package runesutil

// ToRunes convertit une chaîne UTF-8 en slice de runes (caractères Unicode).
func ToRunes(s string) []rune {
    return []rune(s)
}
```

---

### ✅ Fichier de test : `runes_test.go`

```go
package runesutil

import (
    "reflect"
    "testing"
)

func TestToRunes(t *testing.T) {
    input := "Go ✅ éà"
    expected := []rune{'G', 'o', ' ', '✅', ' ', 'é', 'à'}

    result := ToRunes(input)

    if !reflect.DeepEqual(result, expected) {
        t.Errorf("ToRunes(%q) = %v, want %v", input, result, expected)
    }
}
```

---

### 🧪 Pour lancer le test :

```sh
go test
```

### Pourquoi `reflect.DeepEqual` ?

* Car comparer deux slices rune par rune manuellement serait plus verbeux.
* `DeepEqual` compare le contenu des slices correctement.

Souhaites-tu que je t’aide à écrire des tests inverses (reconstruire la string depuis les runes) ou à utiliser `testing.T.Run` pour plusieurs cas ?

Parfait, ajoutons maintenant :

1. **Un test inverse** : vérifier que la conversion `[]rune → string` est correcte.
2. **Des sous-tests (`t.Run`)** pour tester plusieurs cas dans une seule fonction.

---

### 🧩 Mise à jour du fichier : `runes.go`

```go
package runesutil

// ToRunes convertit une chaîne UTF-8 en slice de runes.
func ToRunes(s string) []rune {
    return []rune(s)
}

// FromRunes reconstruit une chaîne UTF-8 à partir d'une slice de runes.
func FromRunes(r []rune) string {
    return string(r)
}
```

---

### ✅ Mise à jour de `runes_test.go` avec sous-tests

```go
package runesutil

import (
    "reflect"
    "testing"
)

func TestRuneConversions(t *testing.T) {
    tests := []struct {
        name     string
        inputStr string
        expected []rune
    }{
        {"ASCII", "Go", []rune{'G', 'o'}},
        {"Emoji", "✅", []rune{'✅'}},
        {"Accent", "éàç", []rune{'é', 'à', 'ç'}},
        {"Mixed", "Go 🌍 é", []rune{'G', 'o', ' ', '🌍', ' ', 'é'}},
    }

    for _, tt := range tests {
        t.Run(tt.name+"_ToRunes", func(t *testing.T) {
            runes := ToRunes(tt.inputStr)
            if !reflect.DeepEqual(runes, tt.expected) {
                t.Errorf("ToRunes(%q) = %v, want %v", tt.inputStr, runes, tt.expected)
            }
        })

        t.Run(tt.name+"_FromRunes", func(t *testing.T) {
            str := FromRunes(tt.expected)
            if str != tt.inputStr {
                t.Errorf("FromRunes(%v) = %q, want %q", tt.expected, str, tt.inputStr)
            }
        })
    }
}
```

---

### ✅ Résultat attendu avec `go test -v`

```
=== RUN   TestRuneConversions/ASCII_ToRunes
=== RUN   TestRuneConversions/ASCII_FromRunes
=== RUN   TestRuneConversions/Emoji_ToRunes
=== RUN   TestRuneConversions/Emoji_FromRunes
...
--- PASS: TestRuneConversions (0.00s)
```

Souhaites-tu intégrer des cas limites comme une chaîne vide ou des caractères non imprimables ?

Parfait, ajoutons maintenant **des cas limites** à nos tests, y compris :

* Chaîne vide (`""`)
* Caractères non imprimables (`'\n'`, `'\t'`, etc.)
* Rune maximum (`\U0010FFFF`) — le plus grand point de code Unicode valide
* Rune invalide (simulée via `0xFFFD`, caractère de remplacement)

---

### ✅ Mise à jour complète des tests dans `runes_test.go`

```go
package runesutil

import (
    "reflect"
    "testing"
)

func TestRuneConversions(t *testing.T) {
    tests := []struct {
        name     string
        inputStr string
        expected []rune
    }{
        {"ASCII", "Go", []rune{'G', 'o'}},
        {"Emoji", "✅", []rune{'✅'}},
        {"Accent", "éàç", []rune{'é', 'à', 'ç'}},
        {"Mixed", "Go 🌍 é", []rune{'G', 'o', ' ', '🌍', ' ', 'é'}},
        {"Empty", "", []rune{}},
        {"NonPrintable", "\n\t", []rune{'\n', '\t'}},
        {"MaxRune", string([]rune{0x10FFFF}), []rune{0x10FFFF}},
        {"ReplacementChar", string([]rune{0xFFFD}), []rune{0xFFFD}}, // U+FFFD = �
    }

    for _, tt := range tests {
        t.Run(tt.name+"_ToRunes", func(t *testing.T) {
            runes := ToRunes(tt.inputStr)
            if !reflect.DeepEqual(runes, tt.expected) {
                t.Errorf("ToRunes(%q) = %v, want %v", tt.inputStr, runes, tt.expected)
            }
        })

        t.Run(tt.name+"_FromRunes", func(t *testing.T) {
            str := FromRunes(tt.expected)
            if str != tt.inputStr {
                t.Errorf("FromRunes(%v) = %q, want %q", tt.expected, str, tt.inputStr)
            }
        })
    }
}
```

---

### Points importants :

* `0x10FFFF` est la rune valide la plus élevée (Unicode limite).
* `0xFFFD` est utilisée par Go pour représenter un caractère UTF-8 invalide ou mal formé.
* Le test `"NonPrintable"` assure que les caractères invisibles sont bien gérés.

Souhaites-tu également tester des chaînes mal encodées en UTF-8 depuis des `[]byte` (pour les passer dans `string(bytes)` et voir comment Go les traite) ?

Parfait, ajoutons maintenant un test spécifique pour les **chaînes mal encodées en UTF-8**, converties depuis un `[]byte`.

En Go, quand tu convertis des bytes invalides UTF-8 en `string`, les séquences invalides sont automatiquement remplacées par `�` (`U+FFFD`, le **caractère de remplacement Unicode**). Ce comportement est conforme à la spécification Unicode.

---

### ✅ Ajout de test : `TestInvalidUTF8Bytes`

Ajoute ceci dans ton fichier `runes_test.go` :

```go
func TestInvalidUTF8Bytes(t *testing.T) {
    // Une séquence de bytes invalide (UTF-8 tronqué volontairement)
    invalid := []byte{0xff, 0xfe, 0xf0, 0x28, 0x8c, 0xbc} // invalide UTF-8

    // Conversion vers string
    str := string(invalid)

    // Attendu : chaque séquence invalide devient U+FFFD
    expectedRunes := []rune{'�', '�', '�', '(', '�', '�'}

    // Conversion de la string en runes
    actualRunes := []rune(str)

    if !reflect.DeepEqual(actualRunes, expectedRunes) {
        t.Errorf("Invalid UTF-8 bytes parsed as runes: got %v, want %v", actualRunes, expectedRunes)
    }
}
```

---

### 🔍 Explication :

* `string([]byte)` fait la conversion, même si les bytes sont invalides.
* Go n’échoue pas, mais remplace chaque séquence non valide par `�`.
* On vérifie que cette tolérance est bien gérée par le parseur UTF-8 intégré de Go.

---

### Pour aller plus loin :

* Tu pourrais aussi valider le comportement inverse : essayer de construire une chaîne contenant `�` et t'assurer qu'elle se comporte comme prévu.
* Ou explorer comment les bibliothèques comme `unicode/utf8` peuvent détecter les séquences invalides explicitement.

Souhaites-tu un exemple de détection manuelle d'encodage UTF-8 invalide avec le package `unicode/utf8` ?

Parfait, allons plus loin avec la **détection manuelle de séquences UTF-8 invalides** en Go, en utilisant le package standard [`unicode/utf8`](https://pkg.go.dev/unicode/utf8).

---

## 🎯 Objectif

Parcourir un tableau de `[]byte` et détecter où se trouvent les **séquences UTF-8 invalides**, sans les convertir silencieusement en `�`.

---

### ✅ Exemple de code avec détection manuelle

```go
package main

import (
    "fmt"
    "unicode/utf8"
)

func main() {
    // Contient des séquences valides et invalides
    data := []byte{'H', 'e', 'l', 'l', 'o', 0xff, 0xfe, '!', 0xe2, 0x28, 0xa1} // mixte valide/invalide

    i := 0
    for i < len(data) {
        r, size := utf8.DecodeRune(data[i:])
        if r == utf8.RuneError && size == 1 {
            fmt.Printf("Invalid UTF-8 byte at index %d: 0x%X\n", i, data[i])
        } else {
            fmt.Printf("Valid rune at index %d: %c (U+%04X)\n", i, r, r)
        }
        i += size
    }
}
```

---

### 🔍 Explication :

* `utf8.DecodeRune` lit un rune UTF-8 valide à partir d’un `[]byte`.
* Si elle retourne `RuneError` **et** `size == 1`, alors c'est une **séquence invalide**.
* `size` peut être >1 même pour `RuneError` si la séquence semble bonne mais ne forme pas une rune valide (cas avancés).

---

### 🧪 Exemple de sortie :

```
Valid rune at index 0: H (U+0048)
Valid rune at index 1: e (U+0065)
Valid rune at index 2: l (U+006C)
Valid rune at index 3: l (U+006C)
Valid rune at index 4: o (U+006F)
Invalid UTF-8 byte at index 5: 0xFF
Invalid UTF-8 byte at index 6: 0xFE
Valid rune at index 7: ! (U+0021)
Invalid UTF-8 byte at index 8: 0xE2
```

---

Souhaites-tu encapsuler ce code dans une fonction de validation que tu pourrais tester unitairement ?

Parfait ! Voici comment encapsuler cette logique dans une **fonction de validation UTF-8 personnalisée** que tu pourras tester facilement.

---

## ✅ Fonction `FindInvalidUTF8Bytes`

```go
package runesutil

import "unicode/utf8"

// FindInvalidUTF8Bytes renvoie les index des bytes invalides dans un tableau UTF-8.
func FindInvalidUTF8Bytes(data []byte) []int {
    var invalidIndexes []int
    i := 0
    for i < len(data) {
        r, size := utf8.DecodeRune(data[i:])
        if r == utf8.RuneError && size == 1 {
            invalidIndexes = append(invalidIndexes, i)
            i++ // avancer d’un byte car on ne connaît pas la taille exacte
        } else {
            i += size
        }
    }
    return invalidIndexes
}
```

---

## ✅ Test unitaire : `runesutil/utf8check_test.go`

```go
package runesutil

import (
    "reflect"
    "testing"
)

func TestFindInvalidUTF8Bytes(t *testing.T) {
    input := []byte{'H', 0xff, 'e', 0xfe, 0xe2, 0x28, '!', 0xf0, 0x28, 0x8c, 0xbc}
    expected := []int{1, 3, 4, 7} // 0xff, 0xfe, 0xe2 (incomplet), 0xf0 (suite invalide)

    result := FindInvalidUTF8Bytes(input)

    if !reflect.DeepEqual(result, expected) {
        t.Errorf("FindInvalidUTF8Bytes returned %v, expected %v", result, expected)
    }
}
```

---

### 🔎 Explication :

* Le test simule des séquences volontairement invalides.
* Les indices retournés correspondent à **l’emplacement exact** des bytes non valides.
* `utf8.DecodeRune` permet une lecture progressive et fiable du buffer.

---

Souhaites-tu aussi une version qui retourne les **runes valides extraites** malgré les erreurs ?
