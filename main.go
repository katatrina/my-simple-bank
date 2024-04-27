package main

import (
	"database/sql"
	"log"

	"github.com/katatrina/my-simple-bank/api"
	db "github.com/katatrina/my-simple-bank/db/sqlc"
	_ "github.com/lib/pq"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	conn := openDB()
	store := db.NewStore(conn)
	server := api.NewServer(store)
	err := server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}

func openDB() *sql.DB {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	pingErr := conn.Ping()
	if pingErr != nil {
		log.Fatal("cannot connect to db:", pingErr)
	}

	return conn
}
