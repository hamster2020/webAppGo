package main

import (
	"log"
	"net/http"

	"github.com/hamster2020/webAppGo/cache"
	"github.com/hamster2020/webAppGo/sqlite"
	"github.com/hamster2020/webAppGo/web"
)

func main() {
	db, err := sqlite.NewDB(sqlite.DataSourceDriver, sqlite.DataSourceName)
	log.Println("main.go: opening db connection pool...")
	if err != nil {
		log.Println("main.go: db connection failed!")
		panic(err)
	}
	log.Println("main.go: db connection established!")

	c := cache.NewCache("../../cache")

	env := &web.Env{
		DB:           db,
		Cache:        c,
		TemplatePath: "../../ui/templates/",
		FilePath:     "../../files/",
	}

	log.Println("main.go: configuring server...")
	mux := web.NewServer(env)
	log.Println("main.go: server properly initialized!")
	log.Println("main.go: server listing on port 8000...")
	http.ListenAndServe(":8000", mux)
}
