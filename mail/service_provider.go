package mail

import (
	"github.com/sujit-baniya/framework/contracts/mail"
	"github.com/sujit-baniya/framework/contracts/queue"
	"github.com/sujit-baniya/framework/facades"
)

type ServiceProvider struct {
	Mailer mail.Mail
}

func (route *ServiceProvider) Register() {
	app := Application{Mailer: route.Mailer}
	facades.Mail = app.Init()
}

func (route *ServiceProvider) Boot() {
	facades.Queue.Register([]queue.Job{
		&SendMailJob{},
	})
}
