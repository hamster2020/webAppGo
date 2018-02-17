package models

import (
	"database/sql"
	"net/http"
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
