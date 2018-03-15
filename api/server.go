package api

import "net/http"

// NewServer returns a server with predefined routing
func NewServer(env *Env) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/pages", env.Pages)
	mux.HandleFunc("/page", env.Page)
	return mux
}
