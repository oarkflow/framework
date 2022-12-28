package validation

import (
	consolecontract "github.com/sujit-baniya/framework/contracts/console"
	"github.com/sujit-baniya/framework/facades"
	"github.com/sujit-baniya/framework/validation/console"
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
