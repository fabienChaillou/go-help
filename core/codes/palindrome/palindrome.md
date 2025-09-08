Voici un script simple en Go qui vérifie si une chaîne de caractères est un palindrome (c’est-à-dire qu’elle se lit de la même manière à l’endroit et à l’envers) :

```go
package main

import (
	"fmt"
	"strings"
	"unicode"
)

// isPalindrome vérifie si une chaîne est un palindrome (sans tenir compte des espaces, de la casse et de la ponctuation).
func isPalindrome(s string) bool {
	// Nettoyer la chaîne : enlever la ponctuation et mettre en minuscules
	var cleaned []rune
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			cleaned = append(cleaned, unicode.ToLower(r))
		}
	}

	// Vérification du palindrome
	n := len(cleaned)
	for i := 0; i < n/2; i++ {
		if cleaned[i] != cleaned[n-1-i] {
			return false
		}
	}
	return true
}

func main() {
	examples := []string{
		"Radar",
		"Engage le jeu que je le gagne",
		"Able was I, ere I saw Elba",
		"Bonjour",
	}

	for _, text := range examples {
		if isPalindrome(text) {
			fmt.Printf("✔ \"%s\" est un palindrome\n", text)
		} else {
			fmt.Printf("✘ \"%s\" n'est pas un palindrome\n", text)
		}
	}
}
```

### Points clés :

* Ce script ignore la casse, les espaces et la ponctuation.
* Il utilise `unicode` pour être compatible avec des caractères accentués ou non-latins.

Souhaites-tu une version qui traite aussi les caractères accentués (avec normalisation Unicode) ?

Voici comment tu peux ajouter un **test unitaire** en Go pour la fonction `isPalindrome`, en utilisant le package standard `testing` :

### 🧪 Fichier de test : `palindrome_test.go`

```go
package main

import "testing"

func TestIsPalindrome(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"Radar", true},
		{"Engage le jeu que je le gagne", true},
		{"Able was I, ere I saw Elba", true},
		{"Bonjour", false},
		{"Ésope reste ici et se repose", true},
		{"12321", true},
		{"12345", false},
	}

	for _, tt := range tests {
		result := isPalindrome(tt.input)
		if result != tt.expected {
			t.Errorf("isPalindrome(%q) = %v; want %v", tt.input, result, tt.expected)
		}
	}
}
```

### ✅ Instructions pour exécuter le test

Place ce fichier à côté de ton fichier `main.go`, puis exécute la commande suivante :

```bash
go test
```

Tu veux aussi que je te montre une version avec **normalisation Unicode** pour mieux gérer les accents ?

Parfait ! Pour gérer correctement les **accents et caractères spéciaux Unicode**, il faut **normaliser** la chaîne (ex. : transformer `é` en `e`). On utilise pour cela le package `golang.org/x/text/unicode/norm`.

---

### ✅ Étapes :

1. **Installer le module nécessaire** (si ce n’est pas déjà fait) :

   ```bash
   go get golang.org/x/text/unicode/norm
   ```

---

### ✨ Mise à jour du code avec normalisation Unicode :

#### `main.go`

```go
package main

import (
	"fmt"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
	"golang.org/x/text/transform"
)

// removeAccents retire les accents des caractères Unicode
func removeAccents(s string) string {
	t := transform.Chain(norm.NFD, transform.RemoveFunc(isNonSpacingMark), norm.NFC)
	result, _, _ := transform.String(t, s)
	return result
}

// isNonSpacingMark permet d'identifier les accents
func isNonSpacingMark(r rune) bool {
	return unicode.Is(unicode.Mn, r) // Mn = Mark, nonspacing
}

// isPalindrome vérifie si une chaîne est un palindrome
func isPalindrome(s string) bool {
	// Nettoyage : minuscules, sans accents, lettres/chiffres uniquement
	s = removeAccents(s)
	var cleaned []rune
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			cleaned = append(cleaned, unicode.ToLower(r))
		}
	}
	n := len(cleaned)
	for i := 0; i < n/2; i++ {
		if cleaned[i] != cleaned[n-1-i] {
			return false
		}
	}
	return true
}

func main() {
	examples := []string{
		"Ésope reste ici et se repose",
		"Engage le jeu que je le gagne",
		"Radar",
		"Able was I, ere I saw Elba",
		"Bonjour",
		"été",
	}

	for _, text := range examples {
		if isPalindrome(text) {
			fmt.Printf("✔ \"%s\" est un palindrome\n", text)
		} else {
			fmt.Printf("✘ \"%s\" n'est pas un palindrome\n", text)
		}
	}
}
```

---

### 🧪 Test mis à jour (dans `palindrome_test.go`) :

```go
package main

import "testing"

func TestIsPalindrome(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"Radar", true},
		{"Engage le jeu que je le gagne", true},
		{"Able was I, ere I saw Elba", true},
		{"Bonjour", false},
		{"Ésope reste ici et se repose", true},
		{"été", true},
		{"12321", true},
		{"12345", false},
	}

	for _, tt := range tests {
		if result := isPalindrome(tt.input); result != tt.expected {
			t.Errorf("isPalindrome(%q) = %v; want %v", tt.input, result, tt.expected)
		}
	}
}
```
