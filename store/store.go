package store

import (
	"log"
	database "urlShortner/database"
	interfaceGo "urlShortner/interface"
)


func SaveURL(value interfaceGo.URLParameters) (bool, error) {

	q := `INSERT INTO urls (location, alias) VALUES ($1,$2)`
	_, err := database.DBClient.Exec(q, value.Location, value.Alias)

	if err != nil {
		log.Printf("ERROR WHILE ADDING: %v", err)
		return false, err
	}

	return true, nil
}

func RedirectURL(_alias string) string {
	q := `SELECT id, alias, location from urls where alias = $1`
	result := database.DBClient.QueryRow(q, _alias)
	var id, alias, location string
	err := result.Scan(&id, &alias, &location)
	if err != nil {
		log.Printf("Error while Fetching %v", err)
	}
	if location != "" {
		return location
	}
	return "/page/404"
}
