package web

import (
	"net/http"
	"strings"

	"github.com/hamster2020/webAppGo"
)

// View is a function handler for handling http requests
func (env *Env) View(res http.ResponseWriter, req *http.Request, title string) {
	env.Log.V(1, "beginning handling of View.")
	env.Log.V(1, "loading requested page from cache.")
	p, err := env.Cache.LoadPageFromCache(title)
	if err != nil {
		env.Log.V(1, "if file from cache not found, then retrieve requested page from db.")
		p, _ = env.DB.LoadPage(title)
	}
	if p.Title == "" {
		env.Log.V(1, "if page title from db is blank, seach again.")
		p, _ = env.DB.LoadPage(title)
	}
	if p == nil {
		env.Log.V(1, "notifying client that the request page was not found.")
		http.NotFound(res, req)
		return
	}
	if strings.Contains(p.Title, "_") {
		p.Title = strings.Replace(p.Title, "_", " ", -1)
	}
	env.Log.V(1, "requested page found, rendering the view template.")
	env.Render(res, "view", p)
}

// Edit is a function handler for handling http requests
func (env *Env) Edit(res http.ResponseWriter, req *http.Request, title string) {
	env.Log.V(1, "beginning handling of Edit.")
	env.Log.V(1, "loading requested page from cache.")
	p, err := env.Cache.LoadPageFromCache(title)
	if err != nil {
		env.Log.V(1, "if file from cache not found, then retrieve requested page from db.")
		p, _ = env.DB.LoadPage(title)
	}
	if p.Title == "" {
		env.Log.V(1, "if page title is blank, then try again.")
		p, _ = env.DB.LoadPage(title)
	}
	if p == nil {
		env.Log.V(1, "notifying client that the request page was not found.")
		http.NotFound(res, req)
		return
	}
	if strings.Contains(p.Title, "_") {
		p.Title = strings.Replace(p.Title, "_", " ", -1)
	}
	env.Log.V(1, "requested page found, rendering the edit template.")
	env.Render(res, "edit", p)
}

// Save is a function handler for making HTTP post requests of Pages to the server
func (env *Env) Save(res http.ResponseWriter, req *http.Request, title string) {
	env.Log.V(1, "beginning hanlding of Save.")
	title = strings.Replace(strings.Title(title), " ", "_", -1)
	body := []byte(req.FormValue("body"))
	page := &webAppGo.Page{Title: title, Body: body}
	err := env.Cache.SaveToCache(page)
	if err != nil {
		env.Log.V(1, "notifying client that an internal error occured. Error is associated with Cache.SaveToCache.")
		http.Error(res, err.Error(), 500)
	}
	err = env.DB.SavePage(page)
	if err != nil {
		env.Log.V(1, "notifying client that an internal error occured. Error is associated with Cache.SavePage.")
		http.Error(res, err.Error(), 500)
	}
	env.Log.V(1, "The requested new page was successully saved, redirecting client to /view/PageTitle.")
	http.Redirect(res, req, "/view/"+title, http.StatusFound)
}

// Create is for creating new pages
func (env *Env) Create(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		env.Log.V(1, "beginning handling of GET request for Create endnode.")
		p := &webAppGo.Page{}
		env.Log.V(1, "rendering the create template.")
		env.Render(res, "create", p)
	case "POST":
		env.Log.V(1, "beginning handling of POST request for Create endnode.")
		title := strings.Title(req.FormValue("title"))
		if strings.Contains(title, " ") {
			title = strings.Replace(title, " ", "_", -1)
		}
		body := req.FormValue("body")
		p := &webAppGo.Page{Title: strings.Title(title), Body: []byte(body)}
		err := env.Cache.SaveToCache(p)
		if err != nil {
			env.Log.V(1, "notifying client that an internal error occured. Error assocaited with Cache.SaveToCache.")
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		err = env.DB.SavePage(p)
		if err != nil {
			env.Log.V(1, "notifying client that an internal error occured. Error assocaited with DB.SavePage")
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		env.Log.V(1, "page successfully saved, redirecting the client to /view/PageTitle.")
		http.Redirect(res, req, "/view/"+title, 302)
		return
	}
}
