package web

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestAccountGET(t *testing.T) {
	req, _ := http.NewRequest("GET", "/account", nil)
	res := httptest.NewRecorder()

	env := Env{
		DB:           &mockDB{},
		Cache:        &mockCache{Path: "../../cache"},
		TemplatePath: "../ui/templates/",
		FilePath:     "../files",
	}

	http.HandlerFunc(env.CheckUUID(env.Account)).ServeHTTP(res, req)

	if res.Result().StatusCode != 200 {
		t.Errorf("obtained invald status, received response: \n%v", *res)
	}
}

func TestAccountPOST(t *testing.T) {
	form := url.Values{}
	form.Add("userName", "bobjoe13")
	form.Add("fName", "bob")
	form.Add("lName", "joe")
	form.Add("email", "bobjoe13@email.com")
	form.Add("password", "123456")
	form.Add("cpassword", "123456")
	req, _ := http.NewRequest("POST", "/account", strings.NewReader(form.Encode()))
	req.PostForm = form
	res := httptest.NewRecorder()

	env := Env{
		DB:           &mockDB{},
		Cache:        &mockCache{Path: "../../cache"},
		TemplatePath: "../ui/templates/",
		FilePath:     "../files",
	}

	http.HandlerFunc(env.CheckUUID(env.Account)).ServeHTTP(res, req)

	if res.Result().StatusCode != 302 {
		t.Errorf("obtained invald status, received response: \n%v", *res)
	}
}

func TestDeleteAccount(t *testing.T) {
	req, _ := http.NewRequest("GET", "/delete", nil)
	res := httptest.NewRecorder()

	env := Env{
		DB:           &mockDB{},
		Cache:        &mockCache{Path: "../../cache"},
		TemplatePath: "../ui/templates/",
		FilePath:     "../files",
	}

	http.HandlerFunc(env.CheckUUID(env.DeleteAccount)).ServeHTTP(res, req)

	if res.Result().StatusCode != 302 {
		t.Errorf("obtained invald status, received response: \n%v", *res)
	}
}
