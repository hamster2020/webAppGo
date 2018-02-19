package web

import (
	"net/http"
	"time"

	"github.com/hamster2020/webAppGo"
)

// Login is the function handler for POST requests to login
func (env *Env) Login(res http.ResponseWriter, req *http.Request) {
	u := &webAppGo.User{
		Username: req.FormValue("uname"),
		Password: req.FormValue("password"),
	}
	redirect := "/"
	if u.Username != "" && u.Password != "" {

		ok, err := env.DB.CheckUserLoginAttempts(u.Username)
		if err != nil {
			http.Error(res, http.StatusText(500), 500)
			return
		}
		if ok != true {
			webAppGo.SetMsg(res, "msg", "Too many incorrect login attempts were made for the provided username, try again in 10 minutes!")
			http.Redirect(res, req, "/", 302)
			return
		}

		ip, err := GetIP(req)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		ok, err = env.DB.CheckIPLoginAttempts(ip)
		if err != nil {
			http.Error(res, http.StatusText(500), 500)
			return
		}
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

		b, userID, err := env.DB.UserExists(u)
		if err != nil {
			http.Error(res, http.StatusText(500), 500)
			return
		}
		if b == true {
			uuid, err := webAppGo.UUID()
			if err != nil {
				http.Error(res, http.StatusText(500), 500)
				return
			}
			s := &webAppGo.Session{
				SessionID: uuid,
				UserID:    userID,
				Time:      int(time.Now().Unix()),
			}
			login.Attempt = true
			env.DB.SetSession(s, res)
			env.DB.SaveLogin(login)
			redirect = "/home"
		} else {
			login.Attempt = false
			env.DB.SaveLogin(login)
			webAppGo.SetMsg(res, "msg", "Please signup or enter a valid username and password!")
		}
	} else {
		webAppGo.SetMsg(res, "msg", "Username or Password field are empty!")
	}
	http.Redirect(res, req, redirect, 302)
}

// Logout merely clears the session cookie and redirects to the index endnode
func (env *Env) Logout(res http.ResponseWriter, req *http.Request) {
	sessionID := webAppGo.GetSessionIDFromCookie(req)
	env.DB.DeleteSession(res, sessionID)
	webAppGo.ClearCookie(res, "session")
	http.Redirect(res, req, "/", 302)
}

// GetIP is only used for auto-generating a fake IP since this app is not live
func GetIP(req *http.Request) (string, error) {
	return "192.168.0.13", nil
	/*
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
	*/
}
