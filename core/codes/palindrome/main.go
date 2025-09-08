package main

import (
	"fmt"
	"log"
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
	log.Println("Cleaned: ", string(cleaned))

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
