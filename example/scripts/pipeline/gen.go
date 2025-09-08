package main

// The first stage, gen, is a function that converts a list of integers to a channel that emits the integers in the list.
// The gen function starts a goroutine that sends the integers on the channel and closes the channel when all the values have been sent:
func gen(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}
