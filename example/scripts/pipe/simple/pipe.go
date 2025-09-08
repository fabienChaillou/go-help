package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	// Create a pipe
	reader, writer, err := os.Pipe()
	if err != nil {
		fmt.Println("Error creating pipe:", err)
		return
	}

	// Write data to the pipe
	go func() {
		defer writer.Close()
		message := []byte("Hello from pipe!")
		_, err := writer.Write(message)
		if err != nil {
			fmt.Println("Error writing to pipe:", err)
		}
	}()

	// Read data from the pipe
	buffer := make([]byte, 100)
	n, err := reader.Read(buffer)
	if err != nil && err != io.EOF {
		fmt.Println("Error reading from pipe:", err)
		return
	}

	// Print the received message
	fmt.Printf("Received: %s\n", buffer[:n])

	// Close the reader
	reader.Close()
}
