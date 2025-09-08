Voici un **résumé détaillé** du livre *"Concurrency in Go: Tools and Techniques for Developers"* de **Katherine Cox-Buday**, accompagné d’**exemples en Go** pour illustrer les concepts principaux.

---

## 📘 Résumé global du livre

**Objectif** : Ce livre enseigne comment écrire des programmes concurrents sûrs, lisibles et performants en Go. Il couvre les fondements de la concurrence, les primitives de synchronisation de Go (goroutines, channels, mutex, etc.), les patterns de concurrence, et les considérations de design pour des systèmes concurrentiels robustes.

---

## 🧠 Chapitre par chapitre (avec exemples)

---

### **1. Pourquoi la concurrence ?**

* La concurrence permet d'exécuter plusieurs tâches de façon indépendante.
* Go est conçu avec la concurrence comme pilier central.

> Exemple :

```go
go fmt.Println("Concurrent Hello") // s'exécute dans une goroutine
```

---

### **2. Goroutines et planification**

* Une **goroutine** est une tâche légère gérée par le runtime Go.
* Elles sont **multiplexées** sur des threads OS par le planificateur.

> Exemple :

```go
func main() {
    for i := 0; i < 5; i++ {
        go fmt.Println("Hello from goroutine", i)
    }
    time.Sleep(time.Second) // Donne le temps aux goroutines de s'exécuter
}
```

---

### **3. Les canaux (channels)**

* Les **channels** permettent la communication entre goroutines.
* Utilisent le modèle CSP (Communicating Sequential Processes).

> Exemple :

```go
func main() {
    ch := make(chan string)
    go func() {
        ch <- "ping"
    }()
    msg := <-ch
    fmt.Println(msg)
}
```

---

### **4. Synchronisation et verrouillage**

* Go propose aussi des **mutexes**, **RWMutexes**, **WaitGroups** pour la synchronisation.

> Exemple avec `sync.Mutex` :

```go
var mu sync.Mutex
var count int

func increment() {
    mu.Lock()
    count++
    mu.Unlock()
}
```

> Exemple avec `sync.WaitGroup` :

```go
var wg sync.WaitGroup

func worker(id int) {
    defer wg.Done()
    fmt.Println("Worker", id)
}

func main() {
    for i := 0; i < 3; i++ {
        wg.Add(1)
        go worker(i)
    }
    wg.Wait()
}
```

---

### **5. Patterns de concurrence**

* **Fan-out / Fan-in** : Plusieurs goroutines traitent des tâches en parallèle, puis leurs résultats sont regroupés.

> Exemple :

```go
func worker(id int, jobs <-chan int, results chan<- int) {
    for j := range jobs {
        results <- j * 2
    }
}

func main() {
    jobs := make(chan int, 5)
    results := make(chan int, 5)

    for w := 1; w <= 3; w++ {
        go worker(w, jobs, results)
    }

    for j := 1; j <= 5; j++ {
        jobs <- j
    }
    close(jobs)

    for a := 1; a <= 5; a++ {
        fmt.Println(<-results)
    }
}
```

* **Pipeline** : Transformer les données à travers plusieurs étapes de traitement.

---

### **6. Annulation, délais et timeouts**

* Utilisation de `context.Context` pour annuler des opérations.

> Exemple :

```go
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
defer cancel()

ch := make(chan string)

go func() {
    time.Sleep(3 * time.Second)
    ch <- "done"
}()

select {
case <-ctx.Done():
    fmt.Println("Timeout")
case msg := <-ch:
    fmt.Println(msg)
}
```

---

### **7. Erreurs courantes à éviter**

* Fuites de goroutines (goroutines qui ne se terminent jamais).
* Deadlocks (attente mutuelle infinie).
* Race conditions.

> Exemple de **race condition** (à éviter) :

```go
var counter int

func main() {
    for i := 0; i < 1000; i++ {
        go func() {
            counter++
        }()
    }
    time.Sleep(time.Second)
    fmt.Println(counter) // résultat non déterministe
}
```

---

### **8. Tests et benchmarks concurrents**

* Go permet de tester la sécurité concurrente avec `-race`.
* Les benchmarks mesurent la performance sous différentes charges.

---

## ✅ En résumé

| Concept             | Go outil                          |
| ------------------- | --------------------------------- |
| Tâches légères      | `goroutines`                      |
| Communication       | `channels`                        |
| Verrouillage        | `sync.Mutex`, `RWMutex`           |
| Synchronisation     | `sync.WaitGroup`, `sync.Once`     |
| Gestion de contexte | `context.Context`                 |
| Patterns            | Fan-in/out, pipeline, worker pool |
| Sécurité            | `go test -race`                   |

---

Souhaites-tu un résumé visuel (ex. : schéma) ou une version PDF de ce résumé ?

