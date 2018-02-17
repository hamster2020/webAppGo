package models

import (
	"log"
	"os/exec"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// EncryptPass will encrypt the password with the bcrypt algorithm
func EncryptPass(password string) string {
	pass := []byte(password)
	hashpw, _ := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
	return string(hashpw)
}

// UUID generates a universally unqiue ID
func UUID() string {
	out, err := exec.Command("uuidgen").Output()
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimSpace(string(out))
}
