package config

import (
	"flag"
	"github.com/sujit-baniya/framework/facades"
)

type ServiceProvider struct {
}

func (config *ServiceProvider) Register() {
	if facades.Config == nil {
		var env *string
		env = flag.String("env", ".env", "custom .env path")
		flag.Parse()
		facades.Config = NewApplication(*env)
	}
}

func (config *ServiceProvider) Boot() {

}
