package web

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestLogin(t *testing.T) {
	form := url.Values{}
	form.Add("uname", "test")
	form.Add("password", "test")
	req, _ := http.NewRequest("POST", "/login", strings.NewReader(form.Encode()))
	req.PostForm = form
	res := httptest.NewRecorder()

	env := Env{
		DB:           &mockDB{},
		Cache:        &mockCache{Path: "../../cache"},
		TemplatePath: "../ui/templates/",
		FilePath:     "../files/",
	}

	http.HandlerFunc(env.Login).ServeHTTP(res, req)

	if res.Result().StatusCode != 302 {
		t.Errorf("invalid status code, response: \n%v", *req)
	}
}

func TestLogout(t *testing.T) {
	req, _ := http.NewRequest("GET", "/logout", nil)
	res := httptest.NewRecorder()

	env := Env{
		DB:           &mockDB{},
		Cache:        &mockCache{Path: "../../cache"},
		TemplatePath: "../ui/templates/",
		FilePath:     "../files/",
	}

	http.HandlerFunc(env.Logout).ServeHTTP(res, req)

	if res.Result().StatusCode != 302 {
		t.Errorf("invalid status code, response: \n%v", req.Form)
	}
}

func TestGetIP(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.RemoteAddr = "[::1]:374742018"

	env := Env{
		DB:           &mockDB{},
		Cache:        &mockCache{Path: "../../cache"},
		TemplatePath: "../ui/templates/",
		FilePath:     "../files/",
	}

	ip, err := env.GetIP(req)
	if err != nil {
		t.Errorf("GetIP returned the following error: %s", err.Error())
	}

	if ip != "192.168.0.13" {
		t.Errorf("invalid ip address, wanted: 192.168.0.13, got: %s", ip)
	}
}
