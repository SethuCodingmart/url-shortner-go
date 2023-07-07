package main

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "1471"
	dbname   = "urlShortner"
)

func main() {
	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname = %s sslmode=disable", host, port, user, password, dbname)

	migrationsDir := "file:///D:/url-shortner-go/migration/migrations/001_create_users_table.up.sql"
	fmt.Println(migrationsDir)

	m, err := migrate.New(migrationsDir, connString)
	if err != nil {
		log.Fatal("ERROR", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	fmt.Println("Migration successful")
}
