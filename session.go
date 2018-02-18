package webAppGo

import (
	"net/http"

	"github.com/gorilla/securecookie"
)

// Session is a type for storing User Sessions
type Session struct {
	SessionID string
	UserID    string
	Time      int
}

// Timeout variable is used to determine if the session should timeout (in seconds)
var Timeout = 300

// CookieHandler contains generated keys for encrypting and decrypting cookies
var CookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32),
)

// GetSessionIDFromCookie extracts the username from the session cookie in the http response
func GetSessionIDFromCookie(req *http.Request) (uuid string) {
	cookie, err := req.Cookie("session")
	if err == nil {
		cookieValue := make(map[string]string)
		if err = CookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			uuid = cookieValue["uuid"]
		}
	}
	return uuid
}

// ClearCookie removes the given cookie from the client's browser
func ClearCookie(res http.ResponseWriter, name string) {
	cookie := &http.Cookie{
		Name:   name,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(res, cookie)
}

// GetMsg attempts to extract a secure cookie name from the request header,
// decodes it, and clears it in the response header
func GetMsg(res http.ResponseWriter, req *http.Request, name string) (msg string) {
	if cookie, err := req.Cookie(name); err == nil {
		cookieValue := make(map[string]string)
		if err = CookieHandler.Decode(name, cookie.Value, &cookieValue); err == nil {
			msg = cookieValue[name]
			ClearCookie(res, name)
		}
	}
	return msg
}

// SetMsg encodes a name-msg pair into a cookie and
// sends it in the response header.
func SetMsg(res http.ResponseWriter, name string, msg string) {
	value := map[string]string{
		name: msg,
	}
	if encoded, err := CookieHandler.Encode(name, value); err == nil {
		cookie := &http.Cookie{
			Name:  name,
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(res, cookie)
	}
}
