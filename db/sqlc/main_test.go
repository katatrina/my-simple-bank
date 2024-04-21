package db

import (
	"database/sql"
	"log"
	"testing"

	_ "github.com/lib/pq"
)

var (
	testQueries *Queries
	dbDriver    = "postgres"
	dbSource    = "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

func TestMain(m *testing.M) {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to database:", err)
	}

	testQueries = New(conn)

	m.Run()
}
