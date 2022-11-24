package auth

import (
	"github.com/sujit-baniya/framework/auth/console"
	"github.com/sujit-baniya/framework/contracts/auth"
	contractconsole "github.com/sujit-baniya/framework/contracts/console"
	"github.com/sujit-baniya/framework/facades"
)

type ServiceProvider struct {
	Auth auth.Auth
}

func (database *ServiceProvider) Register() {
	if database.Auth != nil {
		facades.Auth = database.Auth
		return
	}
	facades.Auth = GetAuth(facades.Config.GetString("auth.defaults.guard"))
}

func (database *ServiceProvider) Boot() {
	database.registerCommands()
}

func (database *ServiceProvider) registerCommands() {
	facades.Artisan.Register([]contractconsole.Command{
		&console.JwtSecretCommand{},
	})
}

func GetAuth(guard string) auth.Auth {
	driver := facades.Config.GetString("auth.guards." + guard + ".driver")
	switch driver {
	case "jwt":
		return NewJwt(guard)
	default:
		return NewSession(guard)
	}
}
