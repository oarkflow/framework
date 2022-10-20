package mock

import (
	cachemocks "github.com/sujit-baniya/framework/contracts/cache/mocks"
	configmocks "github.com/sujit-baniya/framework/contracts/config/mocks"
	consolemocks "github.com/sujit-baniya/framework/contracts/console/mocks"
	ormmocks "github.com/sujit-baniya/framework/contracts/database/orm/mocks"
	eventmocks "github.com/sujit-baniya/framework/contracts/event/mocks"
	mailmocks "github.com/sujit-baniya/framework/contracts/mail/mocks"
	queuemocks "github.com/sujit-baniya/framework/contracts/queue/mocks"
	"github.com/sujit-baniya/framework/facades"
	"github.com/sujit-baniya/framework/log"
)

func Cache() *cachemocks.Store {
	mockCache := &cachemocks.Store{}
	facades.Cache = mockCache

	return mockCache
}

func Config() *configmocks.Config {
	mockConfig := &configmocks.Config{}
	facades.Config = mockConfig

	return mockConfig
}

func Artisan() *consolemocks.Artisan {
	mockArtisan := &consolemocks.Artisan{}
	facades.Artisan = mockArtisan

	return mockArtisan
}

func Orm() (*ormmocks.Orm, *ormmocks.DB, *ormmocks.Transaction) {
	mockOrm := &ormmocks.Orm{}
	facades.Orm = mockOrm

	return mockOrm, &ormmocks.DB{}, &ormmocks.Transaction{}
}

func Event() (*eventmocks.Instance, *eventmocks.Task) {
	mockEvent := &eventmocks.Instance{}
	facades.Event = mockEvent

	return mockEvent, &eventmocks.Task{}
}

func Log() {
	app := log.Application{}
	facades.Log = app.Init()
	facades.Log.Testing(true)
}

func Mail() *mailmocks.Mail {
	mockMail := &mailmocks.Mail{}
	facades.Mail = mockMail

	return mockMail
}

func Queue() (*queuemocks.Queue, *queuemocks.Task) {
	mockQueue := &queuemocks.Queue{}
	facades.Queue = mockQueue

	return mockQueue, &queuemocks.Task{}
}
