package store

import (
	"database/sql"
	"fmt"
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

func SaveOTP(value interfaceGo.SaveOTP) (bool, error) {

	q := `INSERT INTO otp (key, value, type) VALUES ($1,$2, $3)`
	_, err := database.DBClient.Exec(q, value.Key, value.Value, value.Type)

	if err != nil {
		log.Printf("ERROR WHILE ADDING: %v", err)
		return false, err
	}

	return true, nil
}

func CheckUserExsist(gmail string) (bool, error) {
	q := `SELECT id from users where gmail = $1`
	result := database.DBClient.QueryRow(q, gmail)
	var id string
	err := result.Scan(&id)
	switch err {
	case sql.ErrNoRows:
		return false, nil
	case nil:
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

func UpdateUser(_mail string, _password string) (bool, error) {
	q := `UPDATE users SET password = $1, updatedat = NOW() where gmail = $2`
	_, err := database.DBClient.Exec(q, _password, _mail)
	if err != nil {
		log.Printf("ERROR WHILE UPDATING USER: %v", err)
		return false, err
	}
	return true, nil
}

func CreateUser(_mail string, _password string, _username string, _name string, _phone string) (bool, error) {
	tx, err := database.DBClient.Begin()
	if err != nil {
		log.Fatal(err)
		return false, err
	}
	var lastInsertID int
	q := `INSERT INTO users (gmail, password, name, username, phone) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	ierr := tx.QueryRow(q, _mail, _password, _name, _username, _phone).Scan(&lastInsertID)

	if ierr != nil {
		tx.Rollback()
		log.Printf("ERROR WHILE ADDING USER: %v", ierr)
		return false, ierr
	}

	qc := `INSERT INTO credits (available, user_id) VALUES ($1, $2)`
	_, cerr := tx.Exec(qc, 100, lastInsertID)

	if cerr != nil {
		tx.Rollback()
		log.Printf("ERROR WHILE ADDING CREDITS: %v", cerr)
		return false, cerr
	}

	trc := `INSERT INTO transcations (user_id, tfor, type, value) VALUES ($1, $2, $3, $4)`
	_, trcerr := tx.Exec(trc, lastInsertID, "ALL", "CREDITED", 100)

	if trcerr != nil {
		tx.Rollback()
		log.Printf("ERROR WHILE ADDING CREDITS: %v", trcerr)
		return false, trcerr
	}

	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Transaction committed successfully!")

	return true, nil
}

func GetOTP(key string, _type string) (*database.OTP, error) {
	q := `SELECT value from otp where key = $1 AND type = $2 ORDER BY createdat DESC LIMIT 1`
	result := database.DBClient.QueryRow(q, key, _type)
	otp := &database.OTP{}

	err := result.Scan(&otp.Value)

	switch err {
	case sql.ErrNoRows:
		log.Printf("no rows are present")
		return nil, CustomError{"OTP not requested."}
	case nil:
		return otp, nil
	default:
		log.Print("ERROR OCCURS WHILE QUERYING", err)
		return nil, CustomError{"ERROR OCCURS WHILE QUERYING" + err.Error()}
	}
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
	q := `SELECT id from authkey where authkey = $1`
	result := database.DBClient.QueryRow(q, authkey)
	user := &database.Users{}
	err := result.Scan(&user.Id)
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
	q := `SELECT id, gmail, username, phone, name, createdat, updatedat, deletedat from users where id = $1 AND deletedat IS NULL`
	result := database.DBClient.QueryRow(q, id)
	fmt.Print(result)
	user := &database.Users{}
	err := result.Scan(&user.Id, &user.Gmail, &user.Username, &user.Phone, &user.Name, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)
	fmt.Print(err)

	switch err {
	case sql.ErrNoRows:
		log.Printf("no rows are present")
		return nil, CustomError{"User not found."}
	case nil:
		return user, nil
	default:
		log.Print("ERROR OCCURS WHILE QUERYING", err)
		return nil, CustomError{"ERROR OCCURS WHILE QUERYING" + err.Error()}
	}
}

func GetUserWithGmailAndPassword(gmail string, password string) (*database.Users, error) {
	q := `SELECT id, gmail, username, createdat, updatedat, password from users where gmail = $1`
	result := database.DBClient.QueryRow(q, gmail)
	user := &database.Users{}
	err := result.Scan(&user.Id, &user.Gmail, &user.Username, &user.CreatedAt, &user.UpdatedAt, &user.Password)
	switch err {
	case sql.ErrNoRows:
		log.Printf("no rows are present")
		return nil, CustomError{"No User Found."}
	case nil:
		checkPassword := utils.CheckPassword(password, user.Password)
		if checkPassword {
			return user, nil
		} else {
			return nil, CustomError{"Incorrect Password."}
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
