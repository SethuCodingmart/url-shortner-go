package store

import (
	"database/sql"
	"log"
	"time"
	database "urlShortner/database"
	interfaceGo "urlShortner/interface"
	"urlShortner/utils"
)

func SaveURL(value interfaceGo.URLParameters) (bool, error) {

	q := `INSERT INTO urls (location, alias, userid) VALUES ($1,$2, $3)`
	_, err := database.DBClient.Exec(q, value.Location, value.Alias, value.Id)

	if err != nil {
		log.Printf("ERROR WHILE ADDING: %v", err)
		return false, err
	}

	return true, nil
}

func GetUrlsWithUserId(userId int) ([]database.Urls, error) {
	q := `SELECT id, location, alias, userid, isexpired, expiresat, createdat from urls where userId = $1`
	resultRow, err := database.DBClient.Query(q, userId)
	switch err {
	case sql.ErrNoRows:
		log.Printf("ERROR NO ROWS FOUND: %v", err)
		return nil, err
	case nil:
		result := []database.Urls{}
		for resultRow.Next() {
			var id, userId int
			var location, alias string
			var isexpired bool
			var expiresat, createdat *time.Time

			err := resultRow.Scan(&id, &location, &alias, &userId, &isexpired, &expiresat, &createdat)
			if err != nil {
				log.Println("ERROR WHILE SCAN ROW", err.Error())
				continue
			}
			newObject := database.Urls{Id: id, Location: location, Alias: alias, UserId: userId, IsExpired: isexpired, Expiresat: expiresat, Createdat: *createdat}
			result = append(result, newObject)
		}
		return result, nil
	default:
		log.Println("INTERNAL DATABASE ERROR := ", err.Error())
		return nil, err
	}
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

func CreateUser(_mail string, _password string) (bool, error) {
	q := `INSERT INTO users (gmail, password) VALUES ($1, $2)`
	_, err := database.DBClient.Exec(q, _mail, _password)

	if err != nil {
		log.Printf("ERROR WHILE ADDING USER: %v", err)
		return false, err
	}
	return true, nil
}

func GenerateAuthKey(_authkey string, _userId int) (bool, error) {
	q := `UPDATE users SET authkey = $1, updatedat = NOW() where id = $2`
	_, err := database.DBClient.Exec(q, _authkey, _userId)

	if err != nil {
		log.Printf("ERROR WHILE GENERATING AUTH KEY FOR USER: %v", err)
		return false, err
	}
	return true, nil
}

func GetUserWithAuthKey(authkey string) (*database.Users, error) {
	q := `SELECT id, gmail from users where authkey = $1`
	result := database.DBClient.QueryRow(q, authkey)
	user := &database.Users{}
	err := result.Scan(&user.Id, &user.Gmail)
	switch err {
	case sql.ErrNoRows:
		log.Printf("no rows are present")
		return nil, CustomError{"NO ROWS PRESENT"}
	case nil:
		return user, nil
	default:
		log.Print("ERROR OCCURS WHILE QUERYING", err)
		return nil, CustomError{"ERROR OCCURS WHILE QUERYING" + err.Error()}
	}
}

func GetUserWithId(id int) (*database.Users, error) {
	q := `SELECT id, gmail, authkey, createdat, updatedat from users where id = $1`

	result := database.DBClient.QueryRow(q, id)
	user := &database.Users{}
	err := result.Scan(&user.Id, &user.Gmail, &user.Authkey, &user.CreatedAt, &user.UpdatedAt)
	switch err {
	case sql.ErrNoRows:
		log.Printf("no rows are present")
		return nil, CustomError{"NO ROWS PRESENT"}
	case nil:
		return user, nil
	default:
		log.Print("ERROR OCCURS WHILE QUERYING", err)
		return nil, CustomError{"ERROR OCCURS WHILE QUERYING" + err.Error()}
	}
}

func GetUserWithGmailAndPassword(gmail string, password string) (*database.Users, error) {
	q := `SELECT id, gmail, authkey, createdat, updatedat, password from users where gmail = $1`
	result := database.DBClient.QueryRow(q, gmail)
	user := &database.Users{}
	var authKey sql.NullString
	err := result.Scan(&user.Id, &user.Gmail, &authKey, &user.CreatedAt, &user.UpdatedAt, &user.Password)
	if authKey.Valid {
		user.Authkey = &authKey.String
	} else {
		user.Authkey = nil
	}
	switch err {
	case sql.ErrNoRows:
		log.Printf("no rows are present")
		return nil, CustomError{"NO ROWS PRESENT"}
	case nil:
		checkPassword := utils.CheckPassword(password, user.Password)
		if checkPassword {
			return user, nil
		} else {
			return nil, CustomError{"PASSWORD WRONG"}
		}
	default:
		log.Print("ERROR OCCURS WHILE QUERYING", err)
		return nil, CustomError{"ERROR OCCURS WHILE QUERYING" + err.Error()}
	}
}

func HealthCheck() (bool, error) {
	q := `SELECT 1`
	result := database.DBClient.QueryRow(q)
	var value string
	err := result.Scan(&value)
	if err != nil {
		log.Printf("Error while Fetching %v", err)
		return false, err
	}
	if value == "1" {
		return true, nil
	}
	return false, CustomError{"VALUE NOT EQUAL TO 1"}
}

type CustomError struct {
	message string
}

func (e CustomError) Error() string {
	return e.message
}
