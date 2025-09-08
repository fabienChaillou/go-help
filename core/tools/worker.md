Voici un **mini projet en Go** illustrant l’utilisation des **workers**, **channels**, **signals** et **goroutines**. L'objectif est de simuler un système de traitement de tâches (jobs) en parallèle, avec arrêt propre (grâce aux signaux système comme `SIGINT`).

---

### 🔧 Fonctionnalités :

* Plusieurs workers (goroutines) écoutent un `jobQueue` pour exécuter des tâches.
* Utilisation de channels pour la communication.
* Gestion du shutdown propre avec `os.Signal`.

---

### 🗂️ Structure :

```
worker-pool/
│
├── main.go
```

---

### 📄 `main.go`

```go
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Job représente une tâche à exécuter
type Job struct {
	ID int
}

// worker traite les jobs depuis le channel jobQueue
func worker(ctx context.Context, id int, jobQueue <-chan Job, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Worker %d started\n", id)

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Worker %d received shutdown signal\n", id)
			return
		case job, ok := <-jobQueue:
			if !ok {
				fmt.Printf("Worker %d: job queue closed\n", id)
				return
			}
			fmt.Printf("Worker %d: processing job %d\n", id, job.ID)
			time.Sleep(1 * time.Second) // simule une tâche longue
		}
	}
}

func main() {
	const numWorkers = 3
	const numJobs = 10

	jobQueue := make(chan Job, numJobs)
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	// Interception du signal système pour arrêt propre
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Démarre les workers
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(ctx, i, jobQueue, &wg)
	}

	// Envoie des jobs
	go func() {
		for j := 1; j <= numJobs; j++ {
			jobQueue <- Job{ID: j}
		}
		close(jobQueue) // Fermeture après tous les jobs
	}()

	// Attente du signal système
	<-sigChan
	fmt.Println("\nReceived shutdown signal")
	cancel()

	// Attente de la fin des workers
	wg.Wait()
	fmt.Println("All workers stopped. Exiting.")
}
```

---

### ▶️ Exécution

Compile et lance :

```bash
go run main.go
```

Puis presse `CTRL+C` à tout moment pour voir l'arrêt propre.

---

Souhaites-tu que ce projet soit étendu (ex: lire depuis une base de données, API HTTP, gestion d’erreurs plus poussée) ?

Parfait ! Voici une **version étendue** du projet avec les ajouts suivants :

---

### ✅ Objectif Étendu :

* Ajouter une **API HTTP en Go (avec Gin)** pour soumettre des jobs.
* Enregistrer les jobs dans une **base SQLite**.
* Utiliser un **worker pool** pour traiter les jobs enregistrés.
* Arrêt propre via `SIGINT` / `SIGTERM`.

---

### 🛠️ Stack utilisée :

* [Gin](https://github.com/gin-gonic/gin) – HTTP router.
* [SQLite](https://www.sqlite.org/) – Base de données légère.
* \[Goroutines + Channels] – pour le traitement en fond.
* [sqlx](https://github.com/jmoiron/sqlx) – wrapper SQL pratique.

---

### 🗂️ Structure du projet

```
go-worker-api/
├── main.go
├── db.go
├── handler.go
├── job.go
```

---

### 📄 `main.go`

```go
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	db := InitDB("jobs.db")
	defer db.Close()

	jobQueue := make(chan Job, 100)
	ctx, cancel := context.WithCancel(context.Background())

	// Signal OS pour arrêt propre
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Démarre les workers
	StartWorkerPool(ctx, db, jobQueue, 3)

	// Démarre l'API HTTP
	go StartServer(db, jobQueue)

	<-sigChan
	fmt.Println("\n[!] Signal de fermeture reçu")
	cancel()
	time.Sleep(2 * time.Second) // laisse le temps aux workers
	fmt.Println("Bye.")
}
```

---

### 📄 `db.go`

```go
package main

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func InitDB(path string) *sqlx.DB {
	db, err := sqlx.Connect("sqlite3", path)
	if err != nil {
		log.Fatalln("DB connect error:", err)
	}

	schema := `CREATE TABLE IF NOT EXISTS jobs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		status TEXT
	);`
	db.MustExec(schema)
	return db
}
```

---

### 📄 `job.go`

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type Job struct {
	ID     int    `db:"id"`
	Name   string `db:"name"`
	Status string `db:"status"`
}

func StartWorkerPool(ctx context.Context, db *sqlx.DB, jobQueue <-chan Job, numWorkers int) {
	for i := 0; i < numWorkers; i++ {
		go func(workerID int) {
			for {
				select {
				case <-ctx.Done():
					fmt.Printf("[Worker %d] Stop signal reçu\n", workerID)
					return
				case job := <-jobQueue:
					fmt.Printf("[Worker %d] Traitement du job %d: %s\n", workerID, job.ID, job.Name)
					time.Sleep(2 * time.Second) // simule traitement
					db.MustExec("UPDATE jobs SET status = ? WHERE id = ?", "done", job.ID)
				}
			}
		}(i + 1)
	}
}
```

---

### 📄 `handler.go`

```go
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func StartServer(db *sqlx.DB, jobQueue chan<- Job) {
	r := gin.Default()

	r.POST("/jobs", func(c *gin.Context) {
		var body struct {
			Name string `json:"name"`
		}
		if err := c.BindJSON(&body); err != nil || body.Name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Nom requis"})
			return
		}

		res := db.MustExec("INSERT INTO jobs (name, status) VALUES (?, ?)", body.Name, "queued")
		id, _ := res.LastInsertId()
		jobQueue <- Job{ID: int(id), Name: body.Name, Status: "queued"}

		c.JSON(http.StatusOK, gin.H{"message": "Job ajouté", "id": id})
	})

	r.GET("/jobs", func(c *gin.Context) {
		var jobs []Job
		db.Select(&jobs, "SELECT * FROM jobs")
		c.JSON(http.StatusOK, jobs)
	})

	r.Run(":8080")
}
```

---

### ▶️ Exécution

1. **Installe les dépendances :**

```bash
go mod init go-worker-api
go get github.com/gin-gonic/gin github.com/jmoiron/sqlx github.com/mattn/go-sqlite3
```

2. **Lance le projet :**

```bash
go run .
```

3. **Utilise l’API :**

```bash
curl -X POST http://localhost:8080/jobs -H "Content-Type: application/json" -d '{"name":"Traitement A"}'
curl http://localhost:8080/jobs
```

---

