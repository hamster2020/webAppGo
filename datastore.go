package webAppGo

import (
	"net/http"
)

// Datastore interface uses its interface-implementation feature to set up
// a clean implementation of dependency-injection for the web methods
type Datastore interface {
	// pages.go db methods
	SavePage(*Page) error
	LoadPage(string) (*Page, error)
	PageExists(string) (bool, error)

	// users.go db methods
	SaveUser(*User) error
	DeleteUser(*User) error
	UpdateUser(*User) error
	GetUserFromUserID(string) (*User, error)
	UserExists(*User) (bool, string, error)
	CheckUser(string) (bool, error)

	// sessions.go db methods
	SaveSession(*Session) error
	GetSessionFromSessionID(string) (*Session, error)
	DeleteSession(http.ResponseWriter, string) error
	IsSessionValid(http.ResponseWriter, string) (bool, string, error)

	// logins.go db methods
	SaveLogin(*Login) error
	CheckUserLoginAttempts(string) (bool, error)
	CheckIPLoginAttempts(string) (bool, error)

	// cookie.go db methods
	SetSession(*Session, http.ResponseWriter) error
}
