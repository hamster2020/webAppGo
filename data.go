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

	_ "github.com/mattn/go-sqlite3"
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
	Time      string
}

// SaveCache is for saving webpages to a sqlite DB first, if fails, then attempts to load a cached file
func (p Page) SaveCache() error {
	var db, _ = sql.Open("sqlite3", "cache/db.sqlite3")
	defer db.Close()
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	if strings.Contains(p.Title, " ") {
		p.Title = strings.Replace(p.Title, " ", "_", -1)
	}
	f := "cache/" + p.Title + ".txt"
	db.Exec("create table if not exists pages (title text, body blob, timestamp text)")
	tx, _ := db.Begin()
	stmt, _ := tx.Prepare("insert into pages (title, body, timestamp) values (?, ?, ?)")
	_, err := stmt.Exec(p.Title, p.Body, timestamp)
	tx.Commit()
	ioutil.WriteFile(f, p.Body, 0600)
	return err
}

// load is for loading webpages from a file
func load(title string) (*Page, error) {
	var db, _ = sql.Open("sqlite3", "cache/db.sqlite3")
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
	var db, _ = sql.Open("sqlite3", "cache/db.sqlite3")
	defer db.Close()
	var name string
	var body []byte
	q, err := db.Query("select title, body from pages where title = '" + title + "' order by timestamp Desc limit 1")
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
	var db, _ = sql.Open("sqlite3", "cache/db.sqlite3")
	defer db.Close()
	db.Exec("create table if not exists users (uuid text not null unique, firstname text not null, lastname text not null, username text not null unique, email text not null, password text not null, primary key(uuid))")
	tx, _ := db.Begin()
	stmt, _ := tx.Prepare("insert into users (uuid, firstname, lastname, username, email, password) values (?, ?, ?, ?, ?, ?)")
	_, err := stmt.Exec(u.UUID, u.Fname, u.Lname, u.Username, u.Email, u.Password)
	tx.Commit()
	return err
}

// saveSession saves a user session to the cache/db.sqlite3 db
func saveSession(s *Session) error {
	var db, _ = sql.Open("sqlite3", "cache/db.sqlite3")
	defer db.Close()
	db.Exec("create table if not exists sessions (sessionid text not null unique, userid text not null, timestamp text not null, primary key(sessionid))")
	tx, _ := db.Begin()
	stmt, _ := tx.Prepare("insert into sessions (sessionid, userid, timestamp) values (?, ?, ?)")
	_, err := stmt.Exec(s.SessionID, s.UserID, s.Time)
	tx.Commit()
	return err
}

func pageExists(title string) (bool, error) {
	var db, _ = sql.Open("sqlite3", "cache/db.sqlite3")
	defer db.Close()
	var pt string
	var pb []byte
	q, err := db.Query("select title, body from pages where title = '" + title + "' order by timestamp Desc limit 1")
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
	var db, _ = sql.Open("sqlite3", "cache/db.sqlite3")
	defer db.Close()
	var sid, uid, time string
	q, err := db.Query("select * from sessions where sessionid = '" + sessionid + "'")
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
	var db, _ = sql.Open("sqlite3", "cache/db.sqlite3")
	defer db.Close()
	db.Exec("create table if not exists sessions (sessionid text not null unique, userid text not null, timestamp text not null, primary key(sessionid))")
	tx, _ := db.Begin()
	stmt, _ := tx.Prepare("delete from sessions where sessionid=?")
	_, err := stmt.Exec(sessionid)
	tx.Commit()
	return err
}

// used to check if a user/password combination exist in the db
func sessionExists(sessionid string) (bool, string) {
	var db, _ = sql.Open("sqlite3", "cache/db.sqlite3")
	defer db.Close()
	var sid, uid, time string
	q, err := db.Query("select * from sessions where sessionid = '" + sessionid + "'")
	if err != nil {
		return false, ""
	}
	for q.Next() {
		q.Scan(&sid, &uid, &time)
	}
	if sid != "" {
		return true, sid
	}
	return false, ""
}

func getUserFromUUID(uuid string) *User {
	var db, _ = sql.Open("sqlite3", "cache/db.sqlite3")
	defer db.Close()
	var uu, fn, ln, un, em, pass string
	q, err := db.Query("select * from users where uuid = '" + uuid + "'")
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
	var db, _ = sql.Open("sqlite3", "cache/db.sqlite3")
	defer db.Close()
	var password, userID string
	q, err := db.Query("select uuid, password from users where username = '" + u.Username + "'")
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
	var db, _ = sql.Open("sqlite3", "cache/db.sqlite3")
	defer db.Close()
	var un string
	q, err := db.Query("select username from users where username = '" + user + "'")
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
