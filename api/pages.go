package api

import (
	"encoding/json"
	"net/http"
	"strings"
)

// StrPage is the same as a Page, but the body is of type string
type StrPage struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

// StrPages is merely a slice of StrPage
type StrPages []StrPage

// Pages is merely a handler for returning the most recent version of each page in the system.
func (env *Env) Pages(res http.ResponseWriter, req *http.Request) {
	env.Log.V(1, "beginning handling of Pages endpoint.")
	if req.Method != "GET" {
		env.Log.V(1, "notifying client of invalid request method to the Pages endpoint.")
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	env.Log.V(1, "pulling all pages from the db.")
	pages, err := env.DB.AllPages()
	if err != nil {
		env.Log.V(1, "notifying client that an Inernal error occured. Error associated with DB.AllPages.")
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	strPages := StrPages{}
	for _, page := range pages {
		strPages = append(strPages, StrPage{Title: page.Title, Body: string(page.Body)})
	}
	if err != nil {
		env.Log.V(1, "notifying client that an internal error occured. Error associated with josn.Marshal")
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
	env.Log.V(1, "converting pages from structs to json.")
	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(strPages)
}

// Page is designed to receive an http request for CRUD operations on pages.
// to test, run the main.go file and then run the following commands in the terminal:
//
// curl -X POST -H "Content-Type: application/json" -d '{"title":"Testing123", "body":"Testing testing 123!"}' http://localhost:8001/page
// curl -X PUT -H "Content-Type: application/json" -d '{"title":"Testing123", "body":"Testing testing 123!"}' http://localhost:8001/page
// curl -X GET -H "Content-Type: application/json" -d '{"title":"Testing123"}' http://localhost:8001/page
// curl -X DEL -H "Content-Type: application/json" -d '{"title":"Testing123"}' http://localhost:8001/page
func (env *Env) Page(res http.ResponseWriter, req *http.Request) {

	page, err := env.parseJSON(req)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	switch req.Method {
	case "GET":
		env.Log.V(1, "beginning handling of Page with GET request method.")
		env.Log.V(1, "loading requested page from cache.")
		p, err := env.Cache.LoadPageFromCache(page.Title)
		if err != nil {
			env.Log.V(1, "if file from cache not found, then retrieve requested page from db.")
			p, _ = env.DB.LoadPage(page.Title)
		}
		if p == nil {
			env.Log.V(1, "notifying client that the request page was not found.")
			http.NotFound(res, req)
			return
		}
		if p.Title == "" {
			env.Log.V(1, "notifying client that the request page was not found.")
			http.NotFound(res, req)
			return
		}
		if strings.Contains(p.Title, "_") {
			p.Title = strings.Replace(p.Title, "_", " ", -1)
		}
		strPage := &StrPage{Title: p.Title, Body: string(p.Body)}
		env.Log.V(1, "handling response.")
		res.Header().Set("Content-Type", "application/json; charset=UTF-8")
		res.WriteHeader(http.StatusOK)
		json.NewEncoder(res).Encode(strPage)

	case "POST":
		env.Log.V(1, "beginning hanlding of http POST request.")
		env.Log.V(1, "saving page to cache.")
		err := env.Cache.SaveToCache(page)
		if err != nil {
			env.Log.V(1, "notifying client that in internal error occured. Error associated with SaveToCache.")
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		env.Log.V(1, "saving page to db.")
		err = env.DB.SavePage(page)
		if err != nil {
			env.Log.V(1, "notifying client that in internal error occured. Error associated with SavePage.")
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		env.Log.V(1, "handling response.")
		res.Header().Set("Content-Type", "application/json; charset=UTF-8")
		res.WriteHeader(http.StatusCreated)

	case "PUT":
		env.Log.V(1, "beginning hanlding of http PUT request.")
		env.Log.V(1, "checking to see if title is in the pages db table.")
		exists, err := env.DB.PageExists(page.Title)
		if err != nil {
			env.Log.V(1, "notifying client that an internal error occured. error associated with PageExists")
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		if exists != true {
			env.Log.V(1, "notifying client that resource was not found.")
			res.Header().Set("Content-Type", "application/json; charset=UTF-8")
			res.WriteHeader(http.StatusNotFound)
			return
		}
		env.Log.V(1, "saving page to cache.")
		err = env.Cache.SaveToCache(page)
		if err != nil {
			env.Log.V(1, "notifying client that in internal error occured. Error associated with SaveToCache.")
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		env.Log.V(1, "saving page to db.")
		err = env.DB.SavePage(page)
		if err != nil {
			env.Log.V(1, "notifying client that in internal error occured. Error associated with SavePage.")
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		env.Log.V(1, "handling response.")
		res.Header().Set("Content-Type", "application/json; charset=UTF-8")
		res.WriteHeader(http.StatusCreated)

	case "DEL":
		env.Log.V(1, "beginning handling DEL request.")
		env.Log.V(1, "removing requested page from cache.")
		err := env.Cache.DeletePageFromCache(page.Title)
		if err != nil {
			env.Log.V(1, "notifying client that an internal error occured. Error associated with DeletePageFromCache.")
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		env.Log.V(1, "removing requested page from db.")
		err = env.DB.DeletePage(page.Title)
		if err != nil {
			env.Log.V(1, "notifying client that an internal error occured. Error associated with DeletePage.")
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		env.Log.V(1, "handling response.")
		res.Header().Set("Content-Type", "application/json; charset=UTF-8")
		res.WriteHeader(http.StatusOK)

	default:
		env.Log.V(1, "handling http request iwth invalid method type.")
		res.Header().Set("Content-Type", "application/json; charset=UTF-8")
		res.WriteHeader(http.StatusMethodNotAllowed)
	}
}
