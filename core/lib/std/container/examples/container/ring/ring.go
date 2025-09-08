package main

import (
	"container/ring"
	"fmt"
	"testing"
)

// Crée une ring de taille n et y insère les entiers 1 à n
func CreateRing(n int) *ring.Ring {
	r := ring.New(n)
	for i := 1; i <= n; i++ {
		r.Value = i
		r = r.Next()
	}
	return r
}

// Convertit une ring en slice
func RingToSlice(r *ring.Ring) []int {
	var result []int
	r.Do(func(v interface{}) {
		if val, ok := v.(int); ok {
			result = append(result, val)
		}
	})
	return result
}

// Fonction principale
func main() {
	r := CreateRing(5)
	r.Do(func(v interface{}) {
		fmt.Println(v)
	})
}

// Test unitaire
func TestCreateRing(t *testing.T) {
	r := CreateRing(4)
	got := RingToSlice(r)
	want := []int{1, 2, 3, 4}

	if len(got) != len(want) {
		t.Fatalf("expected ring of length %d, got %d", len(want), len(got))
	}

	for i := range got {
		if got[i] != want[i] {
			t.Errorf("at index %d: expected %d, got %d", i, want[i], got[i])
		}
	}
}
