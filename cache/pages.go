package cache

import (
	"io/ioutil"
	"strings"

	"github.com/hamster2020/webAppGo"
)

// SaveToCache is for saving pages to a cached file
func (cache Cacher) SaveToCache(p *webAppGo.Page) error {
	if strings.Contains(p.Title, " ") {
		p.Title = strings.Replace(p.Title, " ", "_", -1)
	}
	f := "../../cache/" + p.Title + ".txt"
	err := ioutil.WriteFile(f, p.Body, 0600)
	if err != nil {
		return err
	}
	return nil
}

// LoadPageFromCache is for loading webpages from a file
func (cache Cacher) LoadPageFromCache(title string) (*webAppGo.Page, error) {
	f := "../../cache/" + title + ".txt"
	body, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}
	return &webAppGo.Page{Title: title, Body: body}, nil
}
