package cache

import (
	"github.com/oarkflow/framework/cache/console"
	"github.com/oarkflow/framework/contracts/cache"
	console2 "github.com/oarkflow/framework/contracts/console"
	"github.com/oarkflow/framework/facades"
)

type ServiceProvider struct {
	Store  cache.Store
	Prefix string
}

func (database *ServiceProvider) Register() {
	app := Application{Store: database.Store, Prefix: database.Prefix}
	facades.Cache = app.Init()
}

func (database *ServiceProvider) Boot() {
	database.registerCommands()
}

func (database *ServiceProvider) registerCommands() {
	facades.Artisan.Register([]console2.Command{
		&console.ClearCommand{},
	})
}
