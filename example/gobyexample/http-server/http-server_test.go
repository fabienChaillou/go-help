package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
)

const (
	dns  = "127.0.0.1"
	port = "8090"
)

func TestMain(m *testing.M) {
	log.Println("Start Server of test!")
	existCode := m.Run()
	log.Println("End!")

	os.Exit(existCode)
}

func TestHttpBin(t *testing.T) {
	t.Run("Testing a get requst", func(t *testing.T) {
		t.Log("sending get request to http://localhost:8090")
		resp, err := http.Get(fmt.Sprintf("http://%v:%v/%v", dns, port, "hello"))
		if err != nil {
			t.Error(err.Error())
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("An unxpected status code has been returned: %d", resp.StatusCode)
		}
		defer resp.Body.Close()

		scanner := bufio.NewScanner(resp.Body)
		if scanner.Text() != "hello" {
			t.Error(err.Error())
		}
	})
}
