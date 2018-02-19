package cache

// Cacher is a data type that can load and save cached pages
type Cacher struct {
	//	SaveToCacheFunc       func() error
	//	LoadPageFromCacheFunc func(string) (*webAppGo.Page, error)
}

// NewCache is a constructor for a new cache
func NewCache() *Cacher {
	return &Cacher{}
}
