package console

import (
	"strconv"

	"github.com/gookit/color"

	"github.com/oarkflow/framework/contracts/console"
	"github.com/oarkflow/framework/contracts/console/command"
	"github.com/oarkflow/framework/facades"
)

type MigrateCommand struct {
}

// Signature The name and signature of the console command.
func (receiver *MigrateCommand) Signature() string {
	return "migrate"
}

// Description The console command description.
func (receiver *MigrateCommand) Description() string {
	return "Run the database migrations"
}

// Extend The console command extend.
func (receiver *MigrateCommand) Extend() command.Extend {
	return command.Extend{
		Category: "migrate",
		Flags: []command.Flag{
			{
				Name:    "steps",
				Value:   "0",
				Aliases: []string{"s"},
				Usage:   "Limit sql files for migration",
			},
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
			{
				Name:    "dryrun",
				Value:   "false",
				Aliases: []string{"d"},
				Usage:   "Do not actually execute the query and just print the query",
			},
		},
	}
}

// Handle Execute the console command.
func (receiver *MigrateCommand) Handle(ctx console.Context) error {
	l := ctx.Option("steps")
	d := ctx.Option("dryrun")
	dir := ctx.Option("dir")
	limit, err := strconv.Atoi(l)
	if err != nil {
		limit = 0
	}
	dryrun, err := strconv.ParseBool(d)
	if err != nil {
		dryrun = false
	}
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

	if err := m.Up(limit, dryrun); err != nil {
		color.Redln("Migration failed:", err.Error())

		return nil
	}

	color.Greenln("Migration success")

	return nil
}
