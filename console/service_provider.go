package console

import (
	"github.com/oarkflow/framework/console/console"
	console2 "github.com/oarkflow/framework/contracts/console"
	"github.com/oarkflow/framework/facades"
)

type ServiceProvider struct {
	Name string
}

func (receiver *ServiceProvider) Boot() {
	receiver.registerCommands()
}

func (receiver *ServiceProvider) Register() {
	app := Application{Name: receiver.Name}
	facades.Artisan = app.Init()
}

func (receiver *ServiceProvider) registerCommands() {
	facades.Artisan.Register([]console2.Command{
		&console.ListCommand{},
		&console.KeyGenerateCommand{},
		&console.MakeCommand{},
		&console.ScheduleRunCommand{},
	})
}
