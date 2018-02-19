package webAppGo

// Page is our type for storing webpages in memory
type Page struct {
	Title string
	Body  []byte
}

// PageCache is a container for cachers that can save and load cache pages to
// and from files
type PageCache interface {
	SaveToCache(*Page) error
	LoadPageFromCache(string) (*Page, error)
}
