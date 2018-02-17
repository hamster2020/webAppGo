package models

import (
	"database/sql"
	"net/http"
	"regexp"
	"strings"
	//include the pq driver package here due to the VMC file structure
	_ "github.com/lib/pq"
)

// Datastore interface uses its interface-implementation feature to set up
// a clean implementation of dependency-injection
type Datastore interface {
	// pages.go db methods
	SavePage(*Page) error
	LoadPage(string) (*Page, error)
	PageExists(string) (bool, error)

	// users.go db methods
	SaveUser(*User) error
	DeleteUser(*User) error
	UpdateUser(*User) error
	GetUserFromUserID(string) *User
	UserExists(*User) (bool, string)
	CheckUser(username string) bool

	// sessions.go db methods
	SaveSession(*Session) error
	GetSessionFromSessionID(string) *Session
	DeleteSession(http.ResponseWriter, string) error
	IsSessionValid(http.ResponseWriter, string) (bool, string)

	// logins.go db methods
	SaveLogin(*Login) error
	CheckUserLoginAttempts(string) bool
	CheckIPLoginAttempts(string) bool

	// cookie.go db methods
	SetSession(*Session, http.ResponseWriter)
}

// DB is a shortened version of sql.DB that is designed to implement the
// Datastore interface (methods defined through out webAppGo/models/*.go files)
type DB struct {
	*sql.DB
}

// DataSourceDriver defines the driver type of the db
var DataSourceDriver = "postgres"

// DataSourceName defines the connection details of the db
var DataSourceName = "postgres://hamster2020:password@localhost/webappgo?sslmode=disable"

// sql statements:
// pages table
var createPagesTable = postgresVsSQLite(DataSourceDriver, "CREATE TABLE IF NOT EXISTS pages (title TEXT, body BYTEA, timestamp TEXT)")
var insertIntoPagesTable = postgresVsSQLite(DataSourceDriver, "INSERT INTO pages (title, body, timestamp) VALUES ($1, $2, $3)")
var selectPageFromTable = postgresVsSQLite(DataSourceDriver, "SELECT * FROM pages WHERE title = $1 ORDER BY timestamp DESC LIMIT 1")
var selectTitleBodyFromTable = postgresVsSQLite(DataSourceDriver, "SELECT title, body FROM pages WHERE title = $1 ORDER BY timestamp DESC LIMIT 1")

// users table
var createUsersTable = postgresVsSQLite(DataSourceDriver, "CREATE TABLE IF NOT EXISTS users (userid TEXT NOT NULL UNIQUE, firstname TEXT NOT NULL, lastname TEXT NOT NULL, username TEXT NOT NULL UNIQUE, email TEXT NOT NULL, password TEXT NOT NULL, PRIMARY KEY(userid))")
var insertIntoUsersTable = postgresVsSQLite(DataSourceDriver, "INSERT INTO users (userid, firstname, lastname, username, email, password) VALUES ($1, $2, $3, $4, $5, $6)")
var deleteFromUsersTable = postgresVsSQLite(DataSourceDriver, "DELETE FROM users WHERE userid=$1")
var selectUserFromTable = postgresVsSQLite(DataSourceDriver, "SELECT * FROM users WHERE userid = $1")
var selectUsernamePasswordFromTable = postgresVsSQLite(DataSourceDriver, "SELECT userid, password FROM users WHERE username = $1")
var selectUsernameFromTable = postgresVsSQLite(DataSourceDriver, "SELECT username FROM users WHERE username = $1")

// sessions table
var createSessionsTable = postgresVsSQLite(DataSourceDriver, "CREATE TABLE IF NOT EXISTS sessions (sessionid TEXT NOT NULL UNIQUE, userid TEXT NOT NULL , timestamp INTEGER NOT NULL, PRIMARY KEY(sessionid))")
var insertIntoSessionsTable = postgresVsSQLite(DataSourceDriver, "INSERT INTO sessions (sessionid, userid, timestamp) VALUES ($1, $2, $3)")
var selectSessionFromTable = postgresVsSQLite(DataSourceDriver, "SELECT * FROM sessions WHERE sessionid = $1")
var deleteSessionFromTable = postgresVsSQLite(DataSourceDriver, "DELETE FROM sessions WHERE sessionid=$1")

// logins table
var createLoginsTable = postgresVsSQLite(DataSourceDriver, "CREATE TABLE IF NOT EXISTS logins (ip TEXT NOT NULL, username TEXT NOT NULL, timestamp INTEGER NOT NULL, attempt TEXT NOT NULL)")
var insertIntoLoginsTable = postgresVsSQLite(DataSourceDriver, "INSERT INTO logins (ip, username, timestamp, attempt) VALUES ($1, $2, $3, $4)")
var selectRecentUsernamesFromLoginsTable = postgresVsSQLite(DataSourceDriver, "SELECT username FROM logins WHERE username = $1 AND timestamp > $2 AND attempt = '0'")
var selectRecentIPsFromLoginsTable = postgresVsSQLite(DataSourceDriver, "SELECT ip FROM logins WHERE ip = $1 AND timestamp > $2 AND attempt = '0'")

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

func postgresVsSQLite(driverType, stmt string) string {
	switch driverType {
	case "postgres":
		for i := 0; i < strings.Count(stmt, "?"); i++ {
			stmt = strings.Replace(stmt, "?", "$"+string(i), 1)
		}
		stmt = strings.Replace(stmt, "GLOB", "BYTEA", 0)
		return stmt
	case "sqlite3":
		var re = regexp.MustCompile(`([\$/][\d])`)
		s := re.ReplaceAllString(stmt, `?`)
		s = strings.Replace(s, "BYTEA", "GLOB", 0)
		return s
	}
	return ""
}
