package cache

import (
	"github.com/oarkflow/framework/cache/console"
	console2 "github.com/oarkflow/framework/contracts/console"
	"github.com/oarkflow/framework/facades"
)

type ServiceProvider struct {
}

func (database *ServiceProvider) Register() {
	app := Application{}
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
