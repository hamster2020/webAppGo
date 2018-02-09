package main

import (
	"net/http"

	"github.com/gorilla/securecookie"
)

// cookieHandler contains generated keys for encrypting and decrypting cookies
var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32),
)

// setSession creates a session for a user vie secure cookies
func setSession(u *User, res http.ResponseWriter) {
	value := map[string]string{
		"uuid": u.UUID,
	}
	encoded, err := cookieHandler.Encode("session", value)
	if err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(res, cookie)
	}
}

// getUserName extracts the username from the session cookie in the http response
func getUUID(req *http.Request) (uuid string) {
	cookie, err := req.Cookie("session")
	if err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			uuid = cookieValue["uuid"]
		}
	}
	return uuid
}

// clearSession removes the session cookie
func clearSession(res http.ResponseWriter, name string) {
	cookie := &http.Cookie{
		Name:   name,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(res, cookie)
}

// getMsg attempts to extract a secure cookie name from the request header,
// decodes it, and clears it in the response header
func getMsg(res http.ResponseWriter, req *http.Request, name string) (msg string) {
	if cookie, err := req.Cookie(name); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode(name, cookie.Value, &cookieValue); err == nil {
			msg = cookieValue[name]
			clearSession(res, name)
		}
	}
	return msg
}

// setMsg encodes a name-msg pair into a cookie and
// sends it in the response header.
func setMsg(res http.ResponseWriter, name string, msg string) {
	value := map[string]string{
		name: msg,
	}
	if encoded, err := cookieHandler.Encode(name, value); err == nil {
		cookie := &http.Cookie{
			Name:  name,
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(res, cookie)
	}
}
