package main

import "sync"

// The merge function converts a list of channels to a single channel by starting a goroutine
// for each inbound channel that copies the values to the sole outbound channel.
// Once all the output goroutines have been started, merge starts one more goroutine to close the outbound channel
// after all sends on that channel are done.

// Sends on a closed channel panic, so it’s important to ensure all sends are done before calling close.
// The sync.WaitGroup type provides a simple way to arrange this synchronization:
func merge(cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan int) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
