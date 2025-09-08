package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	s := "héllo"
	fmt.Println(len(s))                    // Longueur en bytes (6)
	fmt.Println(utf8.RuneCountInString(s)) // Longueur en runes (5)

	for i, r := range s {
		fmt.Printf("%d: %c (%U), %d\n", i, r, r, int(r))
	}

	// la fonction utf8.DecodeRune() sert à décoder une rune (caractère Unicode) à partir d'une séquence d’octets encodés en UTF-8.
	str := "hé🦊"
	data := []byte(str)
	for len(data) > 0 {
		r, size := utf8.DecodeRune(data)
		fmt.Printf("Rune: %c, Code point: %U, Code numeric: %d, Bytes used: %d\n", r, int(r), r, size)
		data = data[size:] // avance au caractère suivant
	}
	letter := []byte("é")
	letterE := []byte("e")

	fmt.Printf("Rune: %c, Code point: %U, Bytes used: %d\n", letter, letter, len(letter))
	fmt.Printf("Rune: %c, Code point: %U, Bytes used: %d\n", letterE, letterE, len(letterE))
	// fmt.Printf(len("e"))
}
