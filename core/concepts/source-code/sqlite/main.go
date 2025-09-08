package main

import (
	"go-sqlite-api/db"
	"go-sqlite-api/routes"
)

func main() {
	db.InitDB("users.db")
	r := routes.SetupRouter()
	r.Run(":8080")
}
