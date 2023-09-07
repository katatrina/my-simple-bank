package db

import (
	"database/sql"
	"github.com/katatrina/my-simple-bank/util"
	"log"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
		return
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
		return
	}

	if pingErr := testDB.Ping(); pingErr != nil {
		log.Fatal("cannot ping db:", pingErr)
		return
	}

	testQueries = New(testDB)

	m.Run()
}
