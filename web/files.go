package web

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/hamster2020/webAppGo"
)

// Upload is a function that allows the user to upload files to the server
func (env *Env) Upload(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		env.Log.V(1, "handling GET request for /updload.")
		title := "Upload"
		p := &webAppGo.Page{Title: title}
		env.Log.V(1, "rendering upload template for client.")
		env.Render(res, "upload", p)

	case "POST":
		env.Log.V(1, "handling POST request for /upload")
		err := req.ParseMultipartForm(100000)
		if err != nil {
			env.Log.V(1, "notifying client that an internal error occured. Error is associated with ParseMultiPartForm.")
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		m := req.MultipartForm
		files := m.File["myfiles"]
		for i := range files {
			file, err := files[i].Open()
			defer file.Close()
			if err != nil {
				env.Log.V(1, "notifying client that an internal error occured. Error is associated with file.Open.")
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			f, err := os.Create(env.FilePath + files[i].Filename)
			if err != nil {
				env.Log.V(1, "notifying client that an internal error occured. Error is associated with os.Create.")
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			defer f.Close()
			if err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			if _, err := io.Copy(f, file); err != nil {
				env.Log.V(1, "notifying client that an internal error occured. Error is associated with io.Copy.")
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			env.Log.V(1, "files successfully created, routing user to /view/FileName.")
			http.Redirect(res, req, "/files/"+files[i].Filename, http.StatusFound)
		}
	default:
		env.Log.V(1, "notifying client that the request type is not allowed.")
		res.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// DisplayFiles will query all files in ./files/ and display them to the user
func (env *Env) DisplayFiles(res http.ResponseWriter, req *http.Request) {
	env.Log.V(1, "beginning handling of DisplayFiles.")
	files, err := ioutil.ReadDir(env.FilePath)
	if err != nil {
		env.Log.V(1, "notifying client that an internal error occured. Error is associated with ioutil.ReadDir.")
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
	env.Log.V(1, "successful found files in system, rendering the displayFiles template.")
	env.Render(res, "displayFiles", files)
}

// Download will download the selected file to the client
func (env *Env) Download(res http.ResponseWriter, req *http.Request, title string) {
	env.Log.V(1, "beginning hanlding of Download.")
	file, err := os.Open(env.FilePath)
	if err != nil {
		env.Log.V(1, "notifying client that an internal error occured. Error is associated with os.Open.")
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
	defer file.Close()

	res.Header().Set("Content-Disposition", "attachment; filename="+title)
	res.Header().Set("Content-Type", req.Header.Get("Content-Type"))
	env.Log.V(1, "sending response to client to prompt them to indicate the download path of the file to download.")
	io.Copy(res, file)
}
