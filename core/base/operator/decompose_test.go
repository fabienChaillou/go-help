package decompose

import (
	"testing"
)

func TestPrintStrings(t *testing.T) {
	input := []string{"Go", "est", "génial"}
	expected := "Go est génial\n"

	// Décomposition du slice en arguments
	result := PrintStrings(input...)

	if result != expected {
		t.Errorf("PrintStrings() = %q, want %q", result, expected)
	}
}

func TestSum(t *testing.T) {
	input := []int{1, 2, 3, 4}
	expected := 10

	// Décomposition du slice
	result := Sum(input...)

	if result != expected {
		t.Errorf("Sum() = %d, want %d", result, expected)
	}
}

func TestMergeSlices(t *testing.T) {
	a := []int{1, 2}
	b := []int{3, 4}

	expected := []int{1, 2, 3, 4}
	result := MergeSlices(a, b)

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("MergeSlices()[%d] = %d, want %d", i, v, expected[i])
		}
	}
}
