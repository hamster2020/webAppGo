package web

import (
	"net/http"
	"time"

	"github.com/hamster2020/webAppGo"
)

type mockDB struct{}
type mockCache struct {
	//	SaveToCacheFunc   func() error
	//	LoadPageFromCache func(string) (*webAppGo.Page, error)
}

// Datastore pages methods
func (mdb *mockDB) SavePage(p *webAppGo.Page) error {
	return nil
}

func (mdb *mockDB) LoadPage(p string) (*webAppGo.Page, error) {
	return &webAppGo.Page{Title: p, Body: []byte("testing testing 123")}, nil
}

func (mdb *mockDB) PageExists(p string) (bool, error) {
	return true, nil
}

func (mc *mockCache) SaveToCache() error {
	return nil
	//return mc.SaveToCacheFunc()
}

func (mc *mockCache) LoadPageFromCache(title string) (*webAppGo.Page, error) {
	return &webAppGo.Page{Title: title, Body: []byte("testing testing 123")}, nil
	//return mc.LoadPageFromCache(title)
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
	return &webAppGo.User{UUID: userid, Username: "test", Fname: "test", Lname: "test", Email: "test@email.com", Password: pw}, nil
}

func (mdb *mockDB) UserExists(u *webAppGo.User) (bool, string, error) {
	return true, u.UUID, nil
}

func (mdb *mockDB) CheckUser(username string) (bool, error) {
	return true, nil
}

// Datastore sessions methods
func (mdb *mockDB) SaveSession(s *webAppGo.Session) error {
	return nil
}

func (mdb *mockDB) GetSessionFromSessionID(sid string) (*webAppGo.Session, error) {
	return &webAppGo.Session{SessionID: sid, UserID: "test", Time: int(time.Now().Unix())}, nil
}

func (mdb *mockDB) DeleteSession(res http.ResponseWriter, sid string) error {
	return nil
}

func (mdb *mockDB) IsSessionValid(res http.ResponseWriter, sid string) (bool, string, error) {
	return true, "", nil
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

/*
// main.go tests
func TestView(t *testing.T) {
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/view/test", nil)
	env := Env{DB: &mockDB{}}

	http.HandlerFunc(env.CheckUUID(CheckPath(env.View))).ServeHTTP(res, req)

	if res.Result().StatusCode != 200 {
		t.Errorf("obtained invald status, received response: \n%v", *res)
	}
}

func TestEdit(t *testing.T) {
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/edit/test", nil)
	env := Env{DB: &mockDB{}}

	http.HandlerFunc(env.CheckUUID(CheckPath(env.Edit))).ServeHTTP(res, req)

	if res.Result().StatusCode != 200 {
		t.Errorf("obtained invald status, received response: \n%v", *res)
	}
}

func TestSave(t *testing.T) {
	form := url.Values{}
	form.Add("body", "testing testing 123")
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/save/test", strings.NewReader(form.Encode()))
	env := Env{DB: &mockDB{}}

	http.HandlerFunc(env.CheckUUID(CheckPath(env.Save))).ServeHTTP(res, req)

	if res.Result().StatusCode != 302 {
		t.Errorf("obtained invald status, received response: \n%v", *res)
	}
}

func TestUpload(t *testing.T) {
	req, _ := http.NewRequest("GET", "/upload", nil)
	res := httptest.NewRecorder()
	env := Env{DB: &mockDB{}}

	http.HandlerFunc(env.CheckUUID(Upload)).ServeHTTP(res, req)

	if res.Result().StatusCode != 200 {
		t.Errorf("obtained invalid status, received responsed: \n%v", *res)
	}

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	f, err := os.Open("main_test.go")
	if err != nil {
		return
	}
	defer f.Close()
	fw, err := w.CreateFormFile("text", "main_test.go")
	if err != nil {
		return
	}
	if _, err = io.Copy(fw, f); err != nil {
		return
	}
	if fw, err = w.CreateFormField("myfiles"); err != nil {
		return
	}
	if _, err = fw.Write([]byte("myfiles")); err != nil {
		return
	}
	w.Close()

	req, _ = http.NewRequest("POST", "/upload", &b)
	req.Header.Set("Content-Type", w.FormDataContentType())
	res = httptest.NewRecorder()

	http.HandlerFunc(env.CheckUUID(Upload)).ServeHTTP(res, req)

	if res.Result().StatusCode != 200 {
		t.Errorf("invalid status code, received response: \n%v", *res)
	}
}

func TestIndexPage(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	res := httptest.NewRecorder()
	env := Env{DB: &mockDB{}}

	http.HandlerFunc(env.IndexPage).ServeHTTP(res, req)

	if res.Result().StatusCode != 302 {
		t.Errorf("invalid response status, received response: \n%v", *res)
	}
}

func TestLogin(t *testing.T) {
	form := url.Values{}
	form.Add("uname", "test")
	form.Add("password", "test")
	req, _ := http.NewRequest("POST", "/login", strings.NewReader(form.Encode()))
	req.PostForm = form
	res := httptest.NewRecorder()
	env := Env{DB: &mockDB{}}

	http.HandlerFunc(env.Login).ServeHTTP(res, req)

	if res.Result().StatusCode != 302 {
		t.Errorf("invalid status code, response: \n%v", req.Form)
	}
}
*/
