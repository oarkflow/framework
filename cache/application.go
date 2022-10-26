package cache

import (
	"github.com/sujit-baniya/framework/contracts/cache"
	"github.com/sujit-baniya/framework/facades"
)

type Application struct {
	Store  cache.Store
	Prefix string
}

func (app *Application) Init() cache.Store {
	if app.Store != nil {
		return app.Store
	}
	defaultStore := facades.Config.GetString("cache.default")
	driver := facades.Config.GetString("cache.stores." + defaultStore + ".driver")
	switch driver {
	case "redis":
		return NewRedisCache()
	default:
		return NewMemoryCache(app.Prefix)
	}
}
