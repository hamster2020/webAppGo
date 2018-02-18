package web

import (
	"net/http"
	"strings"

	"github.com/hamster2020/webAppGo"
)

// View is a function handler for handling http requests
func (env *Env) View(res http.ResponseWriter, req *http.Request, title string) {
	p, err := webAppGo.LoadPageFromCache(title)
	if err != nil {
		p, _ = env.DB.LoadPage(title)
	}
	if p.Title == "" {
		p, _ = env.DB.LoadPage(title)
	}
	if p == nil {
		http.NotFound(res, req)
		return
	}
	if strings.Contains(p.Title, "_") {
		p.Title = strings.Replace(p.Title, "_", " ", -1)
	}
	Render(res, "view", p)
}

// Edit is a function handler for handling http requests
func (env *Env) Edit(res http.ResponseWriter, req *http.Request, title string) {
	p, err := webAppGo.LoadPageFromCache(title)
	if err != nil {
		p, _ = env.DB.LoadPage(title)
	}
	if p.Title == "" {
		p, _ = env.DB.LoadPage(title)
	}
	if p == nil {
		http.NotFound(res, req)
		return
	}
	if strings.Contains(p.Title, "_") {
		p.Title = strings.Replace(p.Title, "_", " ", -1)
	}
	Render(res, "edit", p)
}

// Save is a function handler for making HTTP post requests of Pages to the server
func (env *Env) Save(res http.ResponseWriter, req *http.Request, title string) {
	title = strings.Replace(strings.Title(title), " ", "_", -1)
	body := []byte(req.FormValue("body"))
	page := &webAppGo.Page{Title: title, Body: body}
	page.SaveToCache()
	env.DB.SavePage(page)
	http.Redirect(res, req, "/view/"+title, http.StatusFound)
}

// Create is for creating new pages
func (env *Env) Create(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		p := &webAppGo.Page{}
		Render(res, "create", p)
	case "POST":
		title := strings.Title(req.FormValue("title"))
		if strings.Contains(title, " ") {
			title = strings.Replace(title, " ", "_", -1)
		}
		body := req.FormValue("body")
		p := &webAppGo.Page{Title: strings.Title(title), Body: []byte(body)}
		err := p.SaveToCache()
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		err = env.DB.SavePage(p)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(res, req, "/view/"+title, 302)
		return
	}
}
