En Go, le **buffer** fait généralement référence à un espace mémoire temporaire utilisé pour **stocker des données avant traitement ou transmission**. Le concept de buffer est souvent utilisé avec :

1. **Les channels bufferisés** (concurrence)
2. **Les buffers d'entrée/sortie** (fichiers, réseaux, etc.)

---

## 1. 🔹 Buffer dans les **channels**

Un **channel bufferisé** permet de **stocker temporairement plusieurs valeurs**, contrairement à un channel non bufferisé (qui bloque jusqu’à réception).

### ➤ Création d’un channel bufferisé :

```go
ch := make(chan int, 3) // channel de type int avec un buffer de 3
```

### ➤ Fonctionnement :

* Tu peux envoyer jusqu’à 3 valeurs dans le channel sans que personne ne les lise.
* Au 4e envoi, l'envoi **bloque** jusqu'à ce qu'une valeur soit lue.

### ➤ Exemple :

```go
func main() {
    ch := make(chan int, 2)

    ch <- 1
    ch <- 2
    // ch <- 3 // ceci bloquerait car le buffer est plein

    fmt.Println(<-ch) // lit 1
    fmt.Println(<-ch) // lit 2
}
```

---

## 2. 🔹 Buffer en **entrée/sortie (I/O)**

Go fournit des buffers pour optimiser les lectures/écritures sur les fichiers ou réseaux, via le package `bufio`.

### ➤ Exemple : écriture bufferisée dans un fichier

```go
import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    fichier, _ := os.Create("exemple.txt")
    writer := bufio.NewWriter(fichier)

    writer.WriteString("Hello, monde !\n") // écrit dans le buffer, pas directement dans le fichier
    writer.Flush()                        // vide le buffer dans le fichier

    fichier.Close()
}
```

🔸 `bufio.NewWriter` crée un buffer d’écriture.
🔸 `Flush()` est **obligatoire** pour forcer l’écriture dans le fichier.

---

## 🔹 Pourquoi utiliser des buffers ?

* ✅ **Meilleure performance** : moins d'accès directs au disque ou aux threads.
* ✅ **Moins de blocage** avec des channels bufferisés.
* ⚠️ Attention à bien **vider les buffers** (`Flush()`), sinon des données peuvent rester en mémoire.

---

## Example comparatif entre channel bufferisé et non bufferisé dans un contexte de producteur/consommateur ?

Voici un exemple clair et concis en Go illustrant la différence entre un **channel bufferisé** et un **channel non bufferisé** dans un **modèle producteur/consommateur**.

---

## 🔸 1. Channel **non bufferisé**

Dans ce cas, chaque envoi (`chan <- value`) **bloque** jusqu'à ce qu'un consommateur le reçoive (`<-chan`).

```go
package main

import (
	"fmt"
	"time"
)

func producer(ch chan int) {
	for i := 1; i <= 5; i++ {
		fmt.Println("Producteur envoie :", i)
		ch <- i // bloque jusqu'à ce que le consommateur reçoive
	}
	close(ch)
}

func consumer(ch chan int) {
	for v := range ch {
		fmt.Println("Consommateur reçoit :", v)
		time.Sleep(time.Second) // simule un traitement lent
	}
}

func main() {
	ch := make(chan int) // non bufferisé

	go producer(ch)
	consumer(ch)
}
```

**Sortie (approximative)** :

```
Producteur envoie : 1
Consommateur reçoit : 1
Producteur envoie : 2
Consommateur reçoit : 2
...
```

🟡 Le producteur attend que chaque valeur soit lue avant d’en envoyer une autre.

---

## 🔹 2. Channel **bufferisé**

Avec un buffer, le producteur peut envoyer plusieurs messages **sans attendre immédiatement** qu’ils soient consommés, **jusqu'à la taille du buffer**.

```go
func main() {
	ch := make(chan int, 3) // buffer de taille 3

	go producer(ch)
	consumer(ch)
}
```

**Sortie (approximative)** :

```
Producteur envoie : 1
Producteur envoie : 2
Producteur envoie : 3
Consommateur reçoit : 1
Producteur envoie : 4
Consommateur reçoit : 2
...
```

🔵 Le producteur peut envoyer jusqu’à 3 valeurs **rapidement**, sans attendre. Ensuite, il est bloqué jusqu'à ce que le consommateur ait vidé une place.

---

## ✅ Conclusion

| Aspect           | Non Bufferisé                      | Bufferisé (ex: taille 3)        |
| ---------------- | ---------------------------------- | ------------------------------- |
| Envoi bloquant ? | Oui (jusqu’à lecture)              | Non (jusqu’à remplir le buffer) |
| Couplage         | Fort entre producteur/consommateur | Plus flexible                   |
| Idéal pour       | Synchronisation stricte            | Performance / débit             |

## Example avec `sync.WaitGroup`

Parfait ! Voici un exemple complet en Go avec :

* Deux versions : **channel non bufferisé** et **bufferisé**.
* Utilisation de `sync.WaitGroup` pour attendre la fin des goroutines.
* Un test unitaire simple avec `testing`.

---

## 📁 Structure du projet (simplifiée)

```
/producerconsumer
├── main.go
└── main_test.go
```

