// main.go
package main

import (
	"fmt"
	"strings"
)

// SimpleAdd adds two integers and returns the result
func SimpleAdd(a, b int) int {
	return a + b
}

// Greet returns a greeting message for the given name
func Greet(name string) string {
	return fmt.Sprintf("Hello, %s!", name)
}

// WordCount returns the number of words in a string
func WordCount(s string) int {
	return len(strings.Fields(s))
}

func main() {
	fmt.Println(Greet("Gopher"))
	fmt.Println("2 + 3 =", SimpleAdd(2, 3))
	fmt.Println("Words in 'Go is awesome':", WordCount("Go is awesome"))
}
