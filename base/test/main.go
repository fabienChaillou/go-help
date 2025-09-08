package main

import (
	"fmt"
	"log"
)

func Add(a, b int) int {
	log.Println("<---- foo ----->")
	return a + b
}

func main() {

	fmt.Println(Add(2, 3))
}
