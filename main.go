package main

import (
	"html/template"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

var validPath = regexp.MustCompile(`^/(edit|save|view)/([:\w+:]+)$`)

// view is a function handler for handling http requests
func view(res http.ResponseWriter, req *http.Request, title string) {
	p, err := loadSource(title)
	if err != nil {
		p, _ = load(title)
	}
	if p.Title == "" {
		p, _ = load(title)
	}
	if p == nil {
		http.NotFound(res, req)
		return
	}
	if strings.Contains(p.Title, "_") {
		p.Title = strings.Replace(p.Title, "_", " ", -1)
	}
	render(res, "test", p)
}

// edit is a function handler for handling http requests
func edit(res http.ResponseWriter, req *http.Request, title string) {
	p, err := loadSource(title)
	if err != nil {
		p, _ = load(title)
	}
	if p.Title == "" {
		p, _ = load(title)
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
func save(res http.ResponseWriter, req *http.Request, title string) {
	title = strings.Replace(strings.Title(title), " ", "_", -1)
	body := []byte(req.FormValue("body"))
	page := &Page{Title: title, Body: body}
	page.SaveCache()
	http.Redirect(res, req, "/view/"+title, http.StatusFound)
}

// upload is a function that allows the user to upload files to the server
func upload(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		title := "Upload"
		p := &Page{Title: title}
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
			f, err := os.Create("./files/" + files[i].Filename)
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
func indexPage(res http.ResponseWriter, req *http.Request) {
	sessionID := getSessionIDFromCookie(req)
	if foundSessionID, _ := sessionIsValid(res, sessionID); foundSessionID == true {
		http.Redirect(res, req, "/example", 302)
		return
	}
	msg := getMsg(res, req, "msg")
	var u = &User{}
	u.Errors = make(map[string]string)
	if msg != "" {
		u.Errors["message"] = msg
		render(res, "signin", u)
	} else {
		u := &User{}
		render(res, "signin", u)
	}
}

// login is the function handler for POST requests to login
func login(res http.ResponseWriter, req *http.Request) {
	u := &User{
		Username: req.FormValue("uname"),
		Password: req.FormValue("password"),
	}
	redirect := "/"
	if u.Username != "" && u.Password != "" {
		b, userID := userExists(u)
		if b == true {
			s := &Session{
				SessionID: uuid(),
				UserID:    userID,
				Time:      int(time.Now().Unix()),
			}
			setSession(s, res)
			redirect = "/example"
		} else {
			setMsg(res, "msg", "Please signup or enter a valid username and password!")
		}
	} else {
		setMsg(res, "msg", "Username or Password field are empty!")
	}
	http.Redirect(res, req, redirect, 302)
}

// logout merely clears the session cookie and redirects to the index endnode
func logout(res http.ResponseWriter, req *http.Request) {
	sessionID := getSessionIDFromCookie(req)
	clearSession(res, sessionID)
	http.Redirect(res, req, "/", 302)
}

// examplePage is a function hanlder for when the user successfully sets up a session
func examplePage(res http.ResponseWriter, req *http.Request) {
	sessionID := getSessionIDFromCookie(req)
	session := getSessionFromSessionID(sessionID)
	u := getUserFromUUID(session.UserID)
	if session.UserID != "" {
		render(res, "internal", u)
	} else {
		setMsg(res, "msg", "Please login first!")
		http.Redirect(res, req, "/", 302)
	}
}

// signup handles both GEt and POST methods, used for creating and sumitting signups
func signup(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		u := &User{}
		u.Errors = make(map[string]string)
		u.Errors["lname"] = getMsg(res, req, "lname")
		u.Errors["fname"] = getMsg(res, req, "fname")
		u.Errors["username"] = getMsg(res, req, "username")
		u.Errors["email"] = getMsg(res, req, "email")
		u.Errors["password"] = getMsg(res, req, "password")
		render(res, "signup", u)
	case "POST":
		n := checkUser(req.FormValue("userName"))
		if n == true {
			setMsg(res, "username", "User already exists. Please enter a unqiue user name!")
			http.Redirect(res, req, "/signup", 302)
			return
		}
		u := &User{
			UUID:     uuid(),
			Fname:    req.FormValue("fName"),
			Lname:    req.FormValue("lName"),
			Email:    req.FormValue("email"),
			Username: req.FormValue("userName"),
			Password: req.FormValue("password"),
		}
		result, err := ValidateUser(u)
		if err != nil {
			if strings.Contains(err.Error(), "Lname") {
				setMsg(res, "lname", "The name, "+u.Lname+" is not valid!")
			}
			if strings.Contains(err.Error(), "Fname") {
				setMsg(res, "fname", "The name, "+u.Fname+" is not valid!")
			}
			if strings.Contains(err.Error(), "Username") {
				setMsg(res, "username", "The username, "+u.Username+" is not valid!")
			}
			if strings.Contains(err.Error(), "Email") {
				setMsg(res, "email", "The email, "+u.Email+" is not valid!")
			}
			if strings.Contains(err.Error(), "Password") {
				setMsg(res, "password", "Enter a valid passowrd!")
			}
		}
		if req.FormValue("password") != req.FormValue("cpassword") {
			setMsg(res, "password", "The passwords you entered do not match!")
			http.Redirect(res, req, "/signup", 302)
			return
		}
		if result == true {
			u.Password = encryptPass(u.Password)
			saveUserData(u)
			http.Redirect(res, req, "/", 302)
			return
		}
		http.Redirect(res, req, "/signup", 302)
	}
}

func create(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		p := &Page{}
		render(res, "create", p)
	case "POST":
		title := strings.Title(req.FormValue("title"))
		if strings.Contains(title, " ") {
			title = strings.Replace(title, " ", "_", -1)
		}
		body := req.FormValue("body")
		p := &Page{Title: strings.Title(title), Body: []byte(body)}
		err := p.SaveCache()
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
		http.Redirect(res, req, "/view/"+title, 302)
		return
	}
}

func search(res http.ResponseWriter, req *http.Request) {
	sValue := req.FormValue("search")
	sValue = strings.Title(sValue)
	sValue = strings.Replace(sValue, " ", "_", -1)
	if b, _ := pageExists(strings.Title(sValue)); b == true {
		http.Redirect(res, req, "/view/"+sValue, 302)
		return
	}
	render(res, "search", &Page{Title: strings.Title(sValue)})
}

// render is used to write html templates to the response writer
func render(res http.ResponseWriter, name string, data interface{}) {
	funcMap := template.FuncMap{
		"urlize": func(s string) string { return strings.Replace(s, " ", "_", -1) },
	}
	tmpl, err := template.New(name).Funcs(funcMap).ParseGlob("templates/*.html")
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
	tmpl.ExecuteTemplate(res, name, data)
}

func checkUUID(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		sessionID := getSessionIDFromCookie(req)
		if foundSessionID, _ := sessionIsValid(res, sessionID); foundSessionID == true {
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
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir("files"))))
	http.HandleFunc("/", indexPage)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/example", examplePage)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/view/", checkUUID(checkPath(view)))
	http.HandleFunc("/edit/", checkUUID(checkPath(edit)))
	http.HandleFunc("/save/", checkUUID(checkPath(save)))
	http.HandleFunc("/upload/", checkUUID(upload))
	http.HandleFunc("/create/", checkUUID(create))
	http.HandleFunc("/search", checkUUID(search))
	http.ListenAndServe(":8000", nil)
}
