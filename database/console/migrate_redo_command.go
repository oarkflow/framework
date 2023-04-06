package console

import (
	"strconv"

	"github.com/gookit/color"

	"github.com/oarkflow/framework/contracts/console"
	"github.com/oarkflow/framework/contracts/console/command"
)

type MigrateRedoCommand struct {
}

// Signature The name and signature of the console command.
func (receiver *MigrateRedoCommand) Signature() string {
	return "migrate:redo"
}

// Description The console command description.
func (receiver *MigrateRedoCommand) Description() string {
	return "Re-apply the last migration"
}

// Extend The console command extend.
func (receiver *MigrateRedoCommand) Extend() command.Extend {
	return command.Extend{
		Category: "migrate",
		Flags: []command.Flag{
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
func (receiver *MigrateRedoCommand) Handle(ctx console.Context) error {
	m, err := getMigrate()
	if err != nil {
		return err
	}
	if m == nil {
		color.Yellowln("Please fill database config first")
		return nil
	}
	dryrunString := ctx.Option("dryrun")
	dryrun, err := strconv.ParseBool(dryrunString)
	if err != nil {
		dryrun = false
	}
	if err := m.Redo(dryrun); err != nil {
		color.Redln("Migration status check failed:", err.Error())
		return nil
	}

	color.Greenln("Migration check success")

	return nil
}
