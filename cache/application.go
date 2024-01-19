package cache

import (
	"github.com/oarkflow/framework/contracts/cache"
)

type Application struct {
}

func (app *Application) Init() cache.Store {
	defaultCache = NewCache("")
	return defaultCache
}
