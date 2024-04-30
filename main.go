package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/katatrina/my-simple-bank/api"
	db "github.com/katatrina/my-simple-bank/db/sqlc"
	"github.com/katatrina/my-simple-bank/util"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig("config.env")
	if err != nil {
		log.Fatal("cannot load config", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	err = conn.Ping()
	if err != nil {
		log.Fatal("ping to db failed:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)
	err = server.Start(config.ServerAddress)
	fmt.Println(config)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
