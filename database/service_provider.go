package database

import (
	"gorm.io/gorm"

	consolecontract "github.com/oarkflow/framework/contracts/console"
	"github.com/oarkflow/framework/database/console"
	"github.com/oarkflow/framework/facades"
)

type ServiceProvider struct {
	Config     *gorm.Config
	DisableLog bool
	Name       string
}

func (database *ServiceProvider) Register() {
	app := Application{Config: database.Config, DisableLog: database.DisableLog, Name: database.Name}
	facades.Orm = app.Init()
}

func (database *ServiceProvider) Boot() {
	database.registerCommands()
}

func (database *ServiceProvider) registerCommands() {
	facades.Artisan.Register([]consolecontract.Command{
		&console.MigrateMakeCommand{},
		&console.MigrateCommand{},
		&console.MigrateRollbackCommand{},
		&console.MigrateStatusCommand{},
		&console.MigrateRedoCommand{},
	})
}
