package cache

import (
	"github.com/oarkflow/framework/contracts/cache"
)

type Application struct {
	Store  cache.Store
	Prefix string
}

func (app *Application) Init() cache.Store {
	if app.Store != nil {
		return app.Store
	}
	return NewMemoryCache(app.Prefix)
}
