package main

import (
	"fmt"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	fmt.Println("Start check guess game!")

	secret := "pomme"
	guess := "moppe"

	result, correctCount := CheckGuess(secret, guess)
	fmt.Println("Résultats :", result)
	fmt.Println("Lettres bien placées :", correctCount)
	return nil
}
