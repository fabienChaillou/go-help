package main

import (
	"fmt"
	"sync"
	"time"
)

// fetchData simulates fetching data from a service
func fetchData(id int, results chan<- string, wg *sync.WaitGroup) {
	// Defer WaitGroup's Done call to signal this goroutine is complete
	defer wg.Done()

	// Simulate varying response times
	duration := time.Duration(id*300) * time.Millisecond
	time.Sleep(duration)

	// Send result to the channel
	results <- fmt.Sprintf("Data from service %d", id)
	fmt.Printf("Service %d completed in %v\n", id, duration)
}

func main() {
	fmt.Println("Starting concurrent data fetching...")
	start := time.Now()

	// Number of services to fetch from
	numServices := 5

	// Create a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Create a channel for results
	results := make(chan string, numServices)

	// Launch goroutines for each service
	for i := 1; i <= numServices; i++ {
		wg.Add(1) // Increment WaitGroup counter
		go fetchData(i, results, &wg)
	}

	// Launch a goroutine to close the results channel once all fetchData goroutines complete
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect and print results as they arrive
	for result := range results {
		fmt.Println("Received:", result)
	}

	elapsed := time.Since(start)
	fmt.Printf("\nAll services completed in %v\n", elapsed)
	fmt.Println("If this were sequential, it would have taken approximately 4.5 seconds")
}
