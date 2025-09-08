En Go, un **`sync.WaitGroup`** est un outil de synchronisation qui permet d’attendre que plusieurs goroutines aient terminé leur travail.

### À quoi ça sert ?

Quand tu lances plusieurs goroutines, tu ne sais pas à quel moment elles vont se terminer. Le `WaitGroup` te permet de dire au programme principal :
**« Attends que toutes les goroutines aient fini avant de continuer. »**

---

### Fonctions principales :

Voici les trois méthodes essentielles de `sync.WaitGroup` :

* **`Add(n int)`** : indique qu’on attend `n` goroutines.
* **`Done()`** : appelée à la fin d’une goroutine pour signaler qu’elle est terminée.
* **`Wait()`** : bloque l'exécution jusqu'à ce que le compteur atteigne zéro (c’est-à-dire que toutes les goroutines aient appelé `Done()`).

---

### Exemple simple :

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

func worker(id int, wg *sync.WaitGroup) {
    defer wg.Done() // Signale qu’on a fini à la fin de la goroutine

    fmt.Printf("Worker %d started\n", id)
    time.Sleep(1 * time.Second)
    fmt.Printf("Worker %d finished\n", id)
}

func main() {
    var wg sync.WaitGroup

    for i := 1; i <= 3; i++ {
        wg.Add(1) // On attend une goroutine de plus
        go worker(i, &wg)
    }

    wg.Wait() // On attend que toutes les goroutines aient fini
    fmt.Println("All workers done.")
}
```

---

### Ce qu’il faut retenir :

* Toujours faire correspondre **`Add(1)`** avec **`Done()`**.
* Utilise **`defer wg.Done()`** au début d’une goroutine pour éviter les oublis.
* Le `WaitGroup` ne peut pas être réutilisé une fois que `Wait()` a été appelée si d'autres `Add()` sont encore en cours après.

Souhaites-tu un exemple plus avancé (par exemple avec erreurs ou des goroutines parallèles sur des données) ?

### avec plusieurs gorouting parallèles

Parfait ! Voici un exemple plus avancé avec un `WaitGroup` où plusieurs goroutines travaillent **en parallèle** sur des **données partagées** (un tableau d’entiers à doubler). On va aussi utiliser un `mutex` pour protéger l’accès concurrent aux données, car plusieurs goroutines peuvent tenter d’écrire en même temps.

---

### 🔧 Exemple : Traitement parallèle d’un tableau d’entiers

```go
package main

import (
    "fmt"
    "sync"
)

func doubleWorker(id int, numbers []int, results map[int]int, wg *sync.WaitGroup, mu *sync.Mutex) {
    defer wg.Done()

    for _, n := range numbers {
        doubled := n * 2

        // Accès concurrent => protection avec un mutex
        mu.Lock()
        results[n] = doubled
        mu.Unlock()
    }

    fmt.Printf("Worker %d done\n", id)
}

func main() {
    numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
    results := make(map[int]int)

    var wg sync.WaitGroup
    var mu sync.Mutex

    // Diviser les données en 2 blocs pour 2 workers
    chunkSize := len(numbers) / 2
    for i := 0; i < 2; i++ {
        start := i * chunkSize
        end := start + chunkSize
        if i == 1 { // gérer le cas impair
            end = len(numbers)
        }

        wg.Add(1)
        go doubleWorker(i+1, numbers[start:end], results, &wg, &mu)
    }

    wg.Wait()

    // Affichage des résultats
    fmt.Println("Résultats :")
    for k, v := range results {
        fmt.Printf("%d -> %d\n", k, v)
    }
}
```

---

### ✅ Ce que cet exemple montre :

* Plusieurs goroutines traitent une partie d’un tableau (division des tâches).
* Un **`mutex`** protège la **map partagée `results`**, évitant des erreurs de concurrence.
* Le `WaitGroup` garantit que **le programme attend que tous les workers aient fini** avant d’afficher les résultats.
