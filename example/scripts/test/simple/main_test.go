// main_test.go
package main

import (
	"testing"
)

func TestSimpleAdd(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"positive numbers", 2, 3, 5},
		{"negative numbers", -2, -3, -5},
		{"mixed numbers", -2, 5, 3},
		{"zeros", 0, 0, 0},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := SimpleAdd(test.a, test.b)
			if result != test.expected {
				t.Errorf("SimpleAdd(%d, %d) = %d; expected %d",
					test.a, test.b, result, test.expected)
			}
		})
	}
}

func TestGreet(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"normal greeting", "John", "Hello, John!"},
		{"empty name", "", "Hello, !"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := Greet(test.input)
			if result != test.expected {
				t.Errorf("Greet(%q) = %q; expected %q",
					test.input, result, test.expected)
			}
		})
	}
}

func TestWordCount(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"simple sentence", "Go is awesome", 3},
		{"empty string", "", 0},
		{"multiple spaces", "Go  is  awesome", 3},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := WordCount(test.input)
			if result != test.expected {
				t.Errorf("WordCount(%q) = %d; expected %d",
					test.input, result, test.expected)
			}
		})
	}
}
