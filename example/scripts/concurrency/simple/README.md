# Golang conrurrency simple example

I'd be happy to provide a Go concurrency example! Go's concurrency model is based on goroutines and channels, which are lightweight threads managed by the Go runtime and a mechanism for communication between goroutines.

Here's a practical example that demonstrates these concepts:

## Key Concepts

1. **Goroutines**: Lightweight threads created with the go keyword. In the example, each data fetch operation runs concurrently.
2. **Channels**: Used for communication between goroutines. The results channel allows the fetching goroutines to send their results back to the main goroutine.
3. **WaitGroup**: Used to wait for a collection of goroutines to finish. This ensures we close the channel only after all goroutines have completed.

## What the Code Does
The program simulates fetching data from multiple services concurrently:

1. It creates 5 goroutines, each "fetching" data with different delays
2. Each goroutine sends its result to a channel when complete
3. The main goroutine reads from the channel and displays results as they arrive
4. The program shows the total time taken, which is much less than if the operations were sequential

When you run this program, you'll see that results arrive in the order they complete (not the order they were started), and the total time is only as long as the slowest operation rather than the sum of all operations.