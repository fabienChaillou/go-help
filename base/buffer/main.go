package main

import "fmt"

func increment(x, y int, ch chan int) int {
	ch <- 123
	// close(ch)
	return x + y
}

func main() {
	ch := make(chan int, 2) // channel bufferisé

	// ch <- 1
	r1 := increment(1, 2, ch)
	// ch <- 2
	r2 := increment(3, 2, ch)
	// ch <- 3 // ceci bloquerait car le buffer est plein

	fmt.Printf("r1: i %d, ch %d\n", r1, <-ch)
	fmt.Printf("r2: i %d, ch %d\n", r2, <-ch)
	close(ch)
	// fmt.Println(<-ch) // lit 1
	// fmt.Println(<-ch) // lit 2
}
