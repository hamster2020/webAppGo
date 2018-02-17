package models

import (
	"net/http"
	"time"
)

// Session is a type for storing User Sessions
type Session struct {
	SessionID string
	UserID    string
	Time      int
}

// Timeout variable is used to determine if the session should timeout (in seconds)
var Timeout = 300

// SaveSession saves a user session to the cache/db.sqlite3 db
func (db *DB) SaveSession(s *Session) error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS sessions (sessionid TEXT NOT NULL UNIQUE, userid TEXT NOT NULL , timestamp INTEGER NOT NULL, PRIMARY KEY(sessionid));")
	if err != nil {
		return err
	}
	_, err = db.Exec("INSERT INTO sessions (sessionid, userid, timestamp) VALUES ($1, $2, $3)", s.SessionID, s.UserID, s.Time)
	if err != nil {
		return err
	}
	return nil
}

// GetSessionFromSessionID retrieves a session from the db, given a session id
func (db *DB) GetSessionFromSessionID(sessionid string) *Session {
	var sid, uid string
	var time int
	rows, err := db.Query("SELECT * FROM sessions WHERE sessionid = '" + sessionid + "'")
	if err != nil {
		return &Session{}
	}
	for rows.Next() {
		rows.Scan(&sid, &uid, &time)
	}
	return &Session{
		SessionID: sid,
		UserID:    uid,
		Time:      time,
	}
}

// DeleteSession removes the session cookie
func (db *DB) DeleteSession(res http.ResponseWriter, sessionid string) error {
	//	clearCookie(res, "session")
	db.Exec("CREATE TABLE IF NOT EXISTS sessions (sessionid TEXT NOT NULL UNIQUE, userid TEXT NOT NULL, timestamp INTEGER NOT NULL, PRIMARY KEY(sessionid))")
	_, err := db.Exec("DELETE FROM sessions WHERE sessionid=$1", sessionid)
	if err != nil {
		return err
	}
	return nil
}

// IsSessionValid is used to check if a user/password combination exist in the db
func (db *DB) IsSessionValid(res http.ResponseWriter, sessionid string) (bool, string) {
	var sid, uid string
	var tm int
	rows, err := db.Query("SELECT * FROM sessions WHERE sessionid = '" + sessionid + "'")
	if err != nil {
		return false, ""
	}
	for rows.Next() {
		rows.Scan(&sid, &uid, &tm)
	}
	lastActivity := int(time.Now().Unix()) - tm
	if lastActivity > Timeout {
		db.DeleteSession(res, sessionid)
		return false, ""
	}
	if sid != "" {
		return true, sid
	}
	return false, ""
}
