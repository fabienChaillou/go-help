package main

import "testing"

func TestIsPalindrome(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"Radar", true},
		// {"Engage le jeu que je le gagne", true},
		// {"Able was I, ere I saw Elba", true},
		// {"Bonjour", false},
		// {"Esope reste ici et se repose", true},
		// {"12321", true},
		// {"12345", false},
	}

	for _, tt := range tests {
		result := isPalindrome(tt.input)
		if result != tt.expected {
			t.Errorf("isPalindrome(%q) = %v; want %v", tt.input, result, tt.expected)
		}
	}
}
