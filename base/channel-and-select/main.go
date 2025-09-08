package main

import (
	"fmt"
	"time"
)

func main() {
	// Création de deux canaux
	ch1 := make(chan string)
	ch2 := make(chan string)

	// Goroutine qui envoie un message après 2 secondes
	go func() {
		time.Sleep(2 * time.Second)
		ch1 <- "Message de ch1"
		close(ch1)
	}()

	// Goroutine qui envoie un message après 1 seconde
	go func() {
		time.Sleep(1 * time.Second)
		ch2 <- "Message de ch2"
		// ch2 <- "Message de ch3"
		close(ch2)
	}()

	// On écoute avec select
	select {
	case msg1 := <-ch1:
		fmt.Println("Reçu :", msg1)
	case msg2 := <-ch2:
		fmt.Println("Reçu :", msg2)
	}
}
