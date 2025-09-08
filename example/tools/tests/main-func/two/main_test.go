// main_test.go
package main

import (
	"testing"
)

func TestSum(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"Positive numbers", 5, 3, 8},
		{"Zero and positive", 0, 7, 7},
		{"Negative numbers", -2, -3, -5},
		{"Positive and negative", 5, -3, 2},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := Sum(tc.a, tc.b)
			if result != tc.expected {
				t.Errorf("Sum(%d, %d) = %d; expected %d", tc.a, tc.b, result, tc.expected)
			}
		})
	}
}

func TestMultiply(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"Positive numbers", 5, 3, 15},
		{"Zero and positive", 0, 7, 0},
		{"Negative numbers", -2, -3, 6},
		{"Positive and negative", 5, -3, -15},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := Multiply(tc.a, tc.b)
			if result != tc.expected {
				t.Errorf("Multiply(%d, %d) = %d; expected %d", tc.a, tc.b, result, tc.expected)
			}
		})
	}
}

// Example of a benchmark test
func BenchmarkSum(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Sum(5, 3)
	}
}

// Example of a benchmark test
func BenchmarkMultiply(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Multiply(5, 3)
	}
}
