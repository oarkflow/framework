package database

import (
	"context"
	"github.com/oarkflow/framework/contracts/database/orm"
	"gorm.io/gorm"
)

type Application struct {
	Config     *gorm.Config
	Name       string
	DisableLog bool
}

func (app *Application) Init() orm.Orm {
	return NewOrm(context.Background(), app.Name, app.Config, app.DisableLog)
}
