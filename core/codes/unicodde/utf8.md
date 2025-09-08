En Go, pour convertir une chaîne de caractères (`string`) en un tableau d’octets (`[]byte`), tu peux simplement faire :

```go
str := "Bonjour"
bytes := []byte(str)
```

Cela transforme chaque caractère de la chaîne en son équivalent byte (UTF-8).

### Exemple complet :

```go
package main

import (
	"fmt"
)

func main() {
	str := "Bonjour"
	bytes := []byte(str)
	fmt.Println(bytes)         // Affiche : [66 111 110 106 111 117 114]
	fmt.Println(string(bytes)) // Affiche : Bonjour (reconversion possible)
}
```

Souhaites-tu convertir en bytes pour l’écriture dans un fichier, le hashing, ou autre chose ?

Parfait ! Voici quelques cas courants où on convertit une chaîne en `[]byte` en Go, avec des exemples :

---

### 1. **Écriture dans un fichier**

```go
import (
	"os"
)

func main() {
	str := "Bonjour le monde"
	err := os.WriteFile("exemple.txt", []byte(str), 0644)
	if err != nil {
		panic(err)
	}
}
```

---

### 2. **Calcul d’un hash (ex: SHA256)**

```go
import (
	"crypto/sha256"
	"fmt"
)

func main() {
	str := "motdepasse"
	hash := sha256.Sum256([]byte(str))
	fmt.Printf("Hash SHA256 : %x\n", hash)
}
```

---

### 3. **Envoi sur le réseau (ex: via une requête HTTP)**

```go
import (
	"bytes"
	"net/http"
)

func main() {
	data := "Message à envoyer"
	resp, err := http.Post("http://example.com", "text/plain", bytes.NewBuffer([]byte(data)))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}
```

Bien sûr ! En Go, la fonction `utf8.DecodeRune()` sert à décoder **une rune (caractère Unicode)** à partir d'une séquence d’octets encodés en **UTF-8**.

---

### 📌 Prototype :

```go
func DecodeRune(p []byte) (r rune, size int)
```

* **`p`** : un tableau d’octets (`[]byte`) encodé en UTF-8.
* **Retourne** :

  * `r` : la **rune** décodée (un caractère Unicode).
  * `size` : le **nombre d’octets** utilisés pour ce caractère dans `p`.

---

### ✅ Exemple simple :

```go
package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	input := []byte("école") // 'é' est un caractère multioctet en UTF-8
	r, size := utf8.DecodeRune(input)
	fmt.Printf("Rune: %c, Unicode: %U, Size: %d\n", r, r, size)
}
```

**Sortie** :

```
Rune: é, Unicode: U+00E9, Size: 2
```

---

### 🔎 Notes importantes :

* Si `p` est vide, la rune retournée est `utf8.RuneError` et la taille est 0.
* Si `p` ne commence pas par une séquence UTF-8 valide, la rune sera aussi `utf8.RuneError` mais la taille sera 1 (indiquant une erreur de décodage sur 1 byte).

---

### 📘 Cas d’usage typique : Parcourir une chaîne rune par rune

```go
str := "hé🦊"
data := []byte(str)
for len(data) > 0 {
    r, size := utf8.DecodeRune(data)
    fmt.Printf("Rune: %c, Code point: %U, Bytes used: %d\n", r, r, size)
    data = data[size:] // avance au caractère suivant
}
```

---
