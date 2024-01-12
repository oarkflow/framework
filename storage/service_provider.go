package storage

import (
	"github.com/oarkflow/framework/facades"
)

type ServiceProvider struct {
}

func (database *ServiceProvider) Register() {
	facades.Storage = NewStorage("")
}

func (database *ServiceProvider) Boot() {

}
