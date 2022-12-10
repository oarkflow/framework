package schedule

import (
	"github.com/gookit/color"
	"github.com/robfig/cron/v3"
	"github.com/sujit-baniya/framework/contracts/schedule"
	"github.com/sujit-baniya/framework/facades"
	"github.com/sujit-baniya/framework/schedule/support"
)

type Application struct {
	cron *cron.Cron
}

func (app *Application) Call(callback func()) schedule.Event {
	return &support.Event{Callback: callback}
}

func (app *Application) Command(command string) schedule.Event {
	return &support.Event{Command: command}
}

func (app *Application) Register(events []schedule.Event) {
	if app.cron == nil {
		app.cron = cron.New(cron.WithLogger(&Logger{}))
	}

	app.addEvents(events)
}

func (app *Application) RegisterOne(event schedule.Event) (cron.EntryID, error) {
	if app.cron == nil {
		app.cron = cron.New(cron.WithLogger(&Logger{}))
	}

	return app.addEvent(event)
}

func (app *Application) Unregister(id int) {
	if app.cron != nil {
		app.cron.Remove(cron.EntryID(id))
	}
}

func (app *Application) Run() {
	app.cron.Start()
}

func (app *Application) addEvent(event schedule.Event) (cron.EntryID, error) {
	chain := cron.NewChain()
	if event.GetDelayIfStillRunning() {
		chain = cron.NewChain(cron.DelayIfStillRunning(&Logger{}))
	} else if event.GetSkipIfStillRunning() {
		chain = cron.NewChain(cron.SkipIfStillRunning(&Logger{}))
	}
	return app.cron.AddJob(event.GetCron(), chain.Then(app.getJob(event)))
}

func (app *Application) addEvents(events []schedule.Event) {
	for _, event := range events {
		_, err := app.addEvent(event)
		if err != nil {
			facades.Log.Errorf("add schedule error: %v", err)
		}
	}
}

func (app *Application) getJob(event schedule.Event) cron.Job {
	return cron.FuncJob(func() {
		if event.GetCommand() != "" {
			facades.Artisan.Call(event.GetCommand())
		} else {
			event.GetCallback()()
		}
	})
}

type Logger struct{}

func (log *Logger) Info(msg string, keysAndValues ...any) {
	color.Green.Printf("%s %v\n", msg, keysAndValues)
}

func (log *Logger) Error(err error, msg string, keysAndValues ...any) {
	facades.Log.Error(msg, keysAndValues)
}
