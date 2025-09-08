package main

import (
	"container/list"
	"fmt"
	"testing"
)

// Fonction principale
func main() {
	l := CreateList()
	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
	}
}

// Fonction pour créer la liste
func CreateList() *list.List {
	l := list.New()
	l.PushBack("A")
	l.PushFront("B")
	l.PushBack("C")
	return l
}

// Fonction utilitaire pour convertir une liste en slice
func ListToSlice(l *list.List) []string {
	var result []string
	for e := l.Front(); e != nil; e = e.Next() {
		result = append(result, e.Value.(string))
	}
	return result
}

// Test unitaire
func TestCreateList(t *testing.T) {
	l := CreateList()
	got := ListToSlice(l)
	want := []string{"B", "A", "C"}

	if len(got) != len(want) {
		t.Fatalf("expected list of length %d, got %d", len(want), len(got))
	}

	for i := range got {
		if got[i] != want[i] {
			t.Errorf("expected element %d to be %q, got %q", i, want[i], got[i])
		}
	}
}
