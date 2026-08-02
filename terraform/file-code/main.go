package main

import (
	"log"
	"net/http"

	"github.com/example/terraform-module/api"
	"github.com/example/terraform-module/terraform"
)

func main() {
	registry := terraform.NewRegistry()
	dispatcher := terraform.NewDispatcher(registry)
	handler := api.NewHandler(dispatcher)

	mux := http.NewServeMux()
	handler.Register(mux)

	log.Println("Terraform API listening on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
