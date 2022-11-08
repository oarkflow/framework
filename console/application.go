package console

import (
	"github.com/sujit-baniya/framework/contracts/console"
)

type Application struct {
	Name string
}

func (app *Application) Init() console.Artisan {
	return NewCli(app.Name)
}
