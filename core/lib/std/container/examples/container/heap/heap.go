package main

import (
	"fmt"
	"priorityqueue/pq"
)

func main() {
	queue := pq.NewPriorityQueue([]*pq.Item{
		{Value: "clean dishes", Priority: 2},
		{Value: "write report", Priority: 1},
	})

	queue.Insert(&pq.Item{Value: "call mom", Priority: 0})

	for queue.Len() > 0 {
		task := queue.Remove()
		fmt.Printf("Doing task: %s\n", task.Value)
	}
}

// type IntHeap []int

// func (h IntHeap) Len() int           { return len(h) }
// func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] } // min-heap
// func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

// func (h *IntHeap) Push(x any) {
// 	*h = append(*h, x.(int))
// }

// func (h *IntHeap) Pop() any {
// 	old := *h
// 	n := len(old)
// 	x := old[n-1]
// 	*h = old[0 : n-1]
// 	return x
// }

// func main() {
// 	h := &IntHeap{3, 1, 4}
// 	heap.Init(h)
// 	heap.Push(h, 2)
// 	fmt.Println(heap.Pop(h)) // 1 (le plus petit)
// }
