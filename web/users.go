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
		env.Log.V(1, "beginning handling of GET request for Account.")
		sessionID := webAppGo.GetSessionIDFromCookie(req)
		session := env.DB.GetSessionFromSessionID(sessionID)
		u, err := env.DB.GetUserFromUserID(session.UserID)
		if err != nil {
			env.Log.V(1, "notifying client that an internal error occured. Error is related to  DB.GetUserFromUserID")
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
		env.Log.V(1, "beginning handling of POST request for Account.")
		sessionID := webAppGo.GetSessionIDFromCookie(req)
		session := env.DB.GetSessionFromSessionID(sessionID)
		u, err := env.DB.GetUserFromUserID(session.UserID)
		if err != nil {
			env.Log.V(1, "notifying client that an internal error occured. Error is related to  DB.GetUserFromUserID.")
			http.Error(res, http.StatusText(500), 500)
			return
		}

		if u.Username != req.FormValue("userName") {
			n, err := env.DB.CheckUser(req.FormValue("userName"))
			if err != nil {
				env.Log.V(1, "notifying client that an internal error occured. Error is related to  DB.CheckUser.")
				http.Error(res, http.StatusText(500), 500)
				return
			}
			if n == true {
				env.Log.V(1, "Client is attempting to create a new user account where the username already exists.")
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
			env.Log.V(1, "The passwords that the client provided do not match.")
			webAppGo.SetMsg(res, "password", "The passwords you entered do not match!")
			http.Redirect(res, req, "/account", 302)
			return
		}

		if result == true {
			u.Password, err = webAppGo.EncryptPass(u.Password)
			if err != nil {
				env.Log.V(1, "notifying client that an internal error occured. Error is related to webAppGo.EncryptPass.")
				http.Error(res, err.Error(), 500)
				return
			}
			env.Log.V(1, "The requested updates to the specified user has been accepted. Rediecting to /.")
			env.DB.UpdateUser(u)
			http.Redirect(res, req, "/", 302)
			return
		}
		http.Redirect(res, req, "/account", 302)
	}
}

// DeleteAccount will delete the User from the db and redirect them to the home page
func (env *Env) DeleteAccount(res http.ResponseWriter, req *http.Request) {
	env.Log.V(1, "Beginning handling of DeleteAccount.")
	sessionID := webAppGo.GetSessionIDFromCookie(req)
	session := env.DB.GetSessionFromSessionID(sessionID)
	user, err := env.DB.GetUserFromUserID(session.UserID)
	if err != nil {
		env.Log.V(1, "notifying client that an internal error occured. Error is related to DB.GetUserFromUserID.")
		http.Error(res, http.StatusText(500), 500)
		return
	}
	err = env.DB.DeleteUser(user)
	if err != nil {
		env.Log.V(1, "notifying client that an internal error occured. Error is related to DB.DeleteUser.")
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	env.DB.DeleteSession(res, sessionID)
	env.Log.V(1, "successfully removed both username and session from db. Redirecting to /.")
	http.Redirect(res, req, "/", 302)
}
