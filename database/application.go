package database

import (
	"context"

	"gorm.io/gorm"

	"github.com/oarkflow/framework/contracts/database/orm"
)

type Application struct {
	Config     *gorm.Config
	Name       string
	DisableLog bool
}

func (app *Application) Init() orm.Orm {
	return NewOrm(context.Background(), app.Name, app.Config, app.DisableLog)
}
