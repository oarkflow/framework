package log

import "github.com/sujit-baniya/framework/facades"

type ServiceProvider struct {
}

func (log *ServiceProvider) Register() {
	app := Application{}
	facades.Log = app.Init()
}

func (log *ServiceProvider) Boot() {

}
