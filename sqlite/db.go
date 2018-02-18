package sqlite

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// DataSourceDriver defines the driver type of the db
// choose between "postgres" and "sqlite3"
var DataSourceDriver = "sqlite3"

// DataSourceName defines the connection details of the db
// choose between "postgres://hamster2020:password@localhost/webappgo?sslmode=disable" and "cache/db.sqlite3"
var DataSourceName = "../../sqlite/db.sqlite3"

// sql statements:
// pages table
var createPagesTable = "CREATE TABLE IF NOT EXISTS pages (title TEXT, body GLOB, timestamp TEXT)"
var insertIntoPagesTable = "INSERT INTO pages (title, body, timestamp) VALUES (?, ?, ?)"
var selectPageFromTable = "SELECT * FROM pages WHERE title = ? ORDER BY timestamp DESC LIMIT 1"
var selectTitleBodyFromTable = "SELECT title, body FROM pages WHERE title = ? ORDER BY timestamp DESC LIMIT 1"

// users table
var createUsersTable = "CREATE TABLE IF NOT EXISTS users (userid TEXT NOT NULL UNIQUE, firstname TEXT NOT NULL, lastname TEXT NOT NULL, username TEXT NOT NULL UNIQUE, email TEXT NOT NULL, password TEXT NOT NULL, PRIMARY KEY(userid))"
var insertIntoUsersTable = "INSERT INTO users (userid, firstname, lastname, username, email, password) VALUES (?, ?, ?, ?, ?, ?)"
var deleteFromUsersTable = "DELETE FROM users WHERE userid=?"
var selectUserFromTable = "SELECT * FROM users WHERE userid = ?"
var selectUsernamePasswordFromTable = "SELECT userid, password FROM users WHERE username = ?"
var selectUsernameFromTable = "SELECT username FROM users WHERE username = ?"

// sessions table
var createSessionsTable = "CREATE TABLE IF NOT EXISTS sessions (sessionid TEXT NOT NULL UNIQUE, userid TEXT NOT NULL , timestamp INTEGER NOT NULL, PRIMARY KEY(sessionid))"
var insertIntoSessionsTable = "INSERT INTO sessions (sessionid, userid, timestamp) VALUES (?, ?, ?)"
var selectSessionFromTable = "SELECT * FROM sessions WHERE sessionid = ?"
var deleteSessionFromTable = "DELETE FROM sessions WHERE sessionid=?"

// logins table
var createLoginsTable = "CREATE TABLE IF NOT EXISTS logins (ip TEXT NOT NULL, username TEXT NOT NULL, timestamp INTEGER NOT NULL, attempt TEXT NOT NULL)"
var insertIntoLoginsTable = "INSERT INTO logins (ip, username, timestamp, attempt) VALUES (?, ?, ?, ?)"
var selectRecentUsernamesFromLoginsTable = "SELECT username FROM logins WHERE username = ? AND timestamp > ? AND attempt = '0'"
var selectRecentIPsFromLoginsTable = "SELECT ip FROM logins WHERE ip = ? AND timestamp > ? AND attempt = '0'"

// DB is a shortened version of sql.DB that is designed to implement the
// Datastore interface (methods defined through out webAppGo/models/*.go files)
type DB struct {
	*sql.DB
}

// NewDB initializes DB object, tests the connection, and returns it
func NewDB(dataSourceDriver, dataSourceName string) (*DB, error) {
	db, err := sql.Open(dataSourceDriver, dataSourceName)

	if err != nil {
		return nil, err
	}
	//  log.Println("sqlite/db.go: Successfully opened pool of db connections from SQLite3 db")

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	//  log.Println("sqlite/db.go: Established successful connection from SQLite db")

	return &DB{db}, nil
}
