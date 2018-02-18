package webAppGo

import (
	"io/ioutil"
	"strings"
)

// Page is our type for storing webpages in memory
type Page struct {
	Title string
	Body  []byte
}

// SaveToCache is for saving pages to a cached file
func (p Page) SaveToCache() error {
	if strings.Contains(p.Title, " ") {
		p.Title = strings.Replace(p.Title, " ", "_", -1)
	}
	f := "~/go/src/github.com/hamster2020/webAppGo/cache/" + p.Title + ".txt"
	err := ioutil.WriteFile(f, p.Body, 0600)
	if err != nil {
		return err
	}
	return nil
}

// LoadPageFromCache is for loading webpages from a file
func LoadPageFromCache(title string) (*Page, error) {
	f := "~/go/src/github.com/hamster2020/webAppGo/cache/" + title + ".txt"
	body, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}
