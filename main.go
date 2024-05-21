package main

import (
	"database/sql"
	"log"

	"github.com/katatrina/my-simple-bank/api"
	db "github.com/katatrina/my-simple-bank/db/sqlc"
	"github.com/katatrina/my-simple-bank/util"
	_ "github.com/lib/pq"
)

//	@title			Simple Bank API

//	@contact.name Châu Vĩnh Phước
//	@contact.email	cvphuoc2014@gmail.com

//	@host		localhost:8080

//	@securityDefinitions.apiKey		accessToken
//	@in								header
//	@name							Authorization
//	@description			    	JWT Authorization header using the Bearer scheme.

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
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server")
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
