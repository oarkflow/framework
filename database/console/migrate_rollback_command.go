package console

import (
	"strconv"

	"github.com/oarkflow/framework/contracts/console"
	"github.com/oarkflow/framework/contracts/console/command"

	"github.com/gookit/color"
)

type MigrateRollbackCommand struct {
}

// Signature The name and signature of the console command.
func (receiver *MigrateRollbackCommand) Signature() string {
	return "migrate:rollback"
}

// Description The console command description.
func (receiver *MigrateRollbackCommand) Description() string {
	return "Rollback the database migrations"
}

// Extend The console command extend.
func (receiver *MigrateRollbackCommand) Extend() command.Extend {
	return command.Extend{
		Category: "migrate",
		Flags: []command.Flag{
			{
				Name:    "steps",
				Value:   "0",
				Aliases: []string{"s"},
				Usage:   "rollback steps",
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
func (receiver *MigrateRollbackCommand) Handle(ctx console.Context) error {
	m, err := getMigrate()
	if err != nil {
		return err
	}
	if m == nil {
		color.Yellowln("Please fill database config first")
		return nil
	}

	stepString := ctx.Option("steps")
	dryrunString := ctx.Option("dryrun")
	step, err := strconv.Atoi(stepString)
	if err != nil {
		step = 0
	}
	dryrun, err := strconv.ParseBool(dryrunString)
	if err != nil {
		dryrun = false
	}
	if err := m.Down(step, dryrun); err != nil {
		color.Redln("Migration failed:", err.Error())
		return nil
	}

	color.Greenln("Migration rollback success")

	return nil
}
