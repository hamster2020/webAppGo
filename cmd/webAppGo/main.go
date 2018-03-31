package main

import (
	"log"
	"net/http"

	"github.com/hamster2020/webAppGo"
	"github.com/hamster2020/webAppGo/api"
	"github.com/hamster2020/webAppGo/cache"
	postgres "github.com/hamster2020/webAppGo/postgres"
	"github.com/hamster2020/webAppGo/sqlite"
	"github.com/hamster2020/webAppGo/web"
)

func main() {
	// set up and initialize environment variables
	// dbType: choose either "sqlite3" or "postgres"
	// dbPath: choose either "../../sqlite/db.sqlite3" or "postgres://hamster2020:password@localhost/webappgo?sslmode=disable"
	// logPath: provide "" to write log to console, otherwise provide the file path
	dbType := "sqlite3"
	dbPath := "../../sqlite/db.sqlite3"
	logPath := ""
	cachePath := "../../cache/"
	templatePath := "../../ui/templates/"
	filePath := "../../files/"
	webEnv := InitWebEnv(dbType, dbPath, logPath, cachePath, templatePath, filePath)
	apiEnv := InitApiEnv(dbType, dbPath, logPath, cachePath, templatePath, filePath)

	// set up http and api servers
	webServer := web.NewServer(webEnv)
	apiServer := api.NewServer(apiEnv)
	go http.ListenAndServe(":8000", webServer)
	http.ListenAndServe(":8001", apiServer)
}

// InitWebEnv initializes the web environment variables
func InitWebEnv(dbType, dbPath, logPath, cachePath, templatePath, filePath string) *web.Env {
	var db webAppGo.Datastore
	var err error

	log.SetFlags(log.Ldate | log.Lmicroseconds)
	logger := webAppGo.Logger{Level: 1, FilePath: logPath}
	logger.SetSource()

	c := cache.NewCache(cachePath)
	switch dbType {
	case "sqlite3":
		db, err = sqlite.NewDB(dbType, dbPath)
		if err != nil {
			panic(err)
		}
	case "postgres":
		db, err = postgres.NewDB(dbType, dbPath)
		if err != nil {
			panic(err)
		}
	}

	return &web.Env{
		DB:           db,
		Cache:        c,
		Log:          logger,
		TemplatePath: templatePath,
		FilePath:     filePath,
	}
}

// InitApiEnv sets up the api environement variables
func InitApiEnv(dbType, dbPath, logPath, cachePath, templatePath, filePath string) *api.Env {
	var db webAppGo.Datastore
	var err error

	log.SetFlags(log.Ldate | log.Lmicroseconds)
	logger := webAppGo.Logger{Level: 1, FilePath: logPath}
	logger.SetSource()

	c := cache.NewCache(cachePath)

	switch dbType {
	case "sqlite3":
		db, err = sqlite.NewDB(dbType, dbPath)
		if err != nil {
			panic(err)
		}
	case "postgres":
		db, err = postgres.NewDB(dbType, dbPath)
		if err != nil {
			panic(err)
		}
	}

	return &api.Env{
		DB:           db,
		Cache:        c,
		Log:          logger,
		TemplatePath: templatePath,
		FilePath:     filePath,
	}
}
