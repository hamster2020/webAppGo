package cache

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/hamster2020/webAppGo"
)

// Cacher is a data type that can load and save cached pages
type cache struct {
	Path string
}

// NewCache is a constructor for a new cache
func NewCache(cachePath string) *cache {
	return &cache{Path: cachePath}
}

// SaveToCache is for saving pages to a cached file
func (c *cache) SaveToCache(p *webAppGo.Page) error {
	if strings.Contains(p.Title, " ") {
		p.Title = strings.Replace(p.Title, " ", "_", -1)
	}
	f := c.Path + p.Title + ".txt"
	err := ioutil.WriteFile(f, p.Body, 0600)
	if err != nil {
		return err
	}
	return nil
}

// LoadPageFromCache is for loading webpages from a file
func (c *cache) LoadPageFromCache(title string) (*webAppGo.Page, error) {
	f := c.Path + title + ".txt"
	body, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}
	return &webAppGo.Page{Title: title, Body: body}, nil
}

// DeletePageFromCache is for loading webpages from a file
func (c *cache) DeletePageFromCache(title string) error {
	f := c.Path + title + ".txt"
	err := os.Remove(f)
	if err != nil {
		return err
	}
	return nil
}
