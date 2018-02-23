package web

import "net/http"

// NewServer returns a server with predefined routing
func NewServer(env *Env) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../../ui/static"))))
	mux.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir(env.FilePath))))
	mux.HandleFunc("/", env.IndexPage)
	mux.HandleFunc("/login", env.Login)
	mux.HandleFunc("/logout", env.Logout)
	mux.HandleFunc("/home", env.HomePage)
	mux.HandleFunc("/signup", env.Signup)
	mux.HandleFunc("/account", env.CheckUUID(env.Account))
	mux.HandleFunc("/delete", env.CheckUUID(env.DeleteAccount))
	mux.HandleFunc("/view/", env.CheckUUID(env.CheckPath(env.View)))
	mux.HandleFunc("/edit/", env.CheckUUID(env.CheckPath(env.Edit)))
	mux.HandleFunc("/save/", env.CheckUUID(env.CheckPath(env.Save)))
	mux.HandleFunc("/upload/", env.CheckUUID(env.Upload))
	mux.HandleFunc("/create/", env.CheckUUID(env.Create))
	mux.HandleFunc("/search", env.CheckUUID(env.Search))
	mux.HandleFunc("/display", env.CheckUUID(env.DisplayFiles))
	mux.HandleFunc("/download/", env.CheckUUID(env.CheckPath(env.Download)))

	return mux
}
