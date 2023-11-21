package utils

import (
	"crypto/rand"
	"math/big"

	"golang.org/x/crypto/bcrypt"
)

func GenerateRandomString(length int) (string, error) {
	charset := "abcdefghijklmnopqrstuvwxyz"
	chars := make([]byte, length)

	for i := 0; i < length; i++ {
		randomInt, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		chars[i] = charset[randomInt.Int64()]
	}

	return string(chars), nil
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func GenerateRandomAlphaNumericString(length int) (string, error) {
	charset := "abcdefghijklmnopqrstuvwxyz12345678"
	chars := make([]byte, length)

	for i := 0; i < length; i++ {
		randomInt, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		chars[i] = charset[randomInt.Int64()]
	}

	return string(chars), nil
}