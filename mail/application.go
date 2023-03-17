package mail

import "github.com/oarkflow/framework/facades"

type Application struct {
	Mailer *Mailer
}

func (app *Application) Init() *Mailer {
	if app.Mailer != nil {
		return app.Mailer
	}
	config := GetMailConfig()
	return New(config, facades.Route.HtmlEngine)
}
