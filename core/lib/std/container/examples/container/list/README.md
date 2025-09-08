Pour tester ce code Go (qui manipule une liste chaînée de la bibliothèque standard `container/list`), on peut extraire la logique dans une fonction que l’on pourra tester plus facilement. Voici un exemple de refactorisation et d’un test unitaire associé :

### ✅ Refactorisation du code

```go
// list_util.go
package main

import (
	"container/list"
)

func CreateList() *list.List {
	l := list.New()
	l.PushBack("A")
	l.PushFront("B")
	l.PushBack("C")
	return l
}

func ListToSlice(l *list.List) []string {
	var result []string
	for e := l.Front(); e != nil; e = e.Next() {
		result = append(result, e.Value.(string))
	}
	return result
}
```

### 🧪 Fichier de test

```go
// list_util_test.go
package main

import (
	"testing"
)

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
```

### 🔧 Pour exécuter les tests

Dans le terminal :

```bash
go test
```
