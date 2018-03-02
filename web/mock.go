package web

import (
	"net/http"
	"time"

	"github.com/hamster2020/webAppGo"
)

type mockDB struct{}
type mockCache struct {
	Path string
}

// Datastore pages methods
func (mdb *mockDB) AllPages() ([]*webAppGo.Page, error) {
	return nil, nil
}

func (mdb *mockDB) SavePage(p *webAppGo.Page) error {
	return nil
}

func (mdb *mockDB) LoadPage(p string) (*webAppGo.Page, error) {
	return &webAppGo.Page{Title: p, Body: []byte("testing testing 123")}, nil
}

func (mdb *mockDB) PageExists(p string) (bool, error) {
	return true, nil
}

func (mc *mockCache) SaveToCache(*webAppGo.Page) error {
	return nil
}

func (mc *mockCache) LoadPageFromCache(title string) (*webAppGo.Page, error) {
	return &webAppGo.Page{Title: title, Body: []byte("testing testing 123")}, nil
}

// Datastore users methods
func (mdb *mockDB) SaveUser(u *webAppGo.User) error {
	return nil
}

func (mdb *mockDB) DeleteUser(u *webAppGo.User) error {
	return nil
}

func (mdb *mockDB) UpdateUser(u *webAppGo.User) error {
	err := mdb.DeleteUser(u)
	if err != nil {
		return err
	}
	err = mdb.SaveUser(u)
	if err != nil {
		return err
	}
	return nil
}

func (mdb *mockDB) GetUserFromUserID(userid string) (*webAppGo.User, error) {
	pw, _ := webAppGo.EncryptPass("test")
	return &webAppGo.User{UserID: userid, Username: "test", Fname: "test", Lname: "test", Email: "test@email.com", Password: pw}, nil
}

func (mdb *mockDB) UserExists(u *webAppGo.User) (bool, string, error) {
	return true, u.UserID, nil
}

func (mdb *mockDB) CheckUser(username string) (bool, error) {
	if username != "bobjoe13" {
		return true, nil
	}
	return false, nil
}

// Datastore sessions methods
func (mdb *mockDB) SaveSession(s *webAppGo.Session) error {
	return nil
}

func (mdb *mockDB) GetSessionFromSessionID(sid string) *webAppGo.Session {
	return &webAppGo.Session{SessionID: sid, UserID: "test", Time: int(time.Now().Unix())}
}

func (mdb *mockDB) DeleteSession(res http.ResponseWriter, sid string) error {
	return nil
}

func (mdb *mockDB) IsSessionValid(res http.ResponseWriter, sid string) (bool, string) {
	return true, ""
}

// Datastore logins methods
func (mdb *mockDB) SaveLogin(login *webAppGo.Login) error {
	return nil
}

func (mdb *mockDB) CheckUserLoginAttempts(userid string) (bool, error) {
	return true, nil
}

func (mdb *mockDB) CheckIPLoginAttempts(ip string) (bool, error) {
	return true, nil
}

func (mdb *mockDB) SetSession(s *webAppGo.Session, res http.ResponseWriter) error {
	return nil
}
