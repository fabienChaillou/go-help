package main

import "fmt"

// The main function sets up the pipeline and runs the final stage:
// it receives values from the second stage and prints each one, until the channel is closed:
func main() {
	fmt.Println("Start!")
	// Set up the pipeline.
	c := gen(2, 3)
	out := sq(c)

	// Consume the output.
	fmt.Println(<-out) // 4
	fmt.Println(<-out) // 9

	// Since sq has the same type for its inbound and outbound channels, we can compose it any number of times.
	// We can also rewrite main as a range loop, like the other stages:
	// Set up the pipeline and consume the output.
	for n := range sq(sq(gen(2, 3))) {
		fmt.Println(n) // 16 then 81
	}
}
