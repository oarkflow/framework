package database

import (
	"context"
	"github.com/oarkflow/framework/contracts/database/orm"
	"gorm.io/gorm"
)

type Application struct {
	Config     *gorm.Config
	DisableLog bool
}

func (app *Application) Init() orm.Orm {
	return NewOrm(context.Background(), app.Config, app.DisableLog)
}
