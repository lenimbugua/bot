package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/lenimbugua/bot/api"
	db "github.com/lenimbugua/bot/db/sqlc"
	"github.com/lenimbugua/bot/util"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to database ", err)
	}

	store := db.NewSQLStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server", err)
	}
	var serverAddress string
	if config.Env == "PRODUCTION" {
		port := os.Getenv("PORT")
		serverAddress = ":" + port
	} else {
		serverAddress = config.ServerAddress
	}
	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("Cannot start server ", err)
	}
}
