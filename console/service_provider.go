package console

import (
	"github.com/sujit-baniya/framework/console/console"
	console2 "github.com/sujit-baniya/framework/contracts/console"
	"github.com/sujit-baniya/framework/facades"
)

type ServiceProvider struct {
	Name string
}

func (receiver *ServiceProvider) Boot() {
	receiver.registerCommands()
}

func (receiver *ServiceProvider) Register() {
	facades.Artisan = NewApplication(receiver.Name)
}

func (receiver *ServiceProvider) registerCommands() {
	facades.Artisan.Register([]console2.Command{
		&console.ListCommand{},
		&console.KeyGenerateCommand{},
		&console.MakeCommand{},
		&console.ScheduleRunCommand{},
	})
}
