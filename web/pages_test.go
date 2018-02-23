package web

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// web/pages.go tests
func TestView(t *testing.T) {
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/view/test", nil)

	env := Env{
		DB:           &mockDB{},
		Cache:        &mockCache{Path: "../../cache"},
		TemplatePath: "../ui/templates/",
		FilePath:     "../files",
	}

	http.HandlerFunc(env.CheckUUID(env.CheckPath(env.View))).ServeHTTP(res, req)

	if res.Result().StatusCode != 200 {
		t.Errorf("obtained invald status, received response: \n%v", *res)
	}
}

func TestEdit(t *testing.T) {
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/edit/test", nil)
	env := Env{DB: &mockDB{}, Cache: &mockCache{Path: "../../cache"}, TemplatePath: "../ui/templates/"}

	http.HandlerFunc(env.CheckUUID(env.CheckPath(env.Edit))).ServeHTTP(res, req)

	if res.Result().StatusCode != 200 {
		t.Errorf("obtained invald status, received response: \n%v", *res)
	}
}

func TestSave(t *testing.T) {
	form := url.Values{}
	form.Add("body", "testing testing 123")
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/save/test", strings.NewReader(form.Encode()))
	req.PostForm = form
	env := Env{DB: &mockDB{}, Cache: &mockCache{Path: "../../cache"}, TemplatePath: "../ui/templates/"}

	http.HandlerFunc(env.CheckUUID(env.CheckPath(env.Save))).ServeHTTP(res, req)

	if res.Result().StatusCode != 302 {
		t.Errorf("obtained invald status, received response: \n%v", *res)
	}
}

func TestCreate(t *testing.T) {
	req, _ := http.NewRequest("GET", "/create/", nil)
	res := httptest.NewRecorder()
	env := Env{DB: &mockDB{}, Cache: &mockCache{Path: "../../cache"}, TemplatePath: "../ui/templates/"}

	http.HandlerFunc(env.CheckUUID(env.Create)).ServeHTTP(res, req)

	if res.Result().StatusCode != 200 {
		t.Errorf("invalid status code, response: \n%v", res)
	}

	form := url.Values{}
	form.Add("title", "test")
	form.Add("body", "testing teseting 123")
	req, _ = http.NewRequest("POST", "/create/", strings.NewReader(form.Encode()))
	res = httptest.NewRecorder()

	http.HandlerFunc(env.CheckUUID(env.Create)).ServeHTTP(res, req)

	if res.Result().StatusCode != 302 {
		t.Errorf("invalid status code, response: \n%v", res)
	}
}
