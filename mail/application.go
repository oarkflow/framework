package mail

import (
	"github.com/sujit-baniya/framework/contracts/mail"
)

type Application struct {
}

func (app *Application) Init() mail.Mail {
	return NewEmail()
}
