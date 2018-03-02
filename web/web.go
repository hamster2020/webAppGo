package web

import (
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/hamster2020/webAppGo"
	"github.com/hamster2020/webAppGo/cache"
	"github.com/hamster2020/webAppGo/sqlite"
)

// Env stores environemnt and application scope data to be easily passed to http handlers
type Env struct {
	DB           webAppGo.Datastore
	Cache        webAppGo.PageCache
	Log          webAppGo.Logger
	TemplatePath string
	FilePath     string
}

var validPath = regexp.MustCompile(`^/(edit|save|view|download)/([:\w+:]+[[.]?[:\w+:]+]?)$`)

func InitEnv(logPath, cachePath, templatePath, filePath string) *Env {
	log.SetFlags(log.Ldate | log.Lmicroseconds)
	logger := webAppGo.Logger{Level: 1, FilePath: logPath} //../../logs/"}
	logger.SetSource()

	c := cache.NewCache(cachePath)

	db, err := sqlite.NewDB(sqlite.DataSourceDriver, sqlite.DataSourceName)
	if err != nil {
		panic(err)
	}

	return &Env{
		DB:           db,
		Cache:        c,
		Log:          logger,
		TemplatePath: templatePath,
		FilePath:     filePath,
	}
}

// IndexPage returns to the client the index.html page
func (env *Env) IndexPage(res http.ResponseWriter, req *http.Request) {
	env.Log.V(1, "handling request for Index Page.")
	sessionID := webAppGo.GetSessionIDFromCookie(req)
	foundSessionID, _ := env.DB.IsSessionValid(res, sessionID)
	if foundSessionID == true {
		env.Log.V(1, "client has a valid sessionID, redirecting to the user's home page.")
		http.Redirect(res, req, "/home", 302)
		return
	}
	msg, _ := webAppGo.GetMsg(res, req, "msg")
	var u = &webAppGo.User{}
	u.Errors = make(map[string]string)
	if msg != "" {
		u.Errors["message"] = msg
		env.Log.V(1, "client provided invalid login, rendering the main page with errors.")
		env.Render(res, "signin", u)
	} else {
		u := &webAppGo.User{}
		env.Log.V(1, "client requested index page with no msg flags in cookies, rendering the main page.")
		env.Render(res, "signin", u)
	}
}

// HomePage is a function hanlder for when the user successfully sets up a session
func (env *Env) HomePage(res http.ResponseWriter, req *http.Request) {
	env.Log.V(1, "begining handling of Home Page.")
	sessionID := webAppGo.GetSessionIDFromCookie(req)
	session := env.DB.GetSessionFromSessionID(sessionID)
	u, err := env.DB.GetUserFromUserID(session.UserID)
	if err != nil {
		env.Log.V(1, "Notifying client of error that occured. Associated with GetUserFromUSERID.")
		http.Error(res, err.Error(), 500)
	}
	if session.UserID != "" {
		env.Log.V(1, "rendering the home page.")
		env.Render(res, "home", u)
	} else {
		env.Log.V(1, "redirecting client to the index page to login first.")
		webAppGo.SetMsg(res, "msg", "Please login first!")
		http.Redirect(res, req, "/", 302)
	}
}

