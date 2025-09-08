package main

import "fmt"

func run() {
	// Une chaîne contenant des caractères Unicode
	// str := "Hello 🌍 éàç"
	str := "radar"

	// Conversion en slice de runes
	runes := []rune(str)

	// Affichage des runes et de leurs codes
	for i, r := range runes {
		// fmt.Printf("Rune %d: %c (U+%04X)\n", i, r, r)
		fmt.Printf("Rune %d: %c %d)\n", i, r, int(r))
	}
}
