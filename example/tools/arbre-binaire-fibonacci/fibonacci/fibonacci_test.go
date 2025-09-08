package fibonacci

import "testing"

func TestFibonacci(t *testing.T) {
	cases := []struct {
		n        int
		expected int
	}{
		{0, 0},
		{1, 1},
		{2, 1},
		{3, 2},
		{4, 3},
		{5, 5},
		{10, 55},
	}

	for _, c := range cases {
		if res := FibonacciRecursive(c.n); res != c.expected {
			t.Errorf("FibonacciRecursive(%d) = %d; want %d", c.n, res, c.expected)
		}
		if res := FibonacciIterative(c.n); res != c.expected {
			t.Errorf("FibonacciIterative(%d) = %d; want %d", c.n, res, c.expected)
		}
	}
}

func BenchmarkFibonacciRecursive(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FibonacciRecursive(20) // Attention : coûteux pour n > 30
	}
}

func BenchmarkFibonacciIterative(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FibonacciIterative(20)
	}
}
