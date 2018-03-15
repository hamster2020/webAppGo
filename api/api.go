package api

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/hamster2020/webAppGo"
)

// Env stores environemnt and application scope data to be easily passed to http handlers
type Env struct {
	DB           webAppGo.Datastore
	Cache        webAppGo.PageCache
	Log          webAppGo.Logger
	TemplatePath string
	FilePath     string
}

func (env *Env) parseJSON(req *http.Request) (*webAppGo.Page, error) {
	var strPage StrPage
	body, err := ioutil.ReadAll(io.LimitReader(req.Body, 1048576))
	if err != nil {
		return nil, err
	}
	if err := req.Body.Close(); err != nil {
		return nil, err
	}
	env.Log.V(1, "decoding json data from body.")
	if err := json.Unmarshal(body, &strPage); err != nil {
		return nil, err
	}
	return &webAppGo.Page{Title: strPage.Title, Body: []byte(strPage.Body)}, nil
}
