package main

import (
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	log.Println("Start TestMain!")
	exitVal := m.Run()
	log.Println("End!")

	os.Exit(exitVal)
}
