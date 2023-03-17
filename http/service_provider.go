package http

import (
	consolecontract "github.com/oarkflow/framework/contracts/console"
	"github.com/oarkflow/framework/facades"
	"github.com/oarkflow/framework/http/console"
)

type ServiceProvider struct {
}

func (database *ServiceProvider) Register() {

}

func (database *ServiceProvider) Boot() {
	database.registerCommands()
}

func (database *ServiceProvider) registerCommands() {
	facades.Artisan.Register([]consolecontract.Command{
		&console.RequestMakeCommand{},
	})
}
