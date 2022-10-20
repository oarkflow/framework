package route

import (
	"github.com/sujit-baniya/framework/contracts/route"
)

type Application struct {
}

func (app *Application) Init() route.Engine {
	return NewGin()
}
