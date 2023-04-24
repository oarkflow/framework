package auth

import (
	"context"
	"sync"

	"github.com/oarkflow/frame"
	"github.com/oarkflow/frame/middlewares/server/session"

	"github.com/oarkflow/framework/contracts/auth"

	"github.com/oarkflow/framework/auth/access"
	"github.com/oarkflow/framework/auth/console"
	contractconsole "github.com/oarkflow/framework/contracts/console"
	"github.com/oarkflow/framework/facades"
)

type ServiceProvider struct {
	Auth   auth.Auth
	Config session.Config
}

func (database *ServiceProvider) Register() {
	facades.Auth = NewJwt(facades.Config.GetString("auth.defaults.guard"))
	session.Default(database.Config)
	facades.Gate = access.NewGate(context.Background())
}

func (database *ServiceProvider) Boot() {
	database.registerCommands()
}

func (database *ServiceProvider) registerCommands() {
	facades.Artisan.Register([]contractconsole.Command{
		&console.JwtSecretCommand{},
		&console.PolicyMakeCommand{},
	})
}

func GetAuth(guard string) auth.Auth {
	a := Drivers.Get(guard)
	if a != nil {
		return a
	}
	return Drivers.Get("session")
}

type drivers struct {
	driver map[string]auth.Auth
	mu     *sync.RWMutex
}

func (d *drivers) Get(guard string) auth.Auth {
	return d.driver[guard]
}

func (d *drivers) Add(guard string, auth2 auth.Auth) {
	d.mu.Lock()
	d.driver[guard] = auth2
	d.mu.Unlock()
}

func (d *drivers) Remove(guard string) {
	d.mu.Lock()
	delete(d.driver, guard)
	d.mu.Unlock()
}

var Drivers *drivers

func init() {
	var store *session.Store
	if facades.Session != nil {
		store = facades.Session
	} else {
		store = session.DefaultStore
	}
	Drivers = &drivers{
		driver: map[string]auth.Auth{
			"session": NewSession("session", store),
			"jwt":     NewJwt("jwt"),
		},
		mu: &sync.RWMutex{},
	}
}

func Logout(ctx *frame.Context) {
	for _, a := range Drivers.driver {
		if a != nil {
			a.Logout(ctx)
		}
	}
}
