package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/lenimbugua/bot/util"
	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:password@localhost:5432/bot?sslmode=disable"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("Cannot load config", err)
	}
	testDB, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to database ", err)
	}

	testQueries = New(testDB)
	os.Exit(m.Run())
}
