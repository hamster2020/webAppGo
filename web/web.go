package web

import (
	"html/template"
	"net/http"
	"regexp"
	"strings"

	"github.com/hamster2020/webAppGo"
)

// Env stores environemnt and application scope data to be easily passed to http handlers
type Env struct {
	DB webAppGo.Datastore
}

var validPath = regexp.MustCompile(`^/(edit|save|view|download)/([:\w+:]+[[.]?[:\w+:]+]?)$`)

// IndexPage returns to the client the index.html page
func (env *Env) IndexPage(res http.ResponseWriter, req *http.Request) {
	//  log.Println("web/web.go: begining handling of IndexPage...")
	sessionID, err := webAppGo.GetSessionIDFromCookie(req)
	if err != nil {
		http.Error(res, http.StatusText(500), 500)
		return
	}
	foundSessionID, _, err := env.DB.IsSessionValid(res, sessionID)
	if err != nil {
		http.Error(res, err.Error(), 500)
		return
	}
	if foundSessionID == true {
		http.Redirect(res, req, "/home", 302)
		return
	}
	msg, _ := webAppGo.GetMsg(res, req, "msg")
	var u = &webAppGo.User{}
	u.Errors = make(map[string]string)
	if msg != "" {
		u.Errors["message"] = msg
		//  log.Println("web/web.go: rendering the signin page with errors...")
		Render(res, "signin", u)
	} else {
		u := &webAppGo.User{}
		//  log.Println("web/web.go: rendering the signin page...")
		Render(res, "signin", u)
	}
}

// HomePage is a function hanlder for when the user successfully sets up a session
func (env *Env) HomePage(res http.ResponseWriter, req *http.Request) {
	//  log.Println("web/web.go: begining handling of HomePage...")
	sessionID, err := webAppGo.GetSessionIDFromCookie(req)
	if err != nil {
		http.Error(res, err.Error(), 500)
	}
	session, err := env.DB.GetSessionFromSessionID(sessionID)
	if err != nil {
		http.Error(res, err.Error(), 500)
	}
	u, err := env.DB.GetUserFromUserID(session.UserID)
	if err != nil {
		http.Error(res, err.Error(), 500)
	}
	if session.UserID != "" {
		//  log.Println("web/web.go: rendering the home page...")
		Render(res, "home", u)
	} else {
		//  log.Println("web/web.go: redirecting client to the index page to login first...")
		webAppGo.SetMsg(res, "msg", "Please login first!")
		http.Redirect(res, req, "/", 302)
	}
}

// Signup handles both GEt and POST methods, used for creating and sumitting signups
func (env *Env) Signup(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		u := &webAppGo.User{}
		u.Errors = make(map[string]string)
		u.Errors["lname"], _ = webAppGo.GetMsg(res, req, "lname")
		u.Errors["fname"], _ = webAppGo.GetMsg(res, req, "fname")
		u.Errors["username"], _ = webAppGo.GetMsg(res, req, "username")
		u.Errors["email"], _ = webAppGo.GetMsg(res, req, "email")
		u.Errors["password"], _ = webAppGo.GetMsg(res, req, "password")
		Render(res, "signup", u)
	case "POST":
		n, err := env.DB.CheckUser(req.FormValue("userName"))
		if err != nil {
			http.Error(res, err.Error(), 500)
		}
		if n == true {
			webAppGo.SetMsg(res, "username", "User already exists. Please enter a unqiue user name!")
			http.Redirect(res, req, "/signup", 302)
			return
		}
		uuid, err := webAppGo.UUID()
		if err != nil {
			http.Error(res, err.Error(), 500)
		}
		u := &webAppGo.User{
			UUID:     uuid,
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
			webAppGo.SetMsg(res, "password", "The passwords you entered do not match!")
			http.Redirect(res, req, "/signup", 302)
			return
		}
		if result == true {
			u.Password, err = webAppGo.EncryptPass(u.Password)
			if err != nil {
				http.Error(res, err.Error(), 500)
			}
			env.DB.SaveUser(u)
			http.Redirect(res, req, "/", 302)
			return
		}
		http.Redirect(res, req, "/signup", 302)
	}
}

// Search returns the view of the provided page client has searched for
func (env *Env) Search(res http.ResponseWriter, req *http.Request) {
	sValue := req.FormValue("search")
	sValue = strings.Title(sValue)
	sValue = strings.Replace(sValue, " ", "_", -1)
	if b, _ := env.DB.PageExists(strings.Title(sValue)); b == true {
		http.Redirect(res, req, "/view/"+sValue, 302)
		return
	}
	Render(res, "search", &webAppGo.Page{Title: strings.Title(sValue)})
}

// Render is used to write html templates to the response writer
func Render(res http.ResponseWriter, name string, data interface{}) {
	//  log.Println("web/web.go: beginning handling of Render...")
	funcMap := template.FuncMap{
		"urlize":   func(s string) string { return strings.Replace(s, " ", "_", -1) },
		"deurlize": func(s string) string { return strings.Replace(s, "_", " ", -1) },
	}
	//  log.Println("web/web.go: parsing GLOB...")
	tmpl, err := template.New(name).Funcs(funcMap).ParseGlob("../../ui/templates/*.html")
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	//  log.Println("web/web.go: Executing template...")
	tmpl.ExecuteTemplate(res, name, data)
}

// CheckUUID checks to see if the client has a cookie with a valid session (loged in)
func (env *Env) CheckUUID(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		sessionID, err := webAppGo.GetSessionIDFromCookie(req)
		if err != nil {
			http.Error(res, err.Error(), 500)
		}
		foundSessionID, _, err := env.DB.IsSessionValid(res, sessionID)
		if err != nil {
			http.Error(res, err.Error(), 500)
		}
		if foundSessionID == true {
			fn(res, req)
			return
		}
		http.Redirect(res, req, "/", 302)
	}
}

// CheckPath checks to see if the path is valid
func CheckPath(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		path := validPath.FindStringSubmatch(req.URL.Path)
		if path == nil {
			http.NotFound(res, req)
			return
		}
		fn(res, req, path[2])
	}
}
