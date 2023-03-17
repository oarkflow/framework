package queue

import (
	"github.com/oarkflow/framework/contracts/console"
	"github.com/oarkflow/framework/facades"
	queueConsole "github.com/oarkflow/framework/queue/console"
)

type ServiceProvider struct {
}

func (receiver *ServiceProvider) Register() {
	facades.Queue = NewApplication()
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
