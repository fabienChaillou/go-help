package pq

import (
	"testing"
)

func TestPriorityQueue(t *testing.T) {
	items := []*Item{
		{Value: "task1", Priority: 3},
		{Value: "task2", Priority: 1},
		{Value: "task3", Priority: 2},
	}

	pq := NewPriorityQueue(items)
	pq.Insert(&Item{Value: "task0", Priority: 0})

	expectedOrder := []string{"task0", "task2", "task3", "task1"}

	for _, expected := range expectedOrder {
		item := pq.Remove()
		if item.Value != expected {
			t.Errorf("Expected %s, got %s", expected, item.Value)
		}
	}
}
