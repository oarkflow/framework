package console

import (
	"github.com/oarkflow/framework/contracts/console"
)

type Application struct {
	Name string
}

func (app *Application) Init() console.Artisan {
	return NewCli(app.Name)
}
