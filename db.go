package webAppGo

import (
	"net/http"
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
