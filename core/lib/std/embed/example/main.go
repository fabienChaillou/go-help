package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// addr is the bind address for the web server.
var file_default = "a.txt"

//go:embed assets/*
var content embed.FS

func GetFileContent(name string) (string, error) {
	data, err := content.ReadFile("assets/" + name)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func ListFiles() ([]fs.DirEntry, error) {
	return content.ReadDir("assets")
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM)
	defer stop()

	// Parse command line flags.
	file_default := flag.String("file", "", "datasource name")
	flag.Parse()
	if *file_default == "" {
		flag.Usage()
		return fmt.Errorf("required: -dsn")
	}

	data, err := GetFileContent(*file_default)
	if err != nil {
		return fmt.Errorf("cannot create table: %w", err)
	}

	Println(data)

	list, err := ListFiles()
	if err != nil {
		return fmt.Errorf("cannot create table: %w", err)
	}

	for _, f := range list {
		Println(f.Name())
	}

	// Wait for signal.
	<-ctx.Done()
	log.Print("Shutting down")

	return nil
}

func Println(data string) {
	fmt.Println(data)
}
