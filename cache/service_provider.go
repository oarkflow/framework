package cache

import (
	"github.com/sujit-baniya/framework/cache/console"
	"github.com/sujit-baniya/framework/contracts/cache"
	console2 "github.com/sujit-baniya/framework/contracts/console"
	"github.com/sujit-baniya/framework/facades"
)

type ServiceProvider struct {
	Store  cache.Store
	Prefix string
}

func (database *ServiceProvider) Register() {
	app := Application{Store: database.Store}
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
