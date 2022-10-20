package route

import (
	"github.com/sujit-baniya/framework/contracts/route"
)

type Application struct {
	Engine route.Engine
}

func (app *Application) Init() route.Engine {
	if app.Engine != nil {
		return app.Engine
	}
	return NewGin()
}
