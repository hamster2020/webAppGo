package models

import (
	"io/ioutil"
	"strconv"
	"strings"
	"time"
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
	f := "cache/" + p.Title + ".txt"
	err := ioutil.WriteFile(f, p.Body, 0600)
	if err != nil {
		return err
	}
	return nil
}

// SavePage is for saving pages to the db
func (db *DB) SavePage(p *Page) error {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	if strings.Contains(p.Title, " ") {
		p.Title = strings.Replace(p.Title, " ", "_", -1)
	}
	_, err := db.Exec(createPagesTable)
	if err != nil {
		return err
	}
	_, err = db.Exec(insertIntoPagesTable, p.Title, p.Body, timestamp)
	if err != nil {
		return err
	}
	return nil
}

// LoadPageFromCache is for loading webpages from a file
func LoadPageFromCache(title string) (*Page, error) {
	f := "cache/" + title + ".txt"
	body, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

// LoadPage is for loading pages from the db
func (db *DB) LoadPage(title string) (*Page, error) {
	var name string
	var body []byte
	rows, err := db.Query(selectPageFromTable, title)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		rows.Scan(&name, &body)
	}
	return &Page{Title: name, Body: body}, nil
}

// PageExists checks to see if a page exists in the db already
func (db *DB) PageExists(title string) (bool, error) {
	var pt string
	var pb []byte
	rows, err := db.Query(selectTitleBodyFromTable, title)
	if err != nil {
		return false, err
	}
	for rows.Next() {
		rows.Scan(&pt, &pb)
	}
	if pt != "" && pb != nil {
		return true, nil
	}
	return false, nil
}
