package support

import (
	"errors"
	queue2 "github.com/sujit-baniya/framework/contracts/queue"

	"github.com/sujit-baniya/machinery"
	"github.com/sujit-baniya/machinery/tasks"
)

type Task struct {
	Job        queue2.Job
	Jobs       []queue2.Jobs
	Chain      bool
	Args       []queue2.Arg
	connection string
	queue      string
	server     *machinery.Server
}

func (receiver *Task) Dispatch() error {
	driver := getDriver(receiver.connection)

	if driver == "" {
		return errors.New("unknown queue driver")
	}
	if driver == DriverSync || driver == "" {
		return receiver.DispatchSync()
	}

	server, err := GetServer(receiver.connection, receiver.queue)
	if err != nil {
		return err
	}
	receiver.server = server

	if receiver.Chain {
		for _, job := range receiver.Jobs {
			if err := receiver.handleAsync(job.Job, job.Args); err != nil {
				return err
			}
		}

		return nil
	} else {
		return receiver.handleAsync(receiver.Job, receiver.Args)
	}
}

func (receiver *Task) DispatchSync() error {
	if receiver.Chain {
		for _, job := range receiver.Jobs {
			if err := receiver.handleSync(job.Job, job.Args); err != nil {
				return err
			}
		}

		return nil
	} else {
		return receiver.handleSync(receiver.Job, receiver.Args)
	}
}

func (receiver *Task) handleSync(job queue2.Job, args []queue2.Arg) error {
	var realArgs []any
	for _, arg := range args {
		realArgs = append(realArgs, arg.Value)
	}

	return job.Handle(realArgs...)
}

func (receiver *Task) handleAsync(job queue2.Job, args []queue2.Arg) error {
	var realArgs []tasks.Arg
	for _, arg := range args {
		realArgs = append(realArgs, tasks.Arg{
			Type:  arg.Type,
			Value: arg.Value,
		})
	}

	_, err := receiver.server.SendTask(&tasks.Signature{
		Name: job.Signature(),
		Args: realArgs,
	})
	if err != nil {
		return err
	}

	return nil
}

func (receiver *Task) OnConnection(connection string) queue2.Task {
	receiver.connection = connection

	return receiver
}

func (receiver *Task) OnQueue(queue string) queue2.Task {
	receiver.queue = queue

	return receiver
}
