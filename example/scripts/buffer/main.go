package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func main() {
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

// BufferOps contains buffer operations we want to test
type BufferOps struct{}

// ConcatStrings concatenates strings using bytes.Buffer
func (bo *BufferOps) ConcatStrings(strs ...string) string {
	var buffer bytes.Buffer
	for _, str := range strs {
		buffer.WriteString(str)
	}
	return buffer.String()
}

// ProcessBytes manipulates bytes in a buffer
func (bo *BufferOps) ProcessBytes(input []byte) string {
	buffer := bytes.NewBuffer(input)

	// Read first 3 bytes if available
	if buffer.Len() >= 3 {
		prefix := make([]byte, 3)
		buffer.Read(prefix)

		// Convert remaining to uppercase (simple processing example)
		remaining := buffer.Bytes()
		for i := 0; i < len(remaining); i++ {
			if remaining[i] >= 'a' && remaining[i] <= 'z' {
				remaining[i] = remaining[i] - 32 // Convert to uppercase
			}
		}

		// Create new result
		result := bytes.NewBuffer(prefix)
		result.Write(remaining)
		return result.String()
	}

	return string(input)
}

// BuildWithFormat builds formatted text using strings.Builder
func (bo *BufferOps) BuildWithFormat(items []string) string {
	var builder strings.Builder

	for i, item := range items {
		if i > 0 {
			builder.WriteString(", ")
		}
		fmt.Fprintf(&builder, "Item %d: %s", i+1, item)
	}

	return builder.String()
}

// TestConcatStrings tests string concatenation with buffer
func TestConcatStrings(t *testing.T) {
	bo := &BufferOps{}

	tests := []struct {
		input    []string
		expected string
	}{
		{[]string{"Hello", ", ", "World"}, "Hello, World"},
		{[]string{"Go", "lang"}, "Golang"},
		{[]string{}, ""},
		{[]string{"Single"}, "Single"},
	}

	for i, test := range tests {
		result := bo.ConcatStrings(test.input...)
		if result != test.expected {
			t.Errorf("Test %d: Expected %q but got %q", i, test.expected, result)
		}
	}
}

// TestProcessBytes tests byte processing using buffer
func TestProcessBytes(t *testing.T) {
	bo := &BufferOps{}

	tests := []struct {
		input    []byte
		expected string
	}{
		{[]byte("golang"), "golANG"},
		{[]byte("UPPER"), "UPPER"},
		{[]byte("ab"), "ab"},         // Too short, returns original
		{[]byte("abc123"), "abc123"}, // Numbers remain unchanged
	}

	for i, test := range tests {
		result := bo.ProcessBytes(test.input)
		if result != test.expected {
			t.Errorf("Test %d: Expected %q but got %q", i, test.expected, result)
		}
	}
}

// TestBuildWithFormat tests string building with formatting
func TestBuildWithFormat(t *testing.T) {
	bo := &BufferOps{}

	tests := []struct {
		input    []string
		expected string
	}{
		{[]string{"apple", "banana"}, "Item 1: apple, Item 2: banana"},
		{[]string{"one"}, "Item 1: one"},
		{[]string{}, ""},
	}

	for i, test := range tests {
		result := bo.BuildWithFormat(test.input)
		if result != test.expected {
			t.Errorf("Test %d: Expected %q but got %q", i, test.expected, result)
		}
	}
}
