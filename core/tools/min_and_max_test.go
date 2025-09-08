package min_and_max

import (
	"fmt"
	"testing"
)

func IntMinBasic(a, b int) int {
	return min(a, b)
}

func StringMinBasic(a, b string) string {
	return min(a, b)
}

func IntMaxBasic(a, b int) int {
	return max(a, b)
}

func StringtMaxBasic(a, b string) string {
	return max(a, b)
}

func TestIntMinComplex(t *testing.T) {
	ans := max(-1, 5, 3, 2, 6)
	if ans != 6 {
		t.Errorf("IntMin(2, -2) = %d; want -2", ans)
	}
}

func TestIntMin(t *testing.T) {
	ans := min(2, -2)
	if ans != -2 {
		t.Errorf("IntMin(2, -2) = %d; want -2", ans)
	}
}

// Benchmark the performance to IntMin() func
func BenchmarkIntMin(b *testing.B) {
	for b.Loop() {
		IntMinBasic(1, 2)
	}
}

func MinAndMaxTest(t *testing.T) {
	var tests = []struct {
		a    int
		b    []int
		want int
	}{
		{1, []int{1, 2, 3}, 0},
		// {[]int{1, 2, 3}, 0},
		// {[]int{1, 2, 3}, -2},
		// {[]int{1, 2, 3}, -1},
		// {[]int{1, 2, 3}, -1},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%d,%v", tt.a)
		t.Run(testname, func(t *testing.T) {
			ans := min(tt.a, tt.b...) // operator de decomposition
			if ans != tt.want {
				t.Errorf("got %d, want %d", ans, tt.want)
			}
		})
	}
}
