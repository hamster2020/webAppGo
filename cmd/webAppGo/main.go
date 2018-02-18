package main

import (
	"net/http"

	"github.com/hamster2020/webAppGo/sqlite"
	"github.com/hamster2020/webAppGo/web"
)

func main() {
	db, err := sqlite.NewDB(sqlite.DataSourceDriver, sqlite.DataSourceName)
	if err != nil {
		panic(err)
	}

	env := &web.Env{DB: db}
	mux := web.NewServer(env)
	http.ListenAndServe(":8000", mux)
}
