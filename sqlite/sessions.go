package sqlite

import (
	"net/http"
	"time"

	"github.com/hamster2020/webAppGo"
)

// SaveSession saves a user session to the db.sqlite3 db
func (db *DB) SaveSession(s *webAppGo.Session) error {
	_, err := db.Exec(createSessionsTable)
	if err != nil {
		return err
	}
	_, err = db.Exec(insertIntoSessionsTable, s.SessionID, s.UserID, s.Time)
	if err != nil {
		return err
	}
	return nil
}

// GetSessionFromSessionID retrieves a session from the db, given a session id
func (db *DB) GetSessionFromSessionID(sessionid string) *webAppGo.Session {
	var sid, uid string
	var time int
	rows, err := db.Query(selectSessionFromTable, sessionid)
	if err != nil {
		return &webAppGo.Session{}
	}
	for rows.Next() {
		rows.Scan(&sid, &uid, &time)
	}
	return &webAppGo.Session{
		SessionID: sid,
		UserID:    uid,
		Time:      time,
	}
}

// DeleteSession removes the session cookie
func (db *DB) DeleteSession(res http.ResponseWriter, sessionid string) error {
	//	clearCookie(res, "session")
	db.Exec(createSessionsTable)
	_, err := db.Exec(deleteSessionFromTable, sessionid)
	if err != nil {
		return err
	}
	return nil
}

// IsSessionValid is used to check if a user/password combination exist in the db
func (db *DB) IsSessionValid(res http.ResponseWriter, sessionid string) (bool, string) {
	var sid, uid string
	var tm int
	rows, err := db.Query(selectSessionFromTable, sessionid)
	if err != nil {
		return false, ""
	}
	for rows.Next() {
		rows.Scan(&sid, &uid, &tm)
	}
	lastActivity := int(time.Now().Unix()) - tm
	if lastActivity > webAppGo.Timeout {
		db.DeleteSession(res, sessionid)
		return false, ""
	}
	if sid != "" {
		return true, sid
	}
	return false, ""
}

// SetSession creates a session for a user vie secure cookies
func (db *DB) SetSession(s *webAppGo.Session, res http.ResponseWriter) {
	value := map[string]string{
		"uuid": s.SessionID,
	}
	db.SaveSession(s)
	encoded, err := webAppGo.CookieHandler.Encode("session", value)
	if err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(res, cookie)
	}
}
