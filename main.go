package main

import (
	"database/sql"
	"github.com/katatrina/my-simple-bank/api"
	"github.com/katatrina/my-simple-bank/applayer"
	db "github.com/katatrina/my-simple-bank/db/sqlc"
	"github.com/katatrina/my-simple-bank/token"
	"github.com/katatrina/my-simple-bank/util"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	conn, err := connectDB(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	store := db.NewStore(conn)

	jwt, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		log.Fatal("cannot create token maker: ", err)
	}

	app := applayer.New(store, jwt, config)

	server, err := api.NewHTTPServer(app, config, jwt)
	if err != nil {
		log.Fatal("cannot create server: ", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}

func connectDB(dbDriver, dbSource string) (*sql.DB, error) {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		return nil, err
	}

	if err = conn.Ping(); err != nil {
		return nil, err
	}

	return conn, nil
}
