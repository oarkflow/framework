package auth

import (
	"github.com/sujit-baniya/frame"
	"github.com/sujit-baniya/framework/auth/console"
	"github.com/sujit-baniya/framework/contracts/auth"
	contractconsole "github.com/sujit-baniya/framework/contracts/console"
	"github.com/sujit-baniya/framework/facades"
	"sync"
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
	Drivers = &drivers{
		driver: map[string]auth.Auth{
			"session": NewSession("session"),
			"jwt":     NewJwt("web"),
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