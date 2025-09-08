package main

import (
	"fmt"
	"time"
)

func worker(name string, out chan<- string, done <-chan struct{}) {
	for i := 1; i <= 5; i++ {
		select {
		case <-done:
			fmt.Println(name, "annulé.")
			return
		case out <- fmt.Sprintf("%s envoie %d", name, i):
			time.Sleep(time.Duration(500+100*i) * time.Millisecond)
		}
	}
}

func main() {
	// ch1 := make(chan string)
	// ch2 := make(chan string)
	// done := make(chan struct{})
	num1 := 1 << 2
	fmt.Println(num1)

	num2 := 1 << 20

	fmt.Println(num2)

	// go worker("Worker1", ch1, done)
	// go worker("Worker2", ch2, done)

	// timeout := time.After(3 * time.Second)

	// for {
	// 	select {
	// 	case msg := <-ch1:
	// 		fmt.Println("Reçu de ch1:", msg)
	// 	case msg := <-ch2:
	// 		fmt.Println("Reçu de ch2:", msg)
	// 	case <-timeout:
	// 		fmt.Println("⏰ Timeout atteint ! On annule les workers.")
	// 		close(done)
	// 		time.Sleep(1 * time.Second) // Laisse le temps aux goroutines de terminer
	// 		return
	// 	}
	// }
}
