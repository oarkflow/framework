package console

import (
	"github.com/gookit/color"
	"github.com/sujit-baniya/framework/contracts/console"
	"github.com/sujit-baniya/framework/contracts/console/command"
	"github.com/sujit-baniya/framework/database/support"
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
	// It's possible for the developer to specify the tables to modify in this
	// schema operation. The developer may also specify if this table needs
	// to be freshly created, so we can create the appropriate migrations.
	name := ctx.Argument(0)
	if name == "" {
		color.Redln("Not enough arguments (missing: name)")

		return nil
	}
	connection := ctx.Option("connection")
	//Write the migration file to disk.
	MigrateCreator{driver: getDriver(connection)}.Create(name)

	color.Green.Printf("Created Migration: %s\n", name)

	return nil
}

func getDriver(connection string) string {
	if connection == "" {
		connection = facades.Config.GetString("database.default")
	}
	driver := facades.Config.GetString("database.connections." + connection + ".driver")
	switch driver {
	case support.Postgresql:
		return "postgres"
	case support.Sqlite:
		return "sqlite3"
	case support.Sqlserver:
		return "sqlserver"
	default:
		return "mysql"
	}
}
