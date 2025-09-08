// main.go
package main

import (
	"fmt"
	"os"
)

func main() {
	code := realMain()
	os.Exit(code)
}

// realMain does the actual work and returns an exit code
// This pattern makes the main function testable
func realMain() int {
	fmt.Println("Hello, World!")
	return 0
}
