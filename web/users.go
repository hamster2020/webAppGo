package web

import (
	"net/http"
	"strings"

	"github.com/hamster2020/webAppGo"
)

// Account handles both GET and POST methods, used for updating account info
func (env *Env) Account(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		sessionID := webAppGo.GetSessionIDFromCookie(req)
		session := env.DB.GetSessionFromSessionID(sessionID)
		u, err := env.DB.GetUserFromUserID(session.UserID)
		if err != nil {
			http.Error(res, http.StatusText(500), 500)
			return
		}
		u.Errors = make(map[string]string)
		u.Errors["lname"], _ = webAppGo.GetMsg(res, req, "lname")
		u.Errors["fname"], _ = webAppGo.GetMsg(res, req, "fname")
		u.Errors["username"], _ = webAppGo.GetMsg(res, req, "username")
		u.Errors["email"], _ = webAppGo.GetMsg(res, req, "email")
		u.Errors["password"], _ = webAppGo.GetMsg(res, req, "password")
		env.Render(res, "account", u)
	case "POST":

		sessionID := webAppGo.GetSessionIDFromCookie(req)
		session := env.DB.GetSessionFromSessionID(sessionID)
		u, err := env.DB.GetUserFromUserID(session.UserID)
		if err != nil {
			http.Error(res, http.StatusText(500), 500)
			return
		}

		if u.Username != req.FormValue("userName") {
			n, err := env.DB.CheckUser(req.FormValue("userName"))
			if err != nil {
				http.Error(res, http.StatusText(500), 500)
				return
			}
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
			u.Password, err = webAppGo.EncryptPass(u.Password)
			if err != nil {
				http.Error(res, err.Error(), 500)
				return
			}
			env.DB.UpdateUser(u)
			http.Redirect(res, req, "/", 302)
			return
		}
		http.Redirect(res, req, "/account", 302)
	}
}

// DeleteAccount will delete the User from the db and redirect them to the home page
func (env *Env) DeleteAccount(res http.ResponseWriter, req *http.Request) {
	sessionID := webAppGo.GetSessionIDFromCookie(req)
	session := env.DB.GetSessionFromSessionID(sessionID)
	user, err := env.DB.GetUserFromUserID(session.UserID)
	if err != nil {
		http.Error(res, http.StatusText(500), 500)
		return
	}
	err = env.DB.DeleteUser(user)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	env.DB.DeleteSession(res, sessionID)
	http.Redirect(res, req, "/", 302)
}
