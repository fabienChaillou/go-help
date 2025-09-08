package decompose

import "fmt"

// PrintStrings affiche une liste de chaînes avec fmt.Println
func PrintStrings(words ...string) string {
	output := fmt.Sprintln(words)
	return output
}

// Sum additionne une liste d'entiers
func Sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

// MergeSlices fusionne deux slices d'entiers
func MergeSlices(a []int, b []int) []int {
	return append(a, b...) // décomposition de b
}
