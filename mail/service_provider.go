package mail

import (
	"github.com/sujit-baniya/framework/contracts/queue"
	"github.com/sujit-baniya/framework/facades"
)

type ServiceProvider struct {
	Mailer *Mailer
}

func (route *ServiceProvider) Register() {
	app := Application{Mailer: route.Mailer}
	DefaultMailer = app.Init()
}

func (route *ServiceProvider) Boot() {
	facades.Queue.Register([]queue.Job{
		&SendMailJob{},
	})
}
