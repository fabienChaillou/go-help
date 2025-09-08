package main

import (
	"fmt"
	"os"
)

func Hello() int {
	fmt.Println("hello world!")
	return 0
}

func main() {
	exitCode := Hello()
	os.Exit(exitCode)
}
