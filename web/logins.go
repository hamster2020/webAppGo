package web

import (
	"net"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/hamster2020/webAppGo"
)

// Login is the function handler for POST requests to login
func (env *Env) Login(res http.ResponseWriter, req *http.Request) {
	env.Log.V(1, "handling login POST request from client.")
	u := &webAppGo.User{
		Username: req.FormValue("uname"),
		Password: req.FormValue("password"),
	}
	redirect := "/"
	if u.Username != "" && u.Password != "" {

		ok, err := env.DB.CheckUserLoginAttempts(u.Username)
		if err != nil {
			env.Log.V(1, "Notifying client that an internal error occured. Error is associated with DB.CheckUserLoginAttempts.")
			http.Error(res, http.StatusText(500), 500)
			return
		}
		if ok != true {
			env.Log.V(1, "Too many failed login request for the given username have occured within the specified timeframe. Redirecting back to /.")
			webAppGo.SetMsg(res, "msg", "Too many incorrect login attempts were made for the provided username, try again in 10 minutes!")
			http.Redirect(res, req, "/", 302)
			return
		}
		ip, err := env.GetIP(req)
		if err != nil {
			env.Log.V(1, "Notifying client that an internal error occured. Error is associated with DB.GetIP.")
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		ok, err = env.DB.CheckIPLoginAttempts(ip)
		if err != nil {
			env.Log.V(1, "Notifying client that an internal error occured. Error is associated with DB.CheckIPLoginAttempts.")
			http.Error(res, http.StatusText(500), 500)
			return
		}
		if ok != true {
			env.Log.V(1, "Too many failed login request for the given username have occured within the specified timeframe. Redirecting back to /.")
			webAppGo.SetMsg(res, "msg", "Too many incorrect login attempts were made from you, try again in 10 minutes!")
			http.Redirect(res, req, "/", 302)
			return
		}

		login := &webAppGo.Login{
			IP:        ip,
			UserName:  u.Username,
			Timestamp: int(time.Now().Unix()),
		}

		b, userID, err := env.DB.UserExists(u)
		if err != nil {
			env.Log.V(1, "Notifying client that an internal error occured. Error is associated with DB.UserExists.")
			http.Error(res, http.StatusText(500), 500)
			return
		}
		if b == true {
			uuid := webAppGo.GenRandID(32)
			s := &webAppGo.Session{
				SessionID: uuid,
				UserID:    userID,
				Time:      int(time.Now().Unix()),
			}
			login.Attempt = true
			env.DB.SetSession(s, res)
			env.DB.SaveLogin(login)
			env.Log.V(1, "Client provided successful login credentials, redirecting to /home.")
			redirect = "/home"
		} else {
			login.Attempt = false
			env.DB.SaveLogin(login)
			env.Log.V(1, "Client failed to provide valid login credentials, redirecting to /.")
			webAppGo.SetMsg(res, "msg", "Please signup or enter a valid username and password!")
		}
	} else {
		env.Log.V(1, "Client failed to provide either a password or username, redirecting to /.")
		webAppGo.SetMsg(res, "msg", "Username or Password field are empty!")
	}
	http.Redirect(res, req, redirect, 302)
}

// Logout merely clears the session cookie and redirects to the index endnode
func (env *Env) Logout(res http.ResponseWriter, req *http.Request) {
	env.Log.V(1, "handling logout endnode.")
	sessionID := webAppGo.GetSessionIDFromCookie(req)
	env.DB.DeleteSession(res, sessionID)
	webAppGo.ClearCookie(res, "session")
	env.Log.V(1, "the client's session is successfully cleared out, redirecting client back to /.")
	http.Redirect(res, req, "/", 302)
}

// GetIP is only used for auto-generating a fake IP since this app is not live
func (env *Env) GetIP(req *http.Request) (string, error) {
	env.Log.V(2, "beginning handling of GetIP.")
	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		env.Log.V(2, "could not find ip in request, had to retrun hard-coded ip (associated with lack of dns)")
		return "192.168.0.13", nil
	}
	if ip == "::1" || ip == "" {
		cmd := "ip route get 8.8.8.8 | awk '{print $NF; exit}'"
		out, err := exec.Command("bash", "-c", cmd).Output()
		if err != nil {
			return "", err
		}
		ip = strings.TrimSpace(string(out))
	}
	env.Log.V(2, "successfuly obtained the client's ip address.")
	return ip, nil
}
