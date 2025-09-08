Le package `container` en Go est un **ensemble de structures de données de type conteneur**, fourni par la bibliothèque standard. Il se trouve sous le chemin `"container"` et propose trois sous-packages principaux :

---

### 🔹 1. `container/list`

Implémente une **liste doublement chaînée** (`doubly linked list`).

#### Principales caractéristiques :

* Chaque élément est un pointeur vers une struct `Element` contenant une valeur (`Value`) et des pointeurs vers l'élément précédent et suivant.
* Opérations rapides d’insertion/suppression n’importe où dans la liste.

#### Exemple d'utilisation :

```go
package main

import (
    "container/list"
    "fmt"
)

func main() {
    l := list.New()
    l.PushBack("A")
    l.PushFront("B")
    l.PushBack("C")

    for e := l.Front(); e != nil; e = e.Next() {
        fmt.Println(e.Value)
    }
}
```

---

### 🔹 2. `container/ring`

Implémente une **liste circulaire** de type anneau (`ring buffer` ou `circular list`).

#### Particularités :

* Chaque élément pointe vers le suivant, et le dernier pointe vers le premier.
* Idéal pour les structures circulaires ou les buffers tournants.

#### Exemple :

```go
package main

import (
    "container/ring"
    "fmt"
)

func main() {
    r := ring.New(3)
    for i := 0; i < 3; i++ {
        r.Value = i
        r = r.Next()
    }

    r.Do(func(p any) {
        fmt.Println(p)
    })
}
```

---

### 🔹 3. `container/heap`

Permet d’utiliser une **file de priorité** sous forme de **tas binaire**.

#### Fonctionne avec une interface (`heap.Interface`) qu’on doit implémenter soi-même (type slice + méthodes `Len`, `Less`, `Swap`, `Push`, `Pop`).

#### Exemple simplifié :

```go
package main

import (
    "container/heap"
    "fmt"
)

type IntHeap []int

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] } // min-heap
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *IntHeap) Push(x any) {
    *h = append(*h, x.(int))
}

func (h *IntHeap) Pop() any {
    old := *h
    n := len(old)
    x := old[n-1]
    *h = old[0 : n-1]
    return x
}

func main() {
    h := &IntHeap{3, 1, 4}
    heap.Init(h)
    heap.Push(h, 2)
    fmt.Println(heap.Pop(h)) // 1 (le plus petit)
}
```

---

### En résumé

| Sous-package     | Type de structure        | Utilisation typique              |
| ---------------- | ------------------------ | -------------------------------- |
| `container/list` | Liste doublement chaînée | Insertion/suppression fréquentes |
| `container/ring` | Liste circulaire         | Buffers circulaires              |
| `container/heap` | File de priorité (tas)   | Tri partiel, scheduling          |

Souhaites-tu un exemple d’utilisation pratique dans un projet ?

Parfait ! Voici un petit **exemple de projet Go structuré**, utilisant le package `container/heap` pour gérer une **file de priorité**, avec **tests unitaires**.

---

## 🧱 Structure du projet

```
priorityqueue/
├── go.mod
├── main.go
├── pq/
│   ├── heap.go
│   └── heap_test.go
```

---

## 📄 `go.mod`

```go
module priorityqueue

go 1.20
```

---

## 📄 `pq/heap.go` — Implémentation d'un min-heap

```go
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
```

---

## 📄 `pq/heap_test.go` — Tests unitaires

```go
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
```

---

## 📄 `main.go` — Exemple simple d'utilisation

```go
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
```

---

### ✅ Pour lancer les tests :

```bash
go test ./pq
```

---


Souhaites-tu un exemple d’utilisation pratique dans un cas d’usage particulier ?

Très bien ! Voici un **cas d’usage concret** d'une file de priorité (via `container/heap`) : **un système de gestion des tâches urgentes** dans une application métier — par exemple pour un centre d'appels, un hôpital ou un support technique.

On va :

1. Simuler un gestionnaire de tâches avec priorité.
2. Ajouter des tests unitaires pour vérifier l’ordre de traitement.

---

