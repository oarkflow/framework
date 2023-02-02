package console

import (
	"github.com/gookit/color"
	"github.com/sujit-baniya/framework/contracts/console"
	"github.com/sujit-baniya/framework/contracts/console/command"
)

type MigrateStatusCommand struct {
}

// Signature The name and signature of the console command.
func (receiver *MigrateStatusCommand) Signature() string {
	return "migrate:status"
}

// Description The console command description.
func (receiver *MigrateStatusCommand) Description() string {
	return "Status of the database migrations"
}

// Extend The console command extend.
func (receiver *MigrateStatusCommand) Extend() command.Extend {
	return command.Extend{
		Category: "migrate",
	}
}

// Handle Execute the console command.
func (receiver *MigrateStatusCommand) Handle(ctx console.Context) error {
	m, err := getMigrate()
	if err != nil {
		return err
	}
	if m == nil {
		color.Yellowln("Please fill database config first")
		return nil
	}
	if err := m.Status(); err != nil {
		color.Redln("Migration status check failed:", err.Error())
		return nil
	}

	color.Greenln("Migration check success")

	return nil
}
