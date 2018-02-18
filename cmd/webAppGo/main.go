package main

import (
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/hamster2020/webAppGo"
	"github.com/hamster2020/webAppGo/sqlite"
)

// Env stores environemnt and application scope data to be easily passed to http handlers
type Env struct {
	db webAppGo.Datastore
}

var validPath = regexp.MustCompile(`^/(edit|save|view|download)/([:\w+:]+[[.]?[:\w+:]+]?)$`)

// view is a function handler for handling http requests
func (env *Env) view(res http.ResponseWriter, req *http.Request, title string) {
	p, err := webAppGo.LoadPageFromCache(title)
	if err != nil {
		p, _ = env.db.LoadPage(title)
	}
	if p.Title == "" {
		p, _ = env.db.LoadPage(title)
	}
	if p == nil {
		http.NotFound(res, req)
		return
	}
	if strings.Contains(p.Title, "_") {
		p.Title = strings.Replace(p.Title, "_", " ", -1)
	}
	render(res, "view", p)
}

// edit is a function handler for handling http requests
func (env *Env) edit(res http.ResponseWriter, req *http.Request, title string) {
	p, err := webAppGo.LoadPageFromCache(title)
	if err != nil {
		p, _ = env.db.LoadPage(title)
	}
	if p.Title == "" {
		p, _ = env.db.LoadPage(title)
	}
	if p == nil {
		http.NotFound(res, req)
		return
	}
	if strings.Contains(p.Title, "_") {
		p.Title = strings.Replace(p.Title, "_", " ", -1)
	}
	render(res, "edit", p)
}

// save is a function handler for making HTTP post requests of Pages to the server
func (env *Env) save(res http.ResponseWriter, req *http.Request, title string) {
	title = strings.Replace(strings.Title(title), " ", "_", -1)
	body := []byte(req.FormValue("body"))
	page := &webAppGo.Page{Title: title, Body: body}
	page.SaveToCache()
	env.db.SavePage(page)
	http.Redirect(res, req, "/view/"+title, http.StatusFound)
}

// upload is a function that allows the user to upload files to the server
func upload(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		title := "Upload"
		p := &webAppGo.Page{Title: title}
		render(res, "upload", p)

	case "POST":
		err := req.ParseMultipartForm(100000)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
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

// indexPage is a function handler to return to the client the index.html page
func (env *Env) indexPage(res http.ResponseWriter, req *http.Request) {
	sessionID := webAppGo.GetSessionIDFromCookie(req)
	if foundSessionID, _ := env.db.IsSessionValid(res, sessionID); foundSessionID == true {
		http.Redirect(res, req, "/home", 302)
		return
	}
	msg := webAppGo.GetMsg(res, req, "msg")
	var u = &webAppGo.User{}
	u.Errors = make(map[string]string)
	if msg != "" {
		u.Errors["message"] = msg
		render(res, "signin", u)
	} else {
		u := &webAppGo.User{}
		render(res, "signin", u)
	}
}

// login is the function handler for POST requests to login
func (env *Env) login(res http.ResponseWriter, req *http.Request) {
	u := &webAppGo.User{
		Username: req.FormValue("uname"),
		Password: req.FormValue("password"),
	}
	redirect := "/"
	if u.Username != "" && u.Password != "" {

		ok := env.db.CheckUserLoginAttempts(u.Username)
		if ok != true {
			webAppGo.SetMsg(res, "msg", "Too many incorrect login attempts were made for the provided username, try again in 10 minutes!")
			http.Redirect(res, req, "/", 302)
			return
		}

		ip, err := getIP(req)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}

		ok = env.db.CheckIPLoginAttempts(ip)
		if ok != true {
			webAppGo.SetMsg(res, "msg", "Too many incorrect login attempts were made from you, try again in 10 minutes!")
			http.Redirect(res, req, "/", 302)
			return
		}

		login := &webAppGo.Login{
			IP:        ip,
			UserName:  u.Username,
			Timestamp: int(time.Now().Unix()),
		}

		b, userID := env.db.UserExists(u)
		if b == true {
			s := &webAppGo.Session{
				SessionID: webAppGo.UUID(),
				UserID:    userID,
				Time:      int(time.Now().Unix()),
			}
			login.Attempt = true
			env.db.SetSession(s, res)
			env.db.SaveLogin(login)
			redirect = "/home"
		} else {
			login.Attempt = false
			env.db.SaveLogin(login)
			webAppGo.SetMsg(res, "msg", "Please signup or enter a valid username and password!")
		}
	} else {
		webAppGo.SetMsg(res, "msg", "Username or Password field are empty!")
	}
	http.Redirect(res, req, redirect, 302)
}

