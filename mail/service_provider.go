package mail

import (
	"github.com/sujit-baniya/framework/contracts/queue"
	"github.com/sujit-baniya/framework/facades"
)

type ServiceProvider struct {
}

func (route *ServiceProvider) Register() {
	app := Application{}
	facades.Mail = app.Init()
}

func (route *ServiceProvider) Boot() {
	facades.Queue.Register([]queue.Job{
		&SendMailJob{},
	})
}
