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
	facades.Auth = NewAuth(facades.Config.GetString("auth.defaults.guard"))
}

func (database *ServiceProvider) Boot() {
	database.registerCommands()
}

func (database *ServiceProvider) registerCommands() {
	facades.Artisan.Register([]contractconsole.Command{
		&console.JwtSecretCommand{},
	})
}

func NewAuth(guard string) auth.Auth {
	driver := facades.Config.GetString("auth.guards." + guard + ".driver")
	switch driver {
	case "jwt":
		if facades.JwtAuth == nil {
			jwtAuth := NewJwt(guard)
			facades.JwtAuth = jwtAuth
		}
		return facades.JwtAuth
	case "key":
		if facades.ApiKeyAuth == nil {
			apiKeyAuth := NewApiKey(guard)
			facades.ApiKeyAuth = apiKeyAuth
		}
		return facades.ApiKeyAuth
	default:
		if facades.SessionAuth == nil {
			sessionAuth := NewSession(guard)
			facades.SessionAuth = sessionAuth
		}
		return facades.SessionAuth
	}
}
