package web

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestIndexPage(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	res := httptest.NewRecorder()

	env := Env{
		DB:           &mockDB{},
		Cache:        &mockCache{Path: "../../cache"},
		TemplatePath: "../ui/templates/",
		FilePath:     "../files",
	}

	http.HandlerFunc(env.IndexPage).ServeHTTP(res, req)

	if res.Result().StatusCode != 302 {
		t.Errorf("invalid response status, received response: \n%v", *res)
	}
}

func TestHomePage(t *testing.T) {
	req, _ := http.NewRequest("GET", "/home", nil)
	res := httptest.NewRecorder()

	env := Env{
		DB:           &mockDB{},
		Cache:        &mockCache{Path: "../../cache"},
		TemplatePath: "../ui/templates/",
		FilePath:     "../files",
	}

	http.HandlerFunc(env.CheckUUID(env.HomePage)).ServeHTTP(res, req)

	if res.Result().StatusCode != 200 {
		t.Errorf("invalid response status, received response: \n%v", *res)
	}
}

func TestSignUpGET(t *testing.T) {
	req, _ := http.NewRequest("GET", "/signup", nil)
	res := httptest.NewRecorder()

	env := Env{
		DB:           &mockDB{},
		Cache:        &mockCache{Path: "../../cache"},
		TemplatePath: "../ui/templates/",
		FilePath:     "../files",
	}

	http.HandlerFunc(env.CheckUUID(env.Signup)).ServeHTTP(res, req)

	if res.Result().StatusCode != 200 {
		t.Errorf("invalid response status, received response: \n%v", *res)
	}
}

func TestSignUpPOST(t *testing.T) {
	form := url.Values{}
	form.Add("userName", "bobjoe13")
	form.Add("fName", "bob")
	form.Add("lName", "joe")
	form.Add("email", "bobjoe13@email.com")
	form.Add("password", "123456")
	form.Add("cpassword", "123456")
	req, _ := http.NewRequest("POST", "/signup", strings.NewReader(form.Encode()))
	req.PostForm = form
	res := httptest.NewRecorder()

	env := Env{
		DB:           &mockDB{},
		Cache:        &mockCache{Path: "../../cache"},
		TemplatePath: "../ui/templates/",
		FilePath:     "../files",
	}

	http.HandlerFunc(env.CheckUUID(env.Signup)).ServeHTTP(res, req)

	if res.Result().StatusCode != 302 {
		t.Errorf("invalid response status, received response: \n%v", *res)
	}
}

func TestSearch(t *testing.T) {
	form := url.Values{}
	form.Add("search", "test")
	req, _ := http.NewRequest("POST", "/search", strings.NewReader(form.Encode()))
	req.PostForm = form
	res := httptest.NewRecorder()

	env := Env{
		DB:           &mockDB{},
		Cache:        &mockCache{Path: "../../cache"},
		TemplatePath: "../ui/templates/",
		FilePath:     "../files",
	}

	http.HandlerFunc(env.CheckUUID(env.Search)).ServeHTTP(res, req)

	if res.Result().StatusCode != 302 {
		t.Errorf("invalid response status, received response: \n%v", *res)
	}
}
