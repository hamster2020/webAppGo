package sqlite

import (
	"database/sql"

	_ "github.com/lib/pq"
)

// DataSourceDriver defines the driver type of the db
var DataSourceDriver = "postgres"

// DataSourceName defines the connection details of the db
var DataSourceName = "postgres://hamster2020:password@localhost/webappgo?sslmode=disable"

// sql statements:
// pages table
var createPagesTable = "CREATE TABLE IF NOT EXISTS pages (title TEXT, body BYTEA, timestamp TEXT)"
var insertIntoPagesTable = "INSERT INTO pages (title, body, timestamp) VALUES ($1, $2, $3)"
var selectPageFromTable = "SELECT * FROM pages WHERE title = ? ORDER BY timestamp DESC LIMIT 1"
var selectTitleBodyFromTable = "SELECT title, body FROM pages WHERE title = $1 ORDER BY timestamp DESC LIMIT 1"

// users table
var createUsersTable = "CREATE TABLE IF NOT EXISTS users (userid TEXT NOT NULL UNIQUE, firstname TEXT NOT NULL, lastname TEXT NOT NULL, username TEXT NOT NULL UNIQUE, email TEXT NOT NULL, password TEXT NOT NULL, PRIMARY KEY(userid))"
var insertIntoUsersTable = "INSERT INTO users (userid, firstname, lastname, username, email, password) VALUES ($1, $2, $3, $4, $5, $6)"
var deleteFromUsersTable = "DELETE FROM users WHERE userid=$1"
var selectUserFromTable = "SELECT * FROM users WHERE userid = $1"
var selectUsernamePasswordFromTable = "SELECT userid, password FROM users WHERE username = $1"
var selectUsernameFromTable = "SELECT username FROM users WHERE username = $1"

// sessions table
var createSessionsTable = "CREATE TABLE IF NOT EXISTS sessions (sessionid TEXT NOT NULL UNIQUE, userid TEXT NOT NULL , timestamp INTEGER NOT NULL, PRIMARY KEY(sessionid))"
var insertIntoSessionsTable = "INSERT INTO sessions (sessionid, userid, timestamp) VALUES ($1, $2, $3)"
var selectSessionFromTable = "SELECT * FROM sessions WHERE sessionid = $1"
var deleteSessionFromTable = "DELETE FROM sessions WHERE sessionid=$1"

// logins table
var createLoginsTable = "CREATE TABLE IF NOT EXISTS logins (ip TEXT NOT NULL, username TEXT NOT NULL, timestamp INTEGER NOT NULL, attempt TEXT NOT NULL)"
var insertIntoLoginsTable = "INSERT INTO logins (ip, username, timestamp, attempt) VALUES ($1, $2, $3, $4)"
var selectRecentUsernamesFromLoginsTable = "SELECT username FROM logins WHERE username = $1 AND timestamp > $2 AND attempt = '0'"
var selectRecentIPsFromLoginsTable = "SELECT ip FROM logins WHERE ip = $1 AND timestamp > $2 AND attempt = '0'"

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

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}
