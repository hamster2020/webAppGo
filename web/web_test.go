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
	env := Env{DB: &mockDB{}, Cache: &mockCache{Path: "../../cache"}, TemplatePath: "../ui/templates/"}

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