// Signup handles both GET and POST methods, used for creating and submitting signups
func (env *Env) Signup(res http.ResponseWriter, req *http.Request) {
	env.Log.V(1, "beginning handling of Signin Page for %s request.", req.Method)
	switch req.Method {
	case "GET":
		u := &webAppGo.User{}
		u.Errors = make(map[string]string)
		u.Errors["lname"], _ = webAppGo.GetMsg(res, req, "lname")
		u.Errors["fname"], _ = webAppGo.GetMsg(res, req, "fname")
		u.Errors["username"], _ = webAppGo.GetMsg(res, req, "username")
		u.Errors["email"], _ = webAppGo.GetMsg(res, req, "email")
		u.Errors["password"], _ = webAppGo.GetMsg(res, req, "password")
		env.Log.V(1, "rendering the signup page with validation notifications.")
		env.Render(res, "signup", u)
	case "POST":
		n, err := env.DB.CheckUser(req.FormValue("userName"))
		if err != nil {
			env.Log.V(1, "Notifying client that error occured. Error associated with CheckUser.")
			http.Error(res, err.Error(), 500)
		}
		if n == true {
			env.Log.V(1, "redirecting to /signup with GET request to nofiy user that the username already exists.")
			webAppGo.SetMsg(res, "username", "User already exists. Please enter a unqiue user name!")
			http.Redirect(res, req, "/signup", 302)
			return
		}
		uuid := webAppGo.GenRandID(32)
		u := &webAppGo.User{
			UserID:   uuid,
			Fname:    req.FormValue("fName"),
			Lname:    req.FormValue("lName"),
			Email:    req.FormValue("email"),
			Username: req.FormValue("userName"),
			Password: req.FormValue("password"),
		}
		result, err := webAppGo.ValidateUser(u)
		if err != nil {
			if strings.Contains(err.Error(), "Lname") {
				webAppGo.SetMsg(res, "lname", "The name, "+u.Lname+" is not valid!")
			}
			if strings.Contains(err.Error(), "Fname") {
				webAppGo.SetMsg(res, "fname", "The name, "+u.Fname+" is not valid!")
			}
			if strings.Contains(err.Error(), "Username") {
				webAppGo.SetMsg(res, "username", "The username, "+u.Username+" is not valid!")
			}
			if strings.Contains(err.Error(), "Email") {
				webAppGo.SetMsg(res, "email", "The email, "+u.Email+" is not valid!")
			}
			if strings.Contains(err.Error(), "Password") {
				webAppGo.SetMsg(res, "password", "Enter a valid passowrd!")
			}
		}
		if req.FormValue("password") != req.FormValue("cpassword") {
			env.Log.V(1, "sending a GET request on the client's behalf to /signup to notify them via temp cookies of validation errors")
			webAppGo.SetMsg(res, "password", "The passwords you entered do not match!")
			http.Redirect(res, req, "/signup", 302)
			return
		}
		if result == true {
			u.Password, err = webAppGo.EncryptPass(u.Password)
			if err != nil {
				env.Log.V(1, "notify client that an internal error occured. Error is associated with webAppGo.EncryptPass")
				http.Error(res, err.Error(), 500)
			}
			env.Log.V(1, "successfully added user to db and redirecting client to /.")
			env.DB.SaveUser(u)
			http.Redirect(res, req, "/", 302)
			return
		}
		env.Log.V(1, "redirecting client to /signup again")
		http.Redirect(res, req, "/signup", 302)
	}
}

// Search returns the view of the provided page client has searched for
func (env *Env) Search(res http.ResponseWriter, req *http.Request) {
	env.Log.V(1, "beginning handling of /seach POST request")
	sValue := req.FormValue("search")
	sValue = strings.Title(sValue)
	sValue = strings.Replace(sValue, " ", "_", -1)
	if b, _ := env.DB.PageExists(strings.Title(sValue)); b == true {
		env.Log.V(1, "page was found, rediecting client to /view/%s.", sValue)
		http.Redirect(res, req, "/view/"+sValue, 302)
		return
	}
	env.Log.V(1, "page was not found, rendering search template.")
	env.Render(res, "search", &webAppGo.Page{Title: strings.Title(sValue)})
}

// Render is used to write html templates to the response writer
func (env *Env) Render(res http.ResponseWriter, name string, data interface{}) {
	env.Log.V(2, "beginning handling of Render handler.")
	funcMap := template.FuncMap{
		"urlize":   func(s string) string { return strings.Replace(s, " ", "_", -1) },
		"deurlize": func(s string) string { return strings.Replace(s, "_", " ", -1) },
	}
	tmpl, err := template.New(name).Funcs(funcMap).ParseGlob(env.TemplatePath + "*.html")
	if err != nil {
		env.Log.V(1, "notifying client that an error occured. Error is associated with parsing templates.")
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.ExecuteTemplate(res, name, data)
}

// CheckUUID checks to see if the client has a cookie with a valid session (loged in)
func (env *Env) CheckUUID(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		sessionID := webAppGo.GetSessionIDFromCookie(req)
		foundSessionID, _ := env.DB.IsSessionValid(res, sessionID)
		env.Log.V(2, "checking if client's session cookie exists and is valid.")
		if foundSessionID == true {
			env.Log.V(2, "sessionID is found and valid. Now calling appropriate handler.")
			fn(res, req)
			return
		}
		env.Log.V(2, "sessionID is not valid, redirecting client back to /.")
		http.Redirect(res, req, "/", 302)
	}
}

// CheckPath checks to see if the path is valid
func (env *Env) CheckPath(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		path := validPath.FindStringSubmatch(req.URL.Path)
		env.Log.V(2, "checking if the url path is valid.")
		if path == nil {
			env.Log.V(2, "the request url path is not valid, returning client an status of not found.")
			http.NotFound(res, req)
			return
		}
		env.Log.V(2, "the requested url path is valid, calling the next appropriate function handler.")
		fn(res, req, path[2])
	}
}
