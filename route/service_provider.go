package route

import (
	"github.com/oarkflow/frame/server"

	"github.com/oarkflow/framework/facades"
)

type ServiceProvider struct {
	Engine *server.Frame
}

func (route *ServiceProvider) Register() {
	app := Application{Engine: route.Engine}
	facades.Route = app.Init()
}

func (route *ServiceProvider) Boot() {

}
