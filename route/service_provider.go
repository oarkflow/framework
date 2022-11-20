package route

import (
	"github.com/sujit-baniya/frame/server"
	"github.com/sujit-baniya/framework/facades"
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
