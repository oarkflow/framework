package schedule

import (
	"github.com/sujit-baniya/framework/facades"
)

type ServiceProvider struct {
}

func (receiver *ServiceProvider) Register() {
	facades.Schedule = &Application{}
}

func (receiver *ServiceProvider) Boot() {

}
