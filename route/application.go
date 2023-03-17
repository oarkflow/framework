package route

import (
	"github.com/oarkflow/frame/server"
	"github.com/oarkflow/framework/facades"
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
	template := facades.Config.GetString("view.template")
	extension := facades.Config.GetString("view.extension")
	h := server.Default(
		server.WithHostPorts(facades.Config.GetString("app.host")),
		server.WithRemoveExtraSlash(true),
		server.WithRedirectTrailingSlash(true),
	)
	h.Name = facades.Config.GetString("app.name")
	h.SetHTMLTemplate(template, extension)
	return h
}
