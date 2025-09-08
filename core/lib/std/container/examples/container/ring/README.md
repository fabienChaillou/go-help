Bien sûr ! Voici un exemple simple d'utilisation du package `container/ring` en Go, avec un test unitaire inclus dans **un seul fichier** :

---

### ✅ Exemple avec `container/ring`

```go
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
```

---

### 🔧 Pour exécuter :

* Lancement du programme :

  ```bash
  go run fichier.go
  ```

* Exécution du test :

  ```bash
  go test fichier.go
  ```

Souhaites-tu un exemple avec un comportement circulaire ou avec des opérations plus avancées sur la ring ?

Parfait ! Voici un exemple de `container/ring` avec **comportement circulaire** — on fait tourner la ring, on modifie les valeurs, et on vérifie son comportement circulaire dans un **seul fichier**, avec test inclus.

---

### 🔁 Exemple circulaire avec test

```go
package main

import (
	"container/ring"
	"fmt"
	"testing"
)

// Crée une ring circulaire avec n éléments initialisés à 1, 2, ..., n
func CreateCircularRing(n int) *ring.Ring {
	r := ring.New(n)
	for i := 1; i <= n; i++ {
		r.Value = i
		r = r.Next()
	}
	return r
}

// Fait tourner la ring de k pas (positif = sens horaire, négatif = antihoraire)
func RotateRing(r *ring.Ring, k int) *ring.Ring {
	if k > 0 {
		return r.Move(k)
	}
	return r.Move(lenToMove(k))
}

// Utilitaire pour convertir une ring en slice depuis un point de départ donné
func RingToSlice(r *ring.Ring) []int {
	var result []int
	r.Do(func(v interface{}) {
		result = append(result, v.(int))
	})
	return result
}

// Convertit un déplacement négatif en déplacement positif équivalent
func lenToMove(k int) int {
	return ((k % 5) + 5) % 5 // Adapté pour une ring de taille 5
}

// Fonction principale : démonstration
func main() {
	r := CreateCircularRing(5) // [1 2 3 4 5]
	r = RotateRing(r, 2)       // Tourne de 2 vers la droite -> [3 4 5 1 2]

	r.Do(func(v interface{}) {
		fmt.Println(v)
	})
}

// Test du comportement circulaire
func TestRotateRing(t *testing.T) {
	r := CreateCircularRing(5)
	r = RotateRing(r, 2)
	got := RingToSlice(r)
	want := []int{3, 4, 5, 1, 2}

	for i := range got {
		if got[i] != want[i] {
			t.Errorf("index %d: expected %d, got %d", i, want[i], got[i])
		}
	}
}
```

---

### 📌 Explication rapide :

* `CreateCircularRing(5)` → crée `[1 2 3 4 5]`
* `RotateRing(r, 2)` → tourne deux pas vers la droite → devient `[3 4 5 1 2]`
* Le test vérifie ce comportement.

---

### 🧪 Lancer les tests :

```bash
go test fichier.go
```
