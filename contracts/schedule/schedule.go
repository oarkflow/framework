package schedule

import (
	"github.com/oarkflow/cron"
)

type Schedule interface {
	//Call Add a new callback event to the schedule.
	Call(callback func()) Event

	//Command Add a new Artisan command event to the schedule.
	Command(command string) Event

	//Register schedules.
	Register(events []Event)

	//RegisterOne schedules.
	RegisterOne(event Event) (cron.EntryID, error)

	//Unregister schedules.
	Unregister(id int)

	//PauseEntry schedules.
	PauseEntry(id int)

	//StartEntry schedules.
	StartEntry(id int)

	Logs(id int) []string

	//Run schedules.
	Run()
}
