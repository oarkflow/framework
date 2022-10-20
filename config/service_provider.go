package config

import (
	"github.com/sujit-baniya/framework/facades"
)

type ServiceProvider struct {
}

func (config *ServiceProvider) Register() {
	app := Application{}
	facades.Config = app.Init()
}

func (config *ServiceProvider) Boot() {

}
