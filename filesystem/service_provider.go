package filesystem

import (
	"github.com/oarkflow/framework/facades"
)

type ServiceProvider struct {
}

func (database *ServiceProvider) Register() {
	facades.Filesystem = NewStorage()
}

func (database *ServiceProvider) Boot() {

}
