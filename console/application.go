package console

import (
	"github.com/sujit-baniya/framework/contracts/console"
)

type Application struct {
}

func (app *Application) Init() console.Artisan {
	return NewCli()
}
