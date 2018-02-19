package webAppGo

// Cachestore is a container for cachers that can save and load cache pages to
// and from files
type Cachestore interface {
	SaveToCache(*Page) error
	LoadPageFromCache(string) (*Page, error)
}
