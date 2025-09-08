package pipe

import (
	"testing"
)

func TestPipeExample(t *testing.T) {
	message := "Hello, pipe world!"
	result, err := PipeExample(message)

	if err != nil {
		t.Fatalf("PipeExample failed: %v", err)
	}

	if result != message {
		t.Errorf("Expected message %q, got %q", message, result)
	}
}

func TestPipeCommand(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{
			input:    "test data",
			expected: "Processed: test data",
		},
		{
			input:    "",
			expected: "Processed: ",
		},
		{
			input:    "Go pipes are cool",
			expected: "Processed: Go pipes are cool",
		},
	}

	for _, tc := range testCases {
		result, err := PipeCommand(tc.input)

		if err != nil {
			t.Fatalf("PipeCommand failed with input %q: %v", tc.input, err)
		}

		if result != tc.expected {
			t.Errorf("Expected result %q, got %q", tc.expected, result)
		}
	}
}
