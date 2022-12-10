package schedule

import (
	"context"
	"github.com/sujit-baniya/framework/facades"
)

type ServiceProvider struct {
}

func (receiver *ServiceProvider) Register() {
	facades.Schedule = &Application{ctx: context.Background()}
}

func (receiver *ServiceProvider) Boot() {

}
