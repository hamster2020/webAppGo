package main

import (
	"net/http"

	"github.com/hamster2020/webAppGo/web"
)

func main() {
	// set up and initialize environment variables
	logPath := ""
	cachePath := "../../cache"
	templatePath := "../../ui/templates/"
	filePath := "../../files/"
	env := web.InitEnv(logPath, cachePath, templatePath, filePath)

	// set up http server
	mux := web.NewServer(env)
	http.ListenAndServe(":8000", mux)
}
