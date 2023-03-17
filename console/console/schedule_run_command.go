package console

import (
	"github.com/oarkflow/framework/contracts/console"
	"github.com/oarkflow/framework/contracts/console/command"
	"github.com/oarkflow/framework/facades"
)

type ScheduleRunCommand struct {
}

// Signature The name and signature of the console command.
func (receiver *ScheduleRunCommand) Signature() string {
	return "schedule:run"
}

// Description The console command description.
func (receiver *ScheduleRunCommand) Description() string {
	return "Run Scheduled Command"
}

// Extend The console command extend.
func (receiver *ScheduleRunCommand) Extend() command.Extend {
	return command.Extend{}
}

// Handle Execute the console command.
func (receiver *ScheduleRunCommand) Handle(ctx console.Context) error {
	facades.Schedule.Run()
	return nil
}
