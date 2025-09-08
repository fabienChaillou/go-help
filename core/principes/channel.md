## Les channel

En Go (Golang), les “channels” sont un mécanisme fondamental du langage Go pour la concurrence.



Voici une explication simple des channels en Go :

Qu’est-ce qu’un channel ?


Un channel est un moyen de communication entre goroutines (les threads légers de Go). Il permet d’envoyer et de recevoir des valeurs entre goroutines, de manière sécurisée et synchronisée.

Syntaxe de base:


```go
// Création d’un channel :
ch := make(chan int) // Channel de type int

// Envoi d’une valeur dans un channel :
ch <- 42 // Envoie la valeur 42 dans le channel

// Réception d’une valeur depuis un channel :
value := <-ch // Lit une valeur depuis le channel
```

### Exemple simple

```go
package main

import "fmt"

func main() {
    ch := make(chan string)

    go func() {
        ch <- "Bonjour depuis la goroutine"
    }()

    message := <-ch
    fmt.Println(message)
}
```

Explication :

Une goroutine envoie un message dans le channel.

La fonction principale reçoit le message et l’affiche.

------------------------------------------------------------------------------------------------------------------------

Types de channels
    1. Unbuffered (sans tampon) : communication synchrone. L’envoi attend que quelqu’un reçoive.

    2. Buffered (avec tampon) : communication asynchrone jusqu’à remplir le buffer.

```go
ch := make(chan int, 2) // Channel avec buffer de 2
```

------------------------------------------------------------------------------------------------------------------------

Fermeture d’un channel

```go
close(ch)
```

------------------------------------------------------------------------------------------------------------------------

Souhaites-tu un exemple plus complexe avec plusieurs goroutines ou une explication visuelle ?

Parfait ! Voici un exemple plus complet montrant comment plusieurs goroutines peuvent communiquer entre elles via des channels. On va simuler un système simple où plusieurs “travailleurs” (workers) reçoivent des tâches via un channel.

Exemple : Pool de Workers

```go
package main

import (
	"fmt"
	"time"
)

// Le worker traite des jobs envoyés via le channel jobs
func worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Printf("Worker %d commence le job %d\n", id, j)
		time.Sleep(time.Second) // Simule un travail
		fmt.Printf("Worker %d termine le job %d\n", id, j)
		results <- j * 2 // Résultat du job
	}
}

func main() {
	const numJobs = 5
	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)

	// Démarre 3 workers
	for w := 1; w <= 3; w++ {
		go worker(w, jobs, results)
	}

	// Envoie les jobs
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs) // Important pour que les workers sachent quand s’arrêter

	// Récupère les résultats
	for a := 1; a <= numJobs; a++ {
		res := <-results
		fmt.Printf("Résultat reçu: %d\n", res)
	}
}
```


------------------------------------------------------------------------------------------------------------------------

Ce que ce programme fait :
    * Crée 3 goroutines (workers).

    * Envoie 5 tâches dans le channel jobs.

    * Chaque worker récupère un job, le traite, et renvoie un résultat dans le channel results.

    * Le programme principal attend les 5 résultats et les affiche.
