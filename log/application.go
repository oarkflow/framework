package log

import (
	"github.com/sujit-baniya/framework/contracts/log"
)

type Application struct {
}

func (app *Application) Init() log.Log {
	return NewLogrus()
}
