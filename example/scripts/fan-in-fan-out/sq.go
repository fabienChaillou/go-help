package main

// The second stage, sq, receives integers from a channel and returns a channel that emits the square of each received integer.
// After the inbound channel is closed and this stage has sent all the values downstream, it closes the outbound channel:
func sq(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}
