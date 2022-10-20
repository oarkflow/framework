package route

import (
	"github.com/sujit-baniya/framework/facades"
)

type ServiceProvider struct {
}

func (route *ServiceProvider) Register() {
	app := Application{}
	facades.Route = app.Init()
}

func (route *ServiceProvider) Boot() {

}
