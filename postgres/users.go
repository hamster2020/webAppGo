package sqlite

import (
	"github.com/hamster2020/webAppGo"
	"golang.org/x/crypto/bcrypt"
)

// SaveUser saves a user as a record to the users table in the db
func (db *DB) SaveUser(u *webAppGo.User) error {
	_, err := db.Exec(insertIntoUsersTable, u.UUID, u.Fname, u.Lname, u.Username, u.Email, u.Password)
	if err != nil {
		return err
	}
	return err
}

// DeleteUser deletes a user record from the users table in the db
func (db *DB) DeleteUser(u *webAppGo.User) error {
	_, err := db.Exec(deleteFromUsersTable, u.UUID)
	if err != nil {
		return err
	}
	return err
}

// UpdateUser saves a User struct to the db
func (db *DB) UpdateUser(u *webAppGo.User) error {
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
func (db *DB) GetUserFromUserID(userid string) (*webAppGo.User, error) {
	var uid, fn, ln, un, em, pass string
	rows, err := db.Query(selectUserFromTable, userid)
	if err != nil {
		return &webAppGo.User{}, err
	}
	for rows.Next() {
		rows.Scan(&uid, &fn, &ln, &un, &em, &pass)
	}
	return &webAppGo.User{
		Username: un,
		Fname:    fn,
		Lname:    ln,
		Email:    em,
		UUID:     uid,
		Password: pass,
	}, nil
}

// UserExists is used to check if a user/password combination exist in the db
func (db *DB) UserExists(u *webAppGo.User) (bool, string, error) {
	var password, userid string
	rows, err := db.Query(selectUsernamePasswordFromTable, u.Username)
	if err != nil {
		return false, "", err
	}
	for rows.Next() {
		rows.Scan(&userid, &password)
	}
	pwHashMatch := bcrypt.CompareHashAndPassword([]byte(password), []byte(u.Password))
	if userid != "" && pwHashMatch == nil {
		return true, userid, nil
	}
	return false, "", nil
}

// CheckUser checks if a given username is in the users table in the db
func (db *DB) CheckUser(username string) (bool, error) {
	var un string
	rows, err := db.Query(selectUsernameFromTable, username)
	if err != nil {
		return false, err
	}
	for rows.Next() {
		rows.Scan(&un)
	}
	if un == username {
		return true, nil
	}
	return false, nil
}
