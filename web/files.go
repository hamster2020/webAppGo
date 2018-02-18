package web

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/hamster2020/webAppGo"
)

// Upload is a function that allows the user to upload files to the server
func Upload(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		title := "Upload"
		p := &webAppGo.Page{Title: title}
		Render(res, "upload", p)

	case "POST":
		err := req.ParseMultipartForm(100000)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		m := req.MultipartForm
		files := m.File["myfiles"]
		for i := range files {
			file, err := files[i].Open()
			defer file.Close()
			if err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			f, err := os.Create("../../files/" + files[i].Filename)
			if err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			defer f.Close()
			if err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			if _, err := io.Copy(f, file); err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			http.Redirect(res, req, "/files/"+files[i].Filename, http.StatusFound)
		}
	default:
		res.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// DisplayFiles will query all files in ./files/ and display them to the user
func DisplayFiles(res http.ResponseWriter, req *http.Request) {
	files, err := ioutil.ReadDir("../../files/")
	if err != nil {
		log.Fatal(err)
	}
	Render(res, "displayFiles", files)
}

// Download will download the selected file to the client
func Download(res http.ResponseWriter, req *http.Request, title string) {
	file, err := os.Open("../../files/" + title)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	res.Header().Set("Content-Disposition", "attachment; filename="+title)
	res.Header().Set("Content-Type", req.Header.Get("Content-Type"))
	io.Copy(res, file)
}
