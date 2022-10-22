package auth

import (
	"github.com/sujit-baniya/framework/auth/console"
	contractconsole "github.com/sujit-baniya/framework/contracts/console"
	"github.com/sujit-baniya/framework/facades"
)

type ServiceProvider struct {
}

func (database *ServiceProvider) Register() {
	facades.Auth = NewApplication(facades.Config.GetString("auth.defaults.guard"))
}

func (database *ServiceProvider) Boot() {
	database.registerCommands()
}

func (database *ServiceProvider) registerCommands() {
	facades.Artisan.Register([]contractconsole.Command{
		&console.JwtSecretCommand{},
	})
}
