## Simple Golang os.pipe Example

Here's a simple example of using os.Pipe() in Go:

This example creates a pipe using `os.Pipe()`, which returns a pair of connected file descriptors: a reader and a writer. It then:

1. Launches a goroutine to write a message to the pipe and close the writer when done
2. Reads the message from the pipe in the main goroutine
3. Prints the received message

When you run this program, it should output:

`Received: Hello from pipe!`
