package main

import (
	"bytes"
	"fmt"
	"strings"
)

func simple() {
	// Example 1: Using bytes.Buffer for building strings
	fmt.Println("Example 1: Using bytes.Buffer for string building")
	var buffer bytes.Buffer

	// Write strings to the buffer
	buffer.WriteString("Hello")
	buffer.WriteString(", ")
	buffer.WriteString("World!")

	// Get the resulting string
	result := buffer.String()
	fmt.Println("Result:", result)
	fmt.Println("Buffer length:", buffer.Len())

	// Reset the buffer
	buffer.Reset()
	fmt.Println("After reset, length:", buffer.Len())

	// Example 2: Using bytes.Buffer with byte operations
	fmt.Println("\nExample 2: Working with bytes")
	buffer.Write([]byte("Go"))
	buffer.WriteByte(' ')
	buffer.Write([]byte("Programming"))
	fmt.Println("Result:", buffer.String())

	// Example 3: Using strings.Builder (more efficient for string concatenation)
	fmt.Println("\nExample 3: Using strings.Builder")
	var builder strings.Builder
	builder.WriteString("Strings")
	builder.WriteString(" ")
	builder.WriteString("Builder")
	fmt.Println("Result:", builder.String())

	// Example 4: Reading from buffer
	fmt.Println("\nExample 4: Reading from buffer")
	buffer.Reset()
	buffer.WriteString("Hello, Reader!")

	// Read first 5 bytes
	b := make([]byte, 5)
	n, err := buffer.Read(b)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("Read %d bytes: %s\n", n, b)
		fmt.Println("Remaining in buffer:", buffer.String())
	}

	// Example 5: Creating buffer from existing data
	fmt.Println("\nExample 5: Buffer from existing data")
	data := []byte("Pre-populated buffer")
	newBuffer := bytes.NewBuffer(data)
	fmt.Println("New buffer:", newBuffer.String())
}