---

## ✅ `main.go` – Producteur/Consommateur avec WaitGroup

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

func producer(id int, ch chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 1; i <= 3; i++ {
		val := id*10 + i
		fmt.Printf("Producteur %d envoie : %d\n", id, val)
		ch <- val
	}
}

func consumer(id int, ch chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for val := range ch {
		fmt.Printf("Consommateur %d reçoit : %d\n", id, val)
		time.Sleep(200 * time.Millisecond) // simulation traitement
	}
}

func run(buffered bool) {
	var ch chan int
	if buffered {
		ch = make(chan int, 5) // bufferisé
	} else {
		ch = make(chan int) // non bufferisé
	}

	var wg sync.WaitGroup
	wg.Add(3) // 2 producteurs, 1 consommateur

	go producer(1, ch, &wg)
	go producer(2, ch, &wg)

	go func() {
		wg.Wait() // attend la fin des producteurs
		close(ch) // ferme le canal
	}()

	wg.Add(1) // consommateur
	consumer(1, ch, &wg)
}
```

---

## 🧪 `main_test.go` – Test unitaire

```go
package main

import "testing"

func TestRunWithBufferedChannel(t *testing.T) {
	run(true)
}

func TestRunWithUnbufferedChannel(t *testing.T) {
	run(false)
}
```

---

## 🏁 Pour exécuter

Dans le terminal, exécute :

```bash
go test -v
```

---

## 📌 Résultat attendu

* Le **comportement diffère** selon que le canal est bufferisé ou non :

  * Non bufferisé : envoi/lecture très séquentielle.
  * Bufferisé : les producteurs envoient plus rapidement jusqu’à saturation du buffer.


## ajoute un benchmark

Voici comment ajouter un **benchmark** en Go pour comparer les performances entre un **channel bufferisé** et **non bufferisé** dans ton scénario producteur/consommateur.

---

## ✅ Mise à jour de `main.go`

Ajoute cette fonction pour permettre de configurer dynamiquement le nombre de producteurs et consommateurs :

```go
package main

import (
	"sync"
	"time"
)

func BenchmarkRun(buffered bool, numProducers, numConsumers, valuesPerProducer int) {
	var ch chan int
	if buffered {
		ch = make(chan int, 100) // buffer assez grand
	} else {
		ch = make(chan int)
	}

	var wg sync.WaitGroup

	// Producteurs
	for p := 1; p <= numProducers; p++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for i := 0; i < valuesPerProducer; i++ {
				ch <- id*1000 + i
			}
		}(p)
	}

	// Fermeture du canal après tous les producteurs
	go func() {
		wg.Wait()
		close(ch)
	}()

	// Consommateurs
	var consumerWg sync.WaitGroup
	for c := 1; c <= numConsumers; c++ {
		consumerWg.Add(1)
		go func() {
			defer consumerWg.Done()
			for range ch {
				// simule un traitement très léger
				time.Sleep(1 * time.Millisecond)
			}
		}()
	}

	consumerWg.Wait()
}
```

---

## 🧪 Ajoute à `main_test.go`

```go
package main

import "testing"

func BenchmarkUnbuffered(b *testing.B) {
	for i := 0; i < b.N; i++ {
		BenchmarkRun(false, 2, 2, 100)
	}
}

func BenchmarkBuffered(b *testing.B) {
	for i := 0; i < b.N; i++ {
		BenchmarkRun(true, 2, 2, 100)
	}
}
```

---

## 🚀 Exécuter les benchmarks

Dans le terminal, lance :

```bash
go test -bench=.
```

---

## 🧾 Résultat attendu (exemple)

```
goos: linux
goarch: amd64
BenchmarkUnbuffered-8   	      10	 150000000 ns/op
BenchmarkBuffered-8     	      20	  80000000 ns/op
```

👉 Le canal **bufferisé** est souvent plus rapide dans ce genre de tâche car les producteurs ne bloquent pas immédiatement.

---

Tu peux ajouter l'option `-benchmem` pour obtenir des statistiques de **performance mémoire** lors des benchmarks. Voici comment faire :

---

## ✅ Étape 1 : Lancer le benchmark avec mémoire

Dans ton terminal :

```bash
go test -bench=. -benchmem
```

---

## 📊 Exemple de sortie attendue

```bash
goos: linux
goarch: amd64
BenchmarkUnbuffered-8   	       5	 245000000 ns/op	123456 B/op	  234 allocs/op
BenchmarkBuffered-8     	      10	 135000000 ns/op	 98765 B/op	  180 allocs/op
```

### Légende :

* `ns/op` : temps d'exécution moyen par itération (nanosecondes).
* `B/op` : nombre de **bytes alloués** par itération.
* `allocs/op` : nombre d’**allocations mémoire** par itération.

---

## 🧠 Analyse

| Critère        | Non bufferisé                           | Bufferisé             |
| -------------- | --------------------------------------- | --------------------- |
| Temps (ns/op)  | Plus lent                               | Plus rapide           |
| Mémoire (B/op) | Généralement plus haut                  | Généralement plus bas |
| Allocs/op      | Plus de synchronisation ⇒ plus d'allocs | Moins d’allocations   |

---
