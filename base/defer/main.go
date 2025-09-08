package main

import "fmt"

func main() {
	x := 10
	defer fmt.Println("x =", x)
	x = 20
	fmt.Println("x =", x)
}
