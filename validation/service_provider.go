package validation

import (
	consolecontract "github.com/oarkflow/framework/contracts/console"
	"github.com/oarkflow/framework/facades"
	"github.com/oarkflow/framework/validation/console"
)

type ServiceProvider struct {
}

func (database *ServiceProvider) Register() {
	facades.Validation = NewValidation()
}

func (database *ServiceProvider) Boot() {
	database.registerCommands()
}

func (database *ServiceProvider) registerCommands() {
	facades.Artisan.Register([]consolecontract.Command{
		&console.RuleMakeCommand{},
	})
}
