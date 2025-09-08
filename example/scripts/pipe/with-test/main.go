package pipe

import (
	"fmt"
	"io"
	"os"
)

// PipeExample demonstrates a simple use of os.Pipe()
// It writes a message to the pipe and reads it back
func PipeExample(message string) (string, error) {
	// Create a new pipe
	reader, writer, err := os.Pipe()
	if err != nil {
		return "", fmt.Errorf("failed to create pipe: %w", err)
	}

	// Don't forget to close the pipe
	defer reader.Close()
	defer writer.Close()

	// Write the message to the pipe
	_, err = writer.Write([]byte(message))
	if err != nil {
		return "", fmt.Errorf("failed to write to pipe: %w", err)
	}

	// Close the write end to signal EOF to the reader
	writer.Close()

	// Read the message from the pipe
	buf := make([]byte, 1024)
	n, err := reader.Read(buf)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("failed to read from pipe: %w", err)
	}

	return string(buf[:n]), nil
}

// PipeCommand runs a function that writes to the pipe and
// another function that reads from the pipe concurrently
func PipeCommand(input string) (string, error) {
	reader, writer, err := os.Pipe()
	if err != nil {
		return "", fmt.Errorf("failed to create pipe: %w", err)
	}

	// Use a channel to collect result from the goroutine
	resultChan := make(chan string, 1)
	errChan := make(chan error, 1)

	// Start a goroutine to read from the pipe
	go func() {
		defer reader.Close()

		// Read everything from the pipe
		data, err := io.ReadAll(reader)
		if err != nil {
			errChan <- fmt.Errorf("read error: %w", err)
			return
		}

		// Process the data
		resultChan <- fmt.Sprintf("Processed: %s", string(data))
	}()

	// Write data to the pipe in the main goroutine
	_, err = writer.Write([]byte(input))
	if err != nil {
		writer.Close()
		return "", fmt.Errorf("write error: %w", err)
	}

	// Close the write end to signal EOF to the reader
	writer.Close()

	// Wait for the result or an error
	select {
	case result := <-resultChan:
		return result, nil
	case err := <-errChan:
		return "", err
	}
}
