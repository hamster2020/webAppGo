package sqlite

import (
	"strconv"
	"strings"
	"time"

	"github.com/hamster2020/webAppGo"
)

// SavePage is for saving pages to the db
func (db *DB) SavePage(p *webAppGo.Page) error {
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

// LoadPage is for loading pages from the db
func (db *DB) LoadPage(title string) (*webAppGo.Page, error) {
	var name string
	var body []byte
	rows, err := db.Query(selectPageFromTable, title)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		rows.Scan(&name, &body)
	}
	return &webAppGo.Page{Title: name, Body: body}, nil
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
