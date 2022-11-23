package filesystem

import (
	"github.com/sujit-baniya/framework/facades"
)

type ServiceProvider struct {
}

func (database *ServiceProvider) Register() {
	facades.Storage = NewStorage()
}

func (database *ServiceProvider) Boot() {

}