// logout merely clears the session cookie and redirects to the index endnode
func (env *Env) logout(res http.ResponseWriter, req *http.Request) {
	sessionID := webAppGo.GetSessionIDFromCookie(req)
	env.db.DeleteSession(res, sessionID)
	webAppGo.ClearCookie(res, "session")
	http.Redirect(res, req, "/", 302)
}

// homePage is a function hanlder for when the user successfully sets up a session
func (env *Env) homePage(res http.ResponseWriter, req *http.Request) {
	sessionID := webAppGo.GetSessionIDFromCookie(req)
	session := env.db.GetSessionFromSessionID(sessionID)
	u := env.db.GetUserFromUserID(session.UserID)
	if session.UserID != "" {
		render(res, "home", u)
	} else {
		webAppGo.SetMsg(res, "msg", "Please login first!")
		http.Redirect(res, req, "/", 302)
	}
}

// signup handles both GEt and POST methods, used for creating and sumitting signups
func (env *Env) signup(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		u := &webAppGo.User{}
		u.Errors = make(map[string]string)
		u.Errors["lname"] = webAppGo.GetMsg(res, req, "lname")
		u.Errors["fname"] = webAppGo.GetMsg(res, req, "fname")
		u.Errors["username"] = webAppGo.GetMsg(res, req, "username")
		u.Errors["email"] = webAppGo.GetMsg(res, req, "email")
		u.Errors["password"] = webAppGo.GetMsg(res, req, "password")
		render(res, "signup", u)
	case "POST":
		n := env.db.CheckUser(req.FormValue("userName"))
		if n == true {
			webAppGo.SetMsg(res, "username", "User already exists. Please enter a unqiue user name!")
			http.Redirect(res, req, "/signup", 302)
			return
		}
		u := &webAppGo.User{
			UUID:     webAppGo.UUID(),
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
			u.Password = webAppGo.EncryptPass(u.Password)
			env.db.SaveUser(u)
			http.Redirect(res, req, "/", 302)
			return
		}
		http.Redirect(res, req, "/signup", 302)
	}
}

// account handles both GET and POST methods, used for updating account info
func (env *Env) account(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		sessionID := webAppGo.GetSessionIDFromCookie(req)
		session := env.db.GetSessionFromSessionID(sessionID)
		u := env.db.GetUserFromUserID(session.UserID)
		u.Errors = make(map[string]string)
		u.Errors["lname"] = webAppGo.GetMsg(res, req, "lname")
		u.Errors["fname"] = webAppGo.GetMsg(res, req, "fname")
		u.Errors["username"] = webAppGo.GetMsg(res, req, "username")
		u.Errors["email"] = webAppGo.GetMsg(res, req, "email")
		u.Errors["password"] = webAppGo.GetMsg(res, req, "password")
		render(res, "account", u)
	case "POST":

		sessionID := webAppGo.GetSessionIDFromCookie(req)
		session := env.db.GetSessionFromSessionID(sessionID)
		u := env.db.GetUserFromUserID(session.UserID)

		if u.Username != req.FormValue("userName") {
			n := env.db.CheckUser(req.FormValue("userName"))
			if n == true {
				webAppGo.SetMsg(res, "username", "User already exists. Please enter a unqiue user name!")
				http.Redirect(res, req, "/account", 302)
				return
			}
		}

		u.Fname = req.FormValue("fName")
		u.Lname = req.FormValue("lName")
		u.Email = req.FormValue("email")
		u.Username = req.FormValue("userName")
		u.Password = req.FormValue("password")

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
			http.Redirect(res, req, "/account", 302)
			return
		}

		if result == true {
			u.Password = webAppGo.EncryptPass(u.Password)
			env.db.UpdateUser(u)
			http.Redirect(res, req, "/", 302)
			return
		}
		http.Redirect(res, req, "/account", 302)
	}
}

