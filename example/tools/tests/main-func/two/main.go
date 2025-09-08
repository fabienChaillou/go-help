// main.go
package main

import (
	"fmt"
)

// Sum returns the sum of two integers
func Sum(a, b int) int {
	return a + b
}

// Multiply returns the product of two integers
func Multiply(a, b int) int {
	return a * b
}

func main() {
	fmt.Println("Sum of 5 and 3:", Sum(5, 3))
	fmt.Println("Product of 5 and 3:", Multiply(5, 3))
}
