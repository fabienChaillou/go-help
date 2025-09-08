package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	// rand.Seed(time.Now().UnixNano()) // Initialisation avec le timestamp actuel
	seed := time.Now().UnixNano()
	source := rand.NewSource(seed)
	// source.Seed(seed)
	r := rand.New(source)

	n := r.Intn(100) // génère un entier entre 0 et 99
	fmt.Println(n)
}
