package config

import (
	"github.com/oarkflow/framework/facades"
)

type ServiceProvider struct {
}

func (config *ServiceProvider) Register() {
	app := Application{}
	if facades.Config == nil {
		facades.Config = app.Init()
	}
}

func (config *ServiceProvider) Boot() {

}
