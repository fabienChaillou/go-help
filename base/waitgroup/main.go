package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// trace les gorouting
// open http://localhost:6060/debug/pprof/goroutine?debug=2
func init() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
}

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

	// garde le serveur principal actif
	select {}
}
