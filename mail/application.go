package mail

import (
	"github.com/sujit-baniya/framework/contracts/mail"
)

type Application struct {
	Mailer mail.Mail
}

func (app *Application) Init() mail.Mail {
	if app.Mailer != nil {
		return app.Mailer
	}
	return NewEmail()
}
