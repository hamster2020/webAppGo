package main

import (
	"database/sql"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"

	//	_ "github.com/mattn/go-sqlite3"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// Page is our type for storing webpages in memory
type Page struct {
	Title string
	Body  []byte
}

// User contains the basic information on a given user for signing up
type User struct {
	UUID     string            `valid:"uuid"`
	Fname    string            `valid:"req,alpha"`
	Lname    string            `valid:"req,alpha"`
	Username string            `valid:"req,alph-num"`
	Email    string            `valid:"req,email"`
	Password string            `valid:"req"`
	Errors   map[string]string `valid:"-"`
}

// Session is a type for storing User Sessions
type Session struct {
	SessionID string
	UserID    string
	Time      int
}

// Timeout variable is used to determine if the session should timeout (in seconds)
var Timeout = 300

// LoginDetails is for tracking successful and failed login attempts of users and ip addresses
type LoginDetails struct {
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

// SaveCache is for saving webpages to a sqlite DB first, if fails, then attempts to load a cached file
func (p Page) SaveCache() error {
	db, err := sql.Open("postgres", "postgres://hamster2020:password@localhost/webappgo?sslmode=disable")
	defer db.Close()
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	if strings.Contains(p.Title, " ") {
		p.Title = strings.Replace(p.Title, " ", "_", -1)
	}
	f := "cache/" + p.Title + ".txt"
	db.Exec("CREATE TABLE IF NOT EXISTS webappschema.pages (title TEXT, body BYTEA, timestamp TEXT);")
	tx, _ := db.Begin()
	_, err = tx.Exec("INSERT INTO webappschema.pages (title, body, timestamp) VALUES ($1, $2, $3)", p.Title, p.Body, timestamp)
	tx.Commit()
	ioutil.WriteFile(f, p.Body, 0600)
	return err
}

// load is for loading webpages from a file
func load(title string) (*Page, error) {
	db, err := sql.Open("postgres", "postgres://hamster2020:password@localhost/webappgo?sslmode=disable")
	defer db.Close()
	f := "cache/" + title + ".txt"
	body, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

// loadSource is for loading webpages from a database
func loadSource(title string) (*Page, error) {
	db, err := sql.Open("postgres", "postgres://hamster2020:password@localhost/webappgo?sslmode=disable")
	defer db.Close()
	var name string
	var body []byte
	q, err := db.Query("SELECT * FROM webappschema.pages WHERE title = '" + title + "' ORDER BY timestamp DESC LIMIT 1")
	if err != nil {
		return nil, err
	}
	for q.Next() {
		q.Scan(&name, &body)
	}
	return &Page{Title: name, Body: body}, nil
}

// saveUserData saves a User struct to the cache/db.sqlite3 db
func saveUserData(u *User) error {
	db, err := sql.Open("postgres", "postgres://hamster2020:password@localhost/webappgo?sslmode=disable")
	defer db.Close()
	db.Exec("CREATE TABLE IF NOT EXISTS webappschema.users (userid TEXT NOT NULL UNIQUE, firstname TEXT NOT NULL, lastname TEXT NOT NULL, username TEXT NOT NULL UNIQUE, email TEXT NOT NULL, password TEXT NOT NULL, PRIMARY KEY(userid));")
	tx, _ := db.Begin()
	_, err = tx.Exec("INSERT INTO webappschema.users (userid, firstname, lastname, username, email, password) VALUES ($1, $2, $3, $4, $5, $6)", u.UUID, u.Fname, u.Lname, u.Username, u.Email, u.Password)
	tx.Commit()
	return err
}

// deletes a user record from the users db
func deleteUserData(u *User) error {
	db, err := sql.Open("postgres", "postgres://hamster2020:password@localhost/webappgo?sslmode=disable")
	defer db.Close()
	db.Exec("CREATE TABLE IF NOT EXISTS webappschema.users (userid TEXT NOT NULL UNIQUE, firstname TEXT NOT NULL, lastname TEXT NOT NULL, username TEXT NOT NULL UNIQUE, email TEXT NOT NULL, password TEXT NOT NULL, PRIMARY KEY(userid));")
	tx, _ := db.Begin()
	_, err = tx.Exec("DELETE FROM webappschema.users WHERE userid=$1", u.UUID)
	tx.Commit()
	return err
}

// updateUserData saves a User struct to the cache/db.sqlite3 db
func updateUserData(u *User) error {
	err := deleteUserData(u)
	if err != nil {
		return err
	}
	err = saveUserData(u)
	if err != nil {
		return err
	}
	return nil
}

// saveSession saves a user session to the cache/db.sqlite3 db
func saveSession(s *Session) error {
	db, err := sql.Open("postgres", "postgres://hamster2020:password@localhost/webappgo?sslmode=disable")
	defer db.Close()
	db.Exec("CREATE TABLE IF NOT EXISTS webappschema.sessions (sessionid TEXT NOT NULL UNIQUE, userid TEXT NOT NULL , timestamp INTEGER NOT NULL, PRIMARY KEY(sessionid));")
	tx, _ := db.Begin()
	_, err = tx.Exec("INSERT INTO webappschema.sessions (sessionid, userid, timestamp) VALUES ($1, $2, $3)", s.SessionID, s.UserID, s.Time)
	tx.Commit()
	return err
}

func pageExists(title string) (bool, error) {
	db, err := sql.Open("postgres", "postgres://hamster2020:password@localhost/webappgo?sslmode=disable")
	defer db.Close()
	var pt string
	var pb []byte
	q, err := db.Query("SELECT title, body FROM webappschema.pages WHERE title = '" + title + "' ORDER BY timestamp DESC LIMIT 1")
	if err != nil {
		return false, err
	}
	for q.Next() {
		q.Scan(&pt, &pb)
	}
	if pt != "" && pb != nil {
		return true, nil
	}
	return false, nil
}

func getSessionFromSessionID(sessionid string) *Session {
	db, err := sql.Open("postgres", "postgres://hamster2020:password@localhost/webappgo?sslmode=disable")
	defer db.Close()
	var sid, uid string
	var time int
	q, err := db.Query("SELECT * FROM webappschema.sessions WHERE sessionid = '" + sessionid + "'")
	if err != nil {
		return &Session{}
	}
	for q.Next() {
		q.Scan(&sid, &uid, &time)
	}
	return &Session{
		SessionID: sid,
		UserID:    uid,
		Time:      time,
	}
}

// clearSession removes the session cookie
func clearSession(res http.ResponseWriter, sessionid string) error {
	clearCookie(res, "session")
	db, err := sql.Open("postgres", "postgres://hamster2020:password@localhost/webappgo?sslmode=disable")
	defer db.Close()
	db.Exec("CREATE TABLE IF NOT EXISTS webappschema.sessions (sessionid TEXT NOT NULL UNIQUE, userid TEXT NOT NULL, timestamp INTEGER NOT NULL, PRIMARY KEY(sessionid))")
	tx, _ := db.Begin()
	_, err = tx.Exec("DELETE FROM webappschema.sessions WHERE sessionid=$1", sessionid)
	tx.Commit()
	return err
}

// used to check if a user/password combination exist in the db
func sessionIsValid(res http.ResponseWriter, sessionid string) (bool, string) {
	db, err := sql.Open("postgres", "postgres://hamster2020:password@localhost/webappgo?sslmode=disable")
	defer db.Close()
	var sid, uid string
	var tm int
	q, err := db.Query("SELECT * FROM webappschema.sessions WHERE sessionid = '" + sessionid + "'")
	if err != nil {
		return false, ""
	}
	for q.Next() {
		q.Scan(&sid, &uid, &tm)
	}
	lastActivity := int(time.Now().Unix()) - tm
	if lastActivity > Timeout {
		clearSession(res, sessionid)
		return false, ""
	}
	if sid != "" {
		return true, sid
	}
	return false, ""
}

func getUserFromUUID(uuid string) *User {
	db, err := sql.Open("postgres", "postgres://hamster2020:password@localhost/webappgo?sslmode=disable")
	defer db.Close()
	var uu, fn, ln, un, em, pass string
	q, err := db.Query("SELECT * FROM webappschema.users WHERE userid = '" + uuid + "'")
	if err != nil {
		return &User{}
	}
	for q.Next() {
		q.Scan(&uu, &fn, &ln, &un, &em, &pass)
	}
	return &User{
		Username: un,
		Fname:    fn,
		Lname:    ln,
		Email:    em,
		UUID:     uu,
		Password: pass,
	}
}

// used to check if a user/password combination exist in the db
func userExists(u *User) (bool, string) {
	db, err := sql.Open("postgres", "postgres://hamster2020:password@localhost/webappgo?sslmode=disable")
	defer db.Close()
	var password, userID string
	q, err := db.Query("SELECT userid, password FROM webappschema.users WHERE username = '" + u.Username + "'")
	if err != nil {
		return false, ""
	}
	for q.Next() {
		q.Scan(&userID, &password)
	}
	pw := bcrypt.CompareHashAndPassword([]byte(password), []byte(u.Password))
	if userID != "" && pw == nil {
		return true, userID
	}
	return false, ""
}

func checkUser(user string) bool {
	db, err := sql.Open("postgres", "postgres://hamster2020:password@localhost/webappgo?sslmode=disable")
	defer db.Close()
	var un string
	q, err := db.Query("SELECT username FROM webappschema.users WHERE username = '" + user + "'")
	if err != nil {
		return false
	}
	for q.Next() {
		q.Scan(&un)
	}
	if un == user {
		return true
	}
	return false
}

// encryptPass will encrypt the password with the bcrypt algorithm
func encryptPass(password string) string {
	pass := []byte(password)
	hashpw, _ := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
	return string(hashpw)
}

// uuid generates a universally unqiue ID
func uuid() string {
	out, err := exec.Command("uuidgen").Output()
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimSpace(string(out))
}

func storeUserLogin(login *LoginDetails) error {
	db, err := sql.Open("postgres", "postgres://hamster2020:password@localhost/webappgo?sslmode=disable")
	defer db.Close()
	db.Exec("CREATE TABLE IF NOT EXISTS webappschema.logins (ip TEXT NOT NULL, username TEXT NOT NULL, timestamp INTEGER NOT NULL, attempt TEXT NOT NULL)")
	tx, _ := db.Begin()
	_, err = tx.Exec("INSERT INTO webappschema.logins (ip, username, timestamp, attempt) VALUES ($1, $2, $3, $4)", login.IP, login.UserName, login.Timestamp, login.Attempt)
	tx.Commit()
	return err
}

func checkUserLoginAttempts(username string) bool {
	db, err := sql.Open("postgres", "postgres://hamster2020:password@localhost/webappgo?sslmode=disable")
	defer db.Close()
	db.Exec("CREATE TABLE IF NOT EXISTS webappschema.logins (ip TEXT NOT NULL, username TEXT NOT NULL, timestamp INTEGER NOT NULL, attempt TEXT NOT NULL)")
	tm := int(time.Now().Unix()) - LoginAttemptTime
	q, err := db.Query("SELECT username FROM webappschema.logins WHERE username = '" + username + "' AND timestamp > '" + strconv.Itoa(tm) + "' AND attempt = '0'")
	if err != nil {
		log.Fatal(err)
	}
	numFails := 0
	for q.Next() {
		q.Scan()
		numFails++
	}
	if numFails >= MaxUserAttempts {
		return false
	}
	return true
}

func checkIPLoginAttempts(ip string) bool {
	db, err := sql.Open("postgres", "postgres://hamster2020:password@localhost/webappgo?sslmode=disable")
	defer db.Close()
	tm := int(time.Now().Unix()) - LoginAttemptTime
	q, err := db.Query("SELECT ip FROM webappschema.logins WHERE ip = '" + ip + "' AND timestamp > '" + strconv.Itoa(tm) + "' AND attempt = '0'")
	if err != nil {
		log.Fatal(err)
	}
	numFails := 0
	for q.Next() {
		q.Scan()
		numFails++
	}
	if numFails >= MaxIPAttempts {
		return false
	}
	return true
}
