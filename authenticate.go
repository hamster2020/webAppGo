package webAppGo

import (
	"os/exec"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// EncryptPass will encrypt the password with the bcrypt algorithm
func EncryptPass(password string) (string, error) {
	pass := []byte(password)
	hashpw, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashpw), nil
}

// UUID generates a universally unqiue ID
func UUID() (string, error) {
	out, err := exec.Command("uuidgen").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
