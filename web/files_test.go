package web

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestUpload(t *testing.T) {
	req, _ := http.NewRequest("GET", "/upload", nil)
	res := httptest.NewRecorder()

	env := Env{
		DB:           &mockDB{},
		Cache:        &mockCache{Path: "../../cache"},
		TemplatePath: "../ui/templates/",
		FilePath:     "../files/",
	}

	http.HandlerFunc(env.CheckUUID(env.Upload)).ServeHTTP(res, req)

	if res.Result().StatusCode != 200 {
		t.Errorf("obtained invalid status, received responsed: \n%v", *res)
	}

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	f, err := os.Open("files_test.go")
	if err != nil {
		return
	}
	defer f.Close()
	fw, err := w.CreateFormFile("text", "files_test.go")
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

	http.HandlerFunc(env.CheckUUID(env.Upload)).ServeHTTP(res, req)

	if res.Result().StatusCode != 200 {
		t.Errorf("invalid status code, received response: \n%v", *res)
	}
}

func TestDisplayFiles(t *testing.T) {
	req, _ := http.NewRequest("GET", "/display", nil)
	res := httptest.NewRecorder()

	env := Env{
		DB:           &mockDB{},
		Cache:        &mockCache{Path: "../../cache"},
		TemplatePath: "../ui/templates/",
		FilePath:     "../files/",
	}

	http.HandlerFunc(env.CheckUUID(env.DisplayFiles)).ServeHTTP(res, req)

	if res.Result().StatusCode != 200 {
		t.Errorf("invalid status code, received response: \n%v", res)
	}
}

/*
func TestDownload(t *testing) {

}
*/
