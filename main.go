package main

import (
	"database/sql"
	"github.com/katatrina/my-simple-bank/api"
	db "github.com/katatrina/my-simple-bank/db/sqlc"
	_ "github.com/lib/pq"
	"log"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgres://postgres:secret@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "127.0.0.1:8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	if pingErr := conn.Ping(); pingErr != nil {
		log.Fatal("cannot ping db:", pingErr)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
