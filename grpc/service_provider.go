package grpc

import (
	"github.com/sujit-baniya/framework/facades"
)

type ServiceProvider struct {
}

func (route *ServiceProvider) Register() {
	facades.Grpc = NewApplication()
}

func (route *ServiceProvider) Boot() {

}
