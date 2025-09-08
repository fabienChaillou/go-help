package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// addr is the bind address for the web server.
const addr = ":8090"

func main() {
	log.Println("Start server!")

	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM)
	defer stop()

	// Run web server.
	fmt.Printf("listening on %s\n", addr)
	go http.ListenAndServe(addr,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println("ping hello!")

			// Print total page views.
			fmt.Fprintf(w, "hello\n")
		}),
	)

	// Wait for signal.
	<-ctx.Done()
	log.Print("Received signal, shutting down")
	return nil
}
