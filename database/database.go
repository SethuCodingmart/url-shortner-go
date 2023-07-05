package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type DB struct {
	DB *sql.DB
}

var DBClient *sql.DB

const (
	host     = "localhost"
	port     = 5439
	user     = "root"
	password = "1471"
	dbname   = "urlshortner"
)

func ConnectDB() error {
	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname = %s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", connString)
	if err != nil {
		log.Printf("failed to connect to database: %v", err)
		return err
	} else {
		DBClient = db
	}
	return nil
}
