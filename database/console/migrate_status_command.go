package console

import (
	"github.com/gookit/color"

	"github.com/oarkflow/framework/contracts/console"
	"github.com/oarkflow/framework/contracts/console/command"
	"github.com/oarkflow/framework/facades"
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
		Flags: []command.Flag{
			{
				Name:    "connection",
				Value:   "",
				Aliases: []string{"c"},
				Usage:   "connection driver for the database",
			},
			{
				Name:  "dir",
				Value: "",
				Usage: "directory for migration",
			},
		},
	}
}

// Handle Execute the console command.
func (receiver *MigrateStatusCommand) Handle(ctx console.Context) error {
	dir := ctx.Option("dir")
	connection := ctx.Option("connection")
	if connection == "" {
		connection = facades.Config.GetString("database.default")
	}
	m, err := getMigrate(connection, dir)
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
