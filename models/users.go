package models

import "golang.org/x/crypto/bcrypt"

// User contains the basic information on a given user for signing up
type User struct {
	UUID     string            `valid:"uuid"`
	Fname    string            `valid:"req,alpha"`
	Lname    string            `valid:"req,alpha"`
	Username string            `valid:"req,alph-num"`
	Email    string            `valid:"req,email"`
	Password string            `valid:"req"`
	Errors   map[string]string `valid:"-"`
}

// SaveUser saves a user as a record to the users table in the db
func (db *DB) SaveUser(u *User) error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS users (userid TEXT NOT NULL UNIQUE, firstname TEXT NOT NULL, lastname TEXT NOT NULL, username TEXT NOT NULL UNIQUE, email TEXT NOT NULL, password TEXT NOT NULL, PRIMARY KEY(userid));")
	if err != nil {
		return err
	}
	_, err = db.Exec("INSERT INTO users (userid, firstname, lastname, username, email, password) VALUES ($1, $2, $3, $4, $5, $6)", u.UUID, u.Fname, u.Lname, u.Username, u.Email, u.Password)
	if err != nil {
		return err
	}
	return err
}

// DeleteUser deletes a user record from the users table in the db
func (db *DB) DeleteUser(u *User) error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS users (userid TEXT NOT NULL UNIQUE, firstname TEXT NOT NULL, lastname TEXT NOT NULL, username TEXT NOT NULL UNIQUE, email TEXT NOT NULL, password TEXT NOT NULL, PRIMARY KEY(userid));")
	if err != nil {
		return err
	}
	_, err = db.Exec("DELETE FROM users WHERE userid=$1", u.UUID)
	if err != nil {
		return err
	}
	return err
}

// UpdateUser saves a User struct to the db
func (db *DB) UpdateUser(u *User) error {
	err := db.DeleteUser(u)
	if err != nil {
		return err
	}
	err = db.SaveUser(u)
	if err != nil {
		return err
	}
	return nil
}

// GetUserFromUserID retrieves a user record from the db, given the userid
func (db *DB) GetUserFromUserID(userid string) *User {
	var uid, fn, ln, un, em, pass string
	rows, err := db.Query("SELECT * FROM users WHERE userid = '" + userid + "'")
	if err != nil {
		return &User{}
	}
	for rows.Next() {
		rows.Scan(&uid, &fn, &ln, &un, &em, &pass)
	}
	return &User{
		Username: un,
		Fname:    fn,
		Lname:    ln,
		Email:    em,
		UUID:     uid,
		Password: pass,
	}
}

// UserExists is used to check if a user/password combination exist in the db
func (db *DB) UserExists(u *User) (bool, string) {
	var password, userid string
	rows, err := db.Query("SELECT userid, password FROM users WHERE username = '" + u.Username + "'")
	if err != nil {
		return false, ""
	}
	for rows.Next() {
		rows.Scan(&userid, &password)
	}
	pwHashMatch := bcrypt.CompareHashAndPassword([]byte(password), []byte(u.Password))
	if userid != "" && pwHashMatch == nil {
		return true, userid
	}
	return false, ""
}

// CheckUser checks if a given username is in the users table in the db
func (db *DB) CheckUser(username string) bool {
	var un string
	rows, err := db.Query("SELECT username FROM users WHERE username = '" + username + "'")
	if err != nil {
		return false
	}
	for rows.Next() {
		rows.Scan(&un)
	}
	if un == username {
		return true
	}
	return false
}
