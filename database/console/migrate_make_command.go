package console

import (
	"github.com/gookit/color"
	"github.com/sujit-baniya/framework/contracts/console"
	"github.com/sujit-baniya/framework/contracts/console/command"
	"github.com/sujit-baniya/framework/facades"
)

type MigrateMakeCommand struct {
}

// Signature The name and signature of the console command.
func (receiver *MigrateMakeCommand) Signature() string {
	return "make:migration"
}

// Description The console command description.
func (receiver *MigrateMakeCommand) Description() string {
	return "Create a new migration file"
}

// Extend The console command extend.
func (receiver *MigrateMakeCommand) Extend() command.Extend {
	return command.Extend{
		Category: "make",
		Flags: []command.Flag{
			{
				Name:    "connection",
				Value:   "",
				Aliases: []string{"c"},
				Usage:   "connection driver for the database",
			},
		},
	}
}

// Handle Execute the console command.
func (receiver *MigrateMakeCommand) Handle(ctx console.Context) error {
	name := ctx.Argument(0)
	if name == "" {
		color.Redln("Not enough arguments (missing: name)")

		return nil
	}
	connection := ctx.Option("connection")
	if connection == "" {
		connection = facades.Config.GetString("database.default")
	}
	m, err := getMigrate(connection)
	if err != nil {
		return err
	}
	if m == nil {
		color.Yellowln("Please fill database config first")

		return nil
	}
	err = m.New(name)
	color.Green.Printf("Created Migration: %s\n", name)
	return err
}
