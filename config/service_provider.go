package config

import (
	"github.com/sujit-baniya/framework/facades"
)

type ServiceProvider struct {
	EnvPath string
}

func (config *ServiceProvider) Register() {
	if config.EnvPath == "" {
		config.EnvPath = ".env"
	}
	app := Application{EnvPath: config.EnvPath}
	if facades.Config == nil {
		facades.Config = app.Init()
	}
}

func (config *ServiceProvider) Boot() {

}
