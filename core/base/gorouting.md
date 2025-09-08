En Go (Golang), une **goroutine** est une fonction qui s'exécute **concurrentement** avec d'autres goroutines. C’est l’un des piliers de la programmation concurrente en Go. Elles sont très **légères**, bien plus que les threads traditionnels, ce qui permet d’en exécuter des milliers sans problème.

---

### 💡 Définition simple :

> Une **goroutine** est une tâche (fonction) lancée en parallèle, gérée par le runtime de Go.

---

### 🧠 Syntaxe de base :

```go
go fonction()
```

Cela lance `fonction()` dans une nouvelle goroutine.

---

### 🔧 Exemple concret :

```go
package main

import (
    "fmt"
    "time"
)

func direBonjour() {
    fmt.Println("Bonjour !")
}

func main() {
    go direBonjour() // goroutine lancée
    fmt.Println("Programme principal")
    
    time.Sleep(1 * time.Second) // attend que la goroutine finisse (sinon elle peut être tuée trop vite)
}
```

**Sortie possible :**

```
Programme principal
Bonjour !
```

⚠️ Sans le `time.Sleep`, il est possible que le programme se termine avant que la goroutine n’ait eu le temps de s’exécuter.

---

### 🔄 Comparaison avec un appel normal :

| Appel classique         | Goroutine                 |
| ----------------------- | ------------------------- |
| `fonction()`            | `go fonction()`           |
| S’exécute immédiatement | S’exécute en arrière-plan |
| Bloque le reste         | Ne bloque pas le reste    |

---

### ⚙️ Utilisation typique :

* Traitement en parallèle (API, fichiers, calculs)
* Serveurs web (chaque requête dans une goroutine)
* Producteurs/consommateurs avec des **channels**

---

### 🧵 Et les threads ?

* Une goroutine **n’est pas** un thread, mais elle peut utiliser un ou plusieurs threads gérés automatiquement par le **runtime Go**.
* Le runtime de Go fait du **multiplexage** des goroutines sur un petit nombre de threads systèmes.

---

Souhaites-tu un exemple avec des **channels** pour communiquer entre goroutines ?

Parfait ! Voici un **exemple simple** où deux goroutines communiquent via un **channel**, un mécanisme clé de Go pour synchroniser et transmettre des données entre goroutines.

---

### 🔁 Exemple : communication entre goroutines avec un channel

```go
package main

import (
    "fmt"
)

func direBonjour(ch chan string) {
    ch <- "Bonjour depuis la goroutine !" // envoie un message dans le channel
}

func main() {
    ch := make(chan string) // création d'un channel de type string

    go direBonjour(ch) // lance la goroutine

    msg := <-ch // lit (bloquant) un message depuis le channel
    fmt.Println(msg)
}
```

---

### 🧠 Ce qu’il se passe :

1. On crée un **channel** `ch` capable de transporter des `string`.
2. On lance `direBonjour` dans une **goroutine**, en lui passant le channel.
3. La fonction `direBonjour` envoie `"Bonjour depuis la goroutine !"` dans le channel avec `ch <-`.
4. Le `main()` attend ce message avec `<-ch` et l'affiche.

---

### 🔒 Pourquoi c’est utile ?

* Le channel **synchronise** les goroutines : la lecture `<-ch` attend que quelqu’un écrive.
* Tu peux transmettre n’importe quel type : `int`, `struct`, `bool`, etc.

---

### 🚀 Un exemple un peu plus poussé :

```go
package main

import (
    "fmt"
    "time"
)

func count(name string, ch chan string) {
    for i := 1; i <= 3; i++ {
        ch <- fmt.Sprintf("%s: %d", name, i)
        time.Sleep(500 * time.Millisecond)
    }
}

func main() {
    ch := make(chan string)

    go count("A", ch)
    go count("B", ch)

    // On lit 6 messages (3 de A + 3 de B)
    for i := 0; i < 6; i++ {
        fmt.Println(<-ch)
    }
}
```

💬 **Sortie possible (ordre non garanti)** :

```
A: 1
B: 1
A: 2
B: 2
A: 3
B: 3
```

---

Souhaites-tu que je t’explique aussi les **channels bufferisés**, ou les **selects** pour gérer plusieurs channels ?

Parfait ! Commençons par les deux sujets que tu veux découvrir : **channels bufferisés** et le mot-clé **`select`**.

---

## 📦 1. Channels bufferisés

Un **channel bufferisé** permet d’envoyer plusieurs messages sans attendre immédiatement qu’ils soient lus. Il a une **capacité** (le "buffer").

### 🔧 Syntaxe :

```go
ch := make(chan int, 3) // channel bufferisé de capacité 3
```

Tu peux y mettre jusqu’à 3 messages **sans que personne ne les lise tout de suite**.

---

### ✅ Exemple :

