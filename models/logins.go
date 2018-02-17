package models

import (
	"log"
	"strconv"
	"time"
)

// Login is for tracking successful and failed login attempts of users and ip addresses
type Login struct {
	IP        string
	UserName  string
	Timestamp int
	Attempt   bool
}

// MaxUserAttempts is the maximum number of failed login attempts to be made
// on a specific user over a given period of time
var MaxUserAttempts = 3

// MaxIPAttempts is the maximum number of failed login attempts to be made
// from a specific ip address over a given period of time
var MaxIPAttempts = 6

// LoginAttemptTime is the period of time in which the use is allowed to make
// incorrect logins within this time frame
var LoginAttemptTime = 10 * 60

// SaveLogin saves the login info to the db
func (db *DB) SaveLogin(login *Login) error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS logins (ip TEXT NOT NULL, username TEXT NOT NULL, timestamp INTEGER NOT NULL, attempt TEXT NOT NULL)")
	if err != nil {
		return err
	}
	_, err = db.Exec("INSERT INTO logins (ip, username, timestamp, attempt) VALUES ($1, $2, $3, $4)", login.IP, login.UserName, login.Timestamp, login.Attempt)
	if err != nil {
		return err
	}
	return nil
}

// CheckUserLoginAttempts checks to see if the number of failed login attempts is
// greater than the alotted amount per unit time for a given username.
func (db *DB) CheckUserLoginAttempts(username string) bool {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS logins (ip TEXT NOT NULL, username TEXT NOT NULL, timestamp INTEGER NOT NULL, attempt TEXT NOT NULL)")
	if err != nil {
		log.Fatal(err)
	}
	tm := int(time.Now().Unix()) - LoginAttemptTime
	rows, err := db.Query("SELECT username FROM logins WHERE username = '" + username + "' AND timestamp > '" + strconv.Itoa(tm) + "' AND attempt = '0'")
	if err != nil {
		log.Fatal(err)
	}
	numFails := 0
	for rows.Next() {
		rows.Scan()
		numFails++
	}
	if numFails >= MaxUserAttempts {
		return false
	}
	return true
}

// CheckIPLoginAttempts checks to see if the number of failed login attempts is
// greater than the alotted amount per unit time from a given ip address.
func (db *DB) CheckIPLoginAttempts(ip string) bool {
	tm := int(time.Now().Unix()) - LoginAttemptTime
	rows, err := db.Query("SELECT ip FROM logins WHERE ip = '" + ip + "' AND timestamp > '" + strconv.Itoa(tm) + "' AND attempt = '0'")
	if err != nil {
		log.Fatal(err)
	}
	numFails := 0
	for rows.Next() {
		rows.Scan()
		numFails++
	}
	if numFails >= MaxIPAttempts {
		return false
	}
	return true
}
