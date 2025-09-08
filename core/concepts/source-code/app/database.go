package main

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func InitDB() *sqlx.DB {
	db, err := sqlx.Connect("postgres", "user=postgres password=postgres dbname=testdb sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	return db
}