## 🎯 Cas d’usage : Gestionnaire de tickets support IT

Chaque ticket a une priorité :

* **0** = critique (serveur HS)
* **1** = haute (incident bloquant)
* **2+** = normal ou bas (requête utilisateur)

---

### 🧱 Structure du projet

```
support/
├── go.mod
├── main.go
├── support/
│   ├── ticket.go
│   └── ticket_test.go
```

---

### 📄 `go.mod`

```go
module support

go 1.20
```

---

### 📄 `support/ticket.go` — Le gestionnaire de tickets

```go
package support

import (
    "container/heap"
    "fmt"
)

type Ticket struct {
    ID       int
    Subject  string
    Priority int // 0 = critique, 1 = haute, ...
}

type TicketQueue []*Ticket

func (tq TicketQueue) Len() int { return len(tq) }

func (tq TicketQueue) Less(i, j int) bool {
    return tq[i].Priority < tq[j].Priority
}

func (tq TicketQueue) Swap(i, j int) {
    tq[i], tq[j] = tq[j], tq[i]
}

func (tq *TicketQueue) Push(x any) {
    *tq = append(*tq, x.(*Ticket))
}

func (tq *TicketQueue) Pop() any {
    old := *tq
    n := len(old)
    ticket := old[n-1]
    *tq = old[0 : n-1]
    return ticket
}

type Manager struct {
    queue *TicketQueue
}

func NewManager() *Manager {
    tq := make(TicketQueue, 0)
    heap.Init(&tq)
    return &Manager{queue: &tq}
}

func (m *Manager) AddTicket(ticket *Ticket) {
    heap.Push(m.queue, ticket)
}

func (m *Manager) NextTicket() *Ticket {
    if m.queue.Len() == 0 {
        return nil
    }
    return heap.Pop(m.queue).(*Ticket)
}

func (m *Manager) QueueSize() int {
    return m.queue.Len()
}

func (m *Manager) PrintAllTickets() {
    for _, t := range *m.queue {
        fmt.Printf("Ticket #%d: %s (P%d)\n", t.ID, t.Subject, t.Priority)
    }
}
```

---

### 📄 `support/ticket_test.go` — Tests du gestionnaire

```go
package support

import (
    "testing"
)

func TestTicketPriority(t *testing.T) {
    manager := NewManager()

    manager.AddTicket(&Ticket{ID: 1, Subject: "Email ne fonctionne plus", Priority: 2})
    manager.AddTicket(&Ticket{ID: 2, Subject: "Serveur inaccessible", Priority: 0})
    manager.AddTicket(&Ticket{ID: 3, Subject: "Demande de nouveau clavier", Priority: 3})
    manager.AddTicket(&Ticket{ID: 4, Subject: "Bug critique en prod", Priority: 1})

    expectedOrder := []int{2, 4, 1, 3} // IDs selon la priorité

    for _, expectedID := range expectedOrder {
        ticket := manager.NextTicket()
        if ticket == nil || ticket.ID != expectedID {
            t.Errorf("Expected ticket ID %d, got %v", expectedID, ticket)
        }
    }

    if manager.QueueSize() != 0 {
        t.Errorf("Queue should be empty, got size %d", manager.QueueSize())
    }
}
```

---

### 📄 `main.go` — Simulation de traitement de tickets

```go
package main

import (
    "fmt"
    "support/support"
)

func main() {
    manager := support.NewManager()

    manager.AddTicket(&support.Ticket{ID: 100, Subject: "PC ne démarre pas", Priority: 2})
    manager.AddTicket(&support.Ticket{ID: 101, Subject: "Fuite mémoire en prod", Priority: 1})
    manager.AddTicket(&support.Ticket{ID: 102, Subject: "Serveur mort", Priority: 0})

    fmt.Println("Traitement des tickets (ordre de priorité) :")
    for manager.QueueSize() > 0 {
        ticket := manager.NextTicket()
        fmt.Printf("-> Ticket #%d : %s [P%d]\n", ticket.ID, ticket.Subject, ticket.Priority)
    }
}
```

---

### ✅ Lancer les tests :

```bash
go test ./support
```

---
