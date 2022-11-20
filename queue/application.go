package queue

import (
	queue2 "github.com/sujit-baniya/framework/contracts/queue"
	"github.com/sujit-baniya/framework/queue/support"
)

type Application struct {
	jobs []queue2.Job
}

func (app *Application) Worker(args *queue2.Args) queue2.Worker {
	if args == nil {
		return &support.Worker{}
	}

	return &support.Worker{
		Connection: args.Connection,
		Queue:      args.Queue,
		Concurrent: args.Concurrent,
	}
}

func (app *Application) Register(jobs []queue2.Job) {
	app.jobs = append(app.jobs, jobs...)
}

func (app *Application) GetJobs() []queue2.Job {
	return app.jobs
}

func (app *Application) Job(job queue2.Job, args []queue2.Arg) queue2.Task {
	return &support.Task{
		Job:  job,
		Args: args,
	}
}

func (app *Application) Chain(jobs []queue2.Jobs) queue2.Task {
	return &support.Task{
		Jobs:  jobs,
		Chain: true,
	}
}
