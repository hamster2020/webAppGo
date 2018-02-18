package sqlite

import (
	"log"
	"strconv"
	"time"

	"github.com/hamster2020/webAppGo"
)

// SaveLogin saves the login info to the db
func (db *DB) SaveLogin(login *webAppGo.Login) error {
	_, err := db.Exec(createLoginsTable)
	if err != nil {
		return err
	}
	_, err = db.Exec(insertIntoLoginsTable, login.IP, login.UserName, login.Timestamp, login.Attempt)
	if err != nil {
		return err
	}
	return nil
}

// CheckUserLoginAttempts checks to see if the number of failed login attempts is
// greater than the alotted amount per unit time for a given username.
func (db *DB) CheckUserLoginAttempts(username string) bool {
	_, err := db.Exec(createLoginsTable)
	if err != nil {
		log.Fatal(err)
	}
	tm := int(time.Now().Unix()) - webAppGo.LoginAttemptTime
	rows, err := db.Query(selectRecentUsernamesFromLoginsTable, username, strconv.Itoa(tm))
	if err != nil {
		log.Fatal(err)
	}
	numFails := 0
	for rows.Next() {
		rows.Scan()
		numFails++
	}
	if numFails >= webAppGo.MaxUserAttempts {
		return false
	}
	return true
}

// CheckIPLoginAttempts checks to see if the number of failed login attempts is
// greater than the alotted amount per unit time from a given ip address.
func (db *DB) CheckIPLoginAttempts(ip string) bool {
	tm := int(time.Now().Unix()) - webAppGo.LoginAttemptTime
	rows, err := db.Query(selectRecentIPsFromLoginsTable, ip, strconv.Itoa(tm))
	if err != nil {
		log.Fatal(err)
	}
	numFails := 0
	for rows.Next() {
		rows.Scan()
		numFails++
	}
	if numFails >= webAppGo.MaxIPAttempts {
		return false
	}
	return true
}