func (env *Env) deleteAccount(res http.ResponseWriter, req *http.Request) {
	sessionID := webAppGo.GetSessionIDFromCookie(req)
	session := env.db.GetSessionFromSessionID(sessionID)
	user := env.db.GetUserFromUserID(session.UserID)
	err := env.db.DeleteUser(user)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
	env.db.DeleteSession(res, sessionID)
	http.Redirect(res, req, "/", 302)
}

func (env *Env) create(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		p := &webAppGo.Page{}
		render(res, "create", p)
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
		}
		err = env.db.SavePage(p)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
		http.Redirect(res, req, "/view/"+title, 302)
		return
	}
}

// DisplayFiles will query all files in ./files/ and display them to the user
func displayFiles(res http.ResponseWriter, req *http.Request) {
	files, err := ioutil.ReadDir("../../files/")
	if err != nil {
		log.Fatal(err)
	}
	render(res, "displayFiles", files)
}

func download(res http.ResponseWriter, req *http.Request, title string) {
	file, err := os.Open("../../files/" + title)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	res.Header().Set("Content-Disposition", "attachment; filename="+title)
	res.Header().Set("Content-Type", req.Header.Get("Content-Type"))
	io.Copy(res, file)
}

func (env *Env) search(res http.ResponseWriter, req *http.Request) {
	sValue := req.FormValue("search")
	sValue = strings.Title(sValue)
	sValue = strings.Replace(sValue, " ", "_", -1)
	if b, _ := env.db.PageExists(strings.Title(sValue)); b == true {
		http.Redirect(res, req, "/view/"+sValue, 302)
		return
	}
	render(res, "search", &webAppGo.Page{Title: strings.Title(sValue)})
}

func getIP(req *http.Request) (string, error) {
	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return "", err
	}
	if ip == "::1" {
		cmd := "ip route get 8.8.8.8 | awk '{print $NF; exit}'"
		out, err := exec.Command("bash", "-c", cmd).Output()
		if err != nil {
			return "", err
		}
		ip = strings.TrimSpace(string(out))
	}
	return ip, nil
}

// render is used to write html templates to the response writer
func render(res http.ResponseWriter, name string, data interface{}) {
	funcMap := template.FuncMap{
		"urlize":   func(s string) string { return strings.Replace(s, " ", "_", -1) },
		"deurlize": func(s string) string { return strings.Replace(s, "_", " ", -1) },
	}
	tmpl, err := template.New(name).Funcs(funcMap).ParseGlob("../../templates/*.html")
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
	tmpl.ExecuteTemplate(res, name, data)
}

func (env *Env) checkUUID(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		sessionID := webAppGo.GetSessionIDFromCookie(req)
		if foundSessionID, _ := env.db.IsSessionValid(res, sessionID); foundSessionID == true {
			fn(res, req)
			return
		}
		http.Redirect(res, req, "/", 302)
	}
}

func checkPath(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		path := validPath.FindStringSubmatch(req.URL.Path)
		if path == nil {
			http.NotFound(res, req)
			return
		}
		fn(res, req, path[2])
	}
}

func main() {
	db, err := sqlite.NewDB(sqlite.DataSourceDriver, sqlite.DataSourceName)
	if err != nil {
		panic(err)
	}

	env := Env{db}

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../../static"))))
	http.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir("../../files"))))
	http.HandleFunc("/", env.indexPage)
	http.HandleFunc("/login", env.login)
	http.HandleFunc("/logout", env.logout)
	http.HandleFunc("/home", env.homePage)
	http.HandleFunc("/signup", env.signup)
	http.HandleFunc("/account", env.checkUUID(env.account))
	http.HandleFunc("/delete", env.checkUUID(env.deleteAccount))
	http.HandleFunc("/view/", env.checkUUID(checkPath(env.view)))
	http.HandleFunc("/edit/", env.checkUUID(checkPath(env.edit)))
	http.HandleFunc("/save/", env.checkUUID(checkPath(env.save)))
	http.HandleFunc("/upload/", env.checkUUID(upload))
	http.HandleFunc("/create/", env.checkUUID(env.create))
	http.HandleFunc("/search", env.checkUUID(env.search))
	http.HandleFunc("/display", env.checkUUID(displayFiles))
	http.HandleFunc("/download/", env.checkUUID(checkPath(download)))
	http.ListenAndServe(":8000", nil)
}
