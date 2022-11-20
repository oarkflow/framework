package route

import (
	"github.com/sujit-baniya/frame/server"
	"github.com/sujit-baniya/framework/facades"
)

type Application struct {
	Engine *server.Frame
}

func (app *Application) Init() *server.Frame {
	if app.Engine != nil {
		return app.Engine
	}
	return NewFrame()
}

func NewFrame() *server.Frame {
	return server.New(
		server.WithHostPorts(facades.Config.GetString("app.host")),
		server.WithRemoveExtraSlash(true),
	)
}
