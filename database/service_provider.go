package database

import (
	consolecontract "github.com/oarkflow/framework/contracts/console"
	"github.com/oarkflow/framework/database/console"
	"github.com/oarkflow/framework/facades"
	"gorm.io/gorm"
)

type ServiceProvider struct {
	Config     *gorm.Config
	DisableLog bool
}

func (database *ServiceProvider) Register() {
	app := Application{Config: database.Config, DisableLog: database.DisableLog}
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
