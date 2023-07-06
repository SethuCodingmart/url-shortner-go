package utils

import (
	"crypto/rand"
	"math/big"
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
