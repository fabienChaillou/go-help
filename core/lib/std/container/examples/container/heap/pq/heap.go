package pq

import "container/heap"

type Item struct {
	Value    string
	Priority int
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Priority < pq[j].Priority // min-heap
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x any) {
	*pq = append(*pq, x.(*Item))
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

func NewPriorityQueue(items []*Item) *PriorityQueue {
	pq := PriorityQueue(items)
	heap.Init(&pq)
	return &pq
}

func (pq *PriorityQueue) Insert(item *Item) {
	heap.Push(pq, item)
}

func (pq *PriorityQueue) Remove() *Item {
	return heap.Pop(pq).(*Item)
}
