En Go, *determining concurrency safety* (déterminer la sécurité en concurrence) signifie s'assurer qu'un morceau de code ou une structure de données peut être utilisé correctement par plusieurs goroutines en même temps sans provoquer de **conditions de course** (*race conditions*) ou d'autres comportements indéterminés.

### ⚠️ Qu’est-ce qu’un problème de concurrence ?

Quand deux ou plusieurs goroutines accèdent ou modifient une même donnée partagée sans synchronisation correcte, le résultat devient imprévisible. Cela peut entraîner des bugs subtils, souvent difficiles à reproduire.

---

## 🔍 Comment déterminer si un code est **concurrent-safe** en Go ?

### 1. **Lire la documentation des types utilisés**

Certains types sont explicitement conçus pour être thread-safe. Par exemple :

* `sync.Mutex`, `sync.RWMutex`, `sync.Map`, `sync.Once`, etc., sont sûrs pour une utilisation concurrente.
* `map` **n’est pas** sûr pour la concurrence par défaut !

### 2. **Analyser les accès à la mémoire partagée**

Demande-toi :

* Est-ce que plusieurs goroutines peuvent lire ou écrire la même variable ?
* Y a-t-il un verrou (mutex) ou une autre forme de synchronisation pour réguler l’accès ?
* Y a-t-il des canaux utilisés pour séquencer les accès ou passer les données sans partage direct ?

### 3. **Utiliser l’outil de détection de race de Go**

Go fournit un outil intégré pour détecter les race conditions :

```bash
go run -race main.go
```

ou pour les tests :

```bash
go test -race ./...
```

Cela affiche les conflits d'accès mémoire entre goroutines pendant l'exécution.

---

## ✅ Bonnes pratiques pour assurer la sécurité en concurrence

* Utiliser des **mutex** pour protéger les variables partagées :

```go
var mu sync.Mutex
var count int

func increment() {
    mu.Lock()
    count++
    mu.Unlock()
}
```

* Utiliser des **channels** pour communiquer entre goroutines sans partager de mémoire directement :

```go
func worker(jobs <-chan int, results chan<- int) {
    for job := range jobs {
        results <- job * 2
    }
}
```

* Éviter de modifier des **slices**, des **maps**, ou des **structs** partagées sans protection.
* Lire la documentation des bibliothèques externes : certaines structures comme `http.Client` ou `sql.DB` sont conçues pour être utilisées de manière concurrente, d'autres non.

---

## 🧠 En résumé :

Déterminer si ton code est **concurrent-safe** revient à :

* Repérer les accès concurrents à une même donnée.
* Vérifier qu’ils sont protégés par des mécanismes de synchronisation.
* Tester avec `-race` pour détecter les problèmes potentiels.

Souhaites-tu un exemple simple illustrant un code non thread-safe et sa version sécurisée ?

Voici un **exemple simple** pour illustrer la différence entre un code **non thread-safe** et un code **thread-safe** en Go.

---

## ❌ Exemple **non thread-safe** (avec condition de course) :

```go
package main

import (
	"fmt"
	"time"
)

var counter int

func increment() {
	for i := 0; i < 1000; i++ {
		counter++ // accès concurrent non protégé !
	}
}

func main() {
	for i := 0; i < 5; i++ {
		go increment()
	}

	time.Sleep(1 * time.Second) // attendre que les goroutines terminent
	fmt.Println("Final counter:", counter)
}
```

### 🛑 Problème :

* Plusieurs goroutines modifient `counter` en même temps.
* Résultat final **non déterministe**.
* Exécute-le avec `go run -race` → tu verras un avertissement de *race condition*.

---

## ✅ Version **thread-safe** avec `sync.Mutex` :

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	counter int
	mu      sync.Mutex
)

func increment() {
	for i := 0; i < 1000; i++ {
		mu.Lock()
		counter++
		mu.Unlock()
	}
}

func main() {
	for i := 0; i < 5; i++ {
		go increment()
	}

	time.Sleep(1 * time.Second)
	fmt.Println("Final counter:", counter)
}
```

### ✅ Avantages :

* `sync.Mutex` garantit qu’une seule goroutine accède à `counter` à la fois.
* Résultat final correct : `5 * 1000 = 5000`.

---

Souhaites-tu aussi voir une version utilisant des **channels** pour éviter le partage direct de la variable ?

Voici une version **concurrent-safe** du même exemple, cette fois en utilisant des **channels**, ce qui est souvent plus idiomatique en Go.

---

## ✅ Version **concurrent-safe avec `chan`** (pas besoin de mutex) :

```go
package main

import (
	"fmt"
	"time"
)

func increment(jobs <-chan int, results chan<- int) {
	count := 0
	for job := range jobs {
		count += job
	}
	results <- count
}

func main() {
	jobs := make(chan int)
	results := make(chan int)

	// Lancer une seule goroutine responsable du comptage
	go increment(jobs, results)

	// Lancer plusieurs workers qui envoient des incréments
	for i := 0; i < 5; i++ {
		go func() {
			for j := 0; j < 1000; j++ {
				jobs <- 1
			}
		}()
	}

	// Attendre que toutes les valeurs soient envoyées
	time.Sleep(1 * time.Second)
	close(jobs) // important : fermer le canal pour terminer la boucle dans `increment`

	// Recevoir le résultat final
	total := <-results
	fmt.Println("Final counter:", total)
}
```

### ✅ Avantages :

* Aucune mémoire partagée entre les goroutines.
* Le comptage est centralisé dans une seule goroutine (`increment`), qui lit depuis un canal.
* Ce modèle suit le principe Go : **"ne pas communiquer en partageant la mémoire, mais partager la mémoire en communiquant."**

Souhaites-tu un schéma pour visualiser comment les goroutines et les channels interagissent ici ?
