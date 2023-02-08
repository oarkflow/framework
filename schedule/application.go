package schedule

import (
	"context"
	"github.com/gookit/color"
	"github.com/sujit-baniya/cron"
	"github.com/sujit-baniya/framework/contracts/schedule"
	"github.com/sujit-baniya/framework/facades"
	"github.com/sujit-baniya/framework/schedule/support"
	logg "github.com/sujit-baniya/log"
)

type Application struct {
	cron *cron.Cron
	ctx  context.Context
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

func (app *Application) PauseEntry(id int) {
	if app.cron != nil {
		app.cron.PauseEntry(cron.EntryID(id))
	}
}

func (app *Application) StartEntry(id int) {
	if app.cron != nil {
		app.cron.StartEntry(cron.EntryID(id))
	}
}

func (app *Application) Logs(id int) []string {
	if app.cron != nil {
		entry := app.cron.Entry(cron.EntryID(id))
		return entry.Logs
	}
	return nil
}

func (app *Application) Run() {
	app.cron.Start(app.ctx)
}

func (app *Application) addEvent(event schedule.Event) (cron.EntryID, error) {
	chain := cron.NewChain()
	if event.GetDelayIfStillRunning() {
		chain = cron.NewChain(cron.DelayIfStillRunning(&Logger{}))
	} else if event.GetSkipIfStillRunning() {
		chain = cron.NewChain(cron.SkipIfStillRunning(&Logger{}))
	}
	return app.cron.AddJob(event.GetTitle(), event.GetCron(), chain.Then(app.getJob(event)))
}

func (app *Application) addEvents(events []schedule.Event) {
	for _, event := range events {
		_, err := app.addEvent(event)
		if err != nil {
			logg.Error().Err(err).Msg("add schedule error")
		}
	}
}

func (app *Application) getJob(event schedule.Event) cron.Job {
	return cron.FuncJob(func(ctx context.Context) error {
		if event.GetCommand() != "" {
			facades.Artisan.Call(event.GetCommand())
		} else {
			event.GetCallback()()
		}
		return nil
	})
}

type Logger struct{}

func (log *Logger) Info(msg string, keysAndValues ...any) {
	color.Green.Printf("%s %v\n", msg, keysAndValues)
}

func (log *Logger) Error(err error, msg string, keysAndValues ...any) {

	logg.Error().Err(err).KeysAndValues(keysAndValues).Msg(msg)
}
