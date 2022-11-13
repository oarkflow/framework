package queue

import (
	"github.com/sujit-baniya/framework/contracts/console"
	"github.com/sujit-baniya/framework/facades"
	queueConsole "github.com/sujit-baniya/framework/queue/console"
)

type ServiceProvider struct {
}

func (receiver *ServiceProvider) Register() {
	facades.Queue = &Application{}
}

func (receiver *ServiceProvider) Boot() {
	receiver.registerCommands()
}

func (receiver *ServiceProvider) registerCommands() {
	facades.Artisan.Register([]console.Command{
		&queueConsole.JobMakeCommand{},
		&queueConsole.QueueWorkCommand{},
	})
}