```go
package main

import "fmt"

func main() {
    ch := make(chan int, 2) // buffer de 2

    ch <- 1
    ch <- 2
    // ch <- 3 // bloquerait ici car le buffer est plein

    fmt.Println(<-ch) // lit 1
    fmt.Println(<-ch) // lit 2
}
```

### 🧠 Quand l’utiliser ?

* Quand tu veux **envoyer plusieurs valeurs d’un coup** sans bloquer l’émetteur.
* Très utile pour faire des **pipelines** de traitement.

---

## 🔀 2. `select` : écouter plusieurs channels à la fois

`select` permet d’**attendre sur plusieurs channels** en même temps : la première opération qui devient prête est exécutée.

---

### ✅ Exemple simple :

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    ch1 := make(chan string)
    ch2 := make(chan string)

    go func() {
        time.Sleep(1 * time.Second)
        ch1 <- "Message de ch1"
    }()

    go func() {
        time.Sleep(2 * time.Second)
        ch2 <- "Message de ch2"
    }()

    select {
    case msg1 := <-ch1:
        fmt.Println("Reçu :", msg1)
    case msg2 := <-ch2:
        fmt.Println("Reçu :", msg2)
    }
}
```

### 🔍 Ce qu’il se passe :

* `select` attend que **ch1 ou ch2** ait un message.
* La première goroutine envoie après 1s, donc `select` lit ce message-là.
* Le second message est ignoré ici (mais tu peux faire une boucle si tu veux lire plusieurs).

---

### ➕ Avec `default` :

```go
select {
case msg := <-ch:
    fmt.Println("Message reçu :", msg)
default:
    fmt.Println("Aucun message, on continue")
}
```

➡️ Permet d’**éviter le blocage** si aucun channel n’est prêt.

---

Souhaites-tu maintenant un exemple **combiné** de goroutines, channels bufferisés et `select`, comme un **système de workers ou de file de traitement** ?

Parfait ! Voici un **exemple complet et testé** d’un petit **système de "workers" concurrents** avec :

* Plusieurs **goroutines workers** qui traitent des jobs,
* Un **channel bufferisé** pour distribuer les jobs,
* Un **`select`** pour la gestion des signaux d’arrêt ou de timeout,
* Et des **tests unitaires** pour vérifier le fonctionnement.

---

## 🛠️ Structure du système

On va créer :

* Un `Job` simple (un int à doubler),
* Un `WorkerPool` qui gère l’envoi des jobs aux workers,
* Un test pour vérifier que les résultats sont corrects.

---

### 📁 `main.go` (ou `worker.go` si tu préfères le séparer)

```go
package worker

import (
    "fmt"
    "sync"
    "time"
)

type Job struct {
    ID     int
    Input  int
    Output chan int
}

type WorkerPool struct {
    JobQueue chan Job
    Workers  int
    wg       sync.WaitGroup
}

func NewWorkerPool(workerCount, queueSize int) *WorkerPool {
    return &WorkerPool{
        JobQueue: make(chan Job, queueSize),
        Workers:  workerCount,
    }
}

func (wp *WorkerPool) Start() {
    for i := 0; i < wp.Workers; i++ {
        wp.wg.Add(1)
        go func(id int) {
            defer wp.wg.Done()
            for job := range wp.JobQueue {
                // Traitement simulé
                time.Sleep(50 * time.Millisecond)
                result := job.Input * 2
                fmt.Printf("Worker %d a traité job %d : %d → %d\n", id, job.ID, job.Input, result)
                job.Output <- result
            }
        }(i)
    }
}

func (wp *WorkerPool) Stop() {
    close(wp.JobQueue)
    wp.wg.Wait()
}
```

---

### 🧪 `worker_test.go`

```go
package worker

import (
    "testing"
    "time"
)

func TestWorkerPool(t *testing.T) {
    wp := NewWorkerPool(3, 5)
    wp.Start()

    jobs := []Job{}
    results := []int{2, 4, 6, 8, 10}

    for i := 1; i <= 5; i++ {
        ch := make(chan int, 1)
        jobs = append(jobs, Job{
            ID:     i,
            Input:  i,
            Output: ch,
        })
    }

    // Envoi des jobs
    for _, job := range jobs {
        wp.JobQueue <- job
    }

    // Récupération des résultats avec timeout de sécurité
    for i, job := range jobs {
        select {
        case res := <-job.Output:
            if res != results[i] {
                t.Errorf("Job %d: attendu %d, obtenu %d", job.ID, results[i], res)
            }
        case <-time.After(1 * time.Second):
            t.Errorf("Job %d: timeout de lecture", job.ID)
        }
    }

    wp.Stop()
}
```

---

### ✅ Pour tester :

1. Place `main.go` et `worker_test.go` dans le même package (ex : `worker`).
2. Lance les tests avec :

```bash
go test -v
```

---
