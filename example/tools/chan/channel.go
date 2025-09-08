package main

import "fmt"

func main() {
	message := make(chan string)

	go func() { message <- "foo" }()

	msg := <-message
	fmt.Println(msg)
}
