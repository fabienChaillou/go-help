En Go (Golang), **Unicode** est bien pris en charge nativement. Voici une explication claire de ce que cela signifie et comment l’utiliser efficacement.

---

### 🔤 Qu’est-ce qu’Unicode ?

**Unicode** est un standard international qui attribue un code unique à chaque caractère, quel que soit le langage ou le système. Par exemple :

* `'A'` → U+0041
* `'é'` → U+00E9
* `'你'` → U+4F60

Unicode couvre des milliers de caractères (latin, chinois, emoji, etc.).

---

### 🧠 Comment Go gère Unicode ?

Go utilise :

* Le type `string`, qui est une **séquence d’octets UTF-8**.
* Le type `rune`, qui représente **un code point Unicode** (un alias de `int32`).

---

### 📦 Exemples de base

#### 1. `string` et `rune`

```go
s := "héllo"
fmt.Println(len(s))        // Longueur en bytes (6)
fmt.Println(utf8.RuneCountInString(s)) // Longueur en runes (5)

for i, r := range s {
    fmt.Printf("%d: %c (%U)\n", i, r, r)
}
```

Sortie :

```
0: h (U+0068)
1: é (U+00E9)
3: l (U+006C)
4: l (U+006C)
5: o (U+006F)
```

> ⚠️ Le `range` sur une `string` lit des runes, pas des bytes.

---

#### 2. Conversion entre `string`, `[]byte` et `[]rune`

```go
s := "你好"
runes := []rune(s)
bytes := []byte(s)

fmt.Println(len(s))       // 6 (UTF-8: 3 bytes par caractère)
fmt.Println(len(runes))   // 2 (2 runes)
fmt.Println(runes)        // [20320 22909]
```

---

### 🧰 Fonctions utiles (dans `unicode` et `unicode/utf8`)

* `utf8.RuneCountInString(s)` : nombre de runes
* `utf8.DecodeRune([]byte)` : décode le premier rune
* `unicode.IsLetter(r)` : est-ce une lettre ?
* `unicode.IsDigit(r)` : est-ce un chiffre ?
* `unicode.ToUpper(r)` : majuscule
* `unicode.ToLower(r)` : minuscule

---

### Exemple complet avec détection de lettres

```go
import (
    "fmt"
    "unicode"
)

func main() {
    s := "Café123"
    for _, r := range s {
        if unicode.IsLetter(r) {
            fmt.Printf("%c is a letter\n", r)
        } else if unicode.IsDigit(r) {
            fmt.Printf("%c is a digit\n", r)
        }
    }
}
```

---
