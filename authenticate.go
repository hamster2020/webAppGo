package webAppGo

import (
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Seed the random number generator on start up
func init() {
	rand.Seed(time.Now().UnixNano())
}

// letterRunes is used to sample random runes for the GenRandID function
var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// EncryptPass will encrypt the password with the bcrypt algorithm
func EncryptPass(password string) (string, error) {
	pass := []byte(password)
	hashpw, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashpw), nil
}

// GenRandID generates a random ID n digits long
func GenRandID(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
