package factorielle

import "testing"

func TestFactorielle(t *testing.T) {
	tests := []struct {
		input    int
		expected int
	}{
		{0, 1},
		{1, 1},
		{2, 2},
		{3, 6},
		{4, 24},
		{5, 120},
	}

	for _, tt := range tests {
		result := Factorielle(tt.input)
		if result != tt.expected {
			t.Errorf("Factorielle(%d) = %d; want %d", tt.input, result, tt.expected)
		}
	}
}

func TestFactorielleIterative(t *testing.T) {
	tests := []struct {
		input    int
		expected int
	}{
		{0, 1},
		{1, 1},
		{2, 2},
		{3, 6},
		{4, 24},
		{5, 120},
	}

	for _, tt := range tests {
		result := FactorielleIterative(tt.input)
		if result != tt.expected {
			t.Errorf("FactorielleIterative(%d) = %d; want %d", tt.input, result, tt.expected)
		}
	}
}

func BenchmarkFactorielleRecursive(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Factorielle(10)
	}
}

func BenchmarkFactorielleIterative(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FactorielleIterative(10)
	}
}
