package console

import (
	"strconv"

	"github.com/oarkflow/framework/contracts/console"
	"github.com/oarkflow/framework/contracts/console/command"
	queue2 "github.com/oarkflow/framework/contracts/queue"
	"github.com/oarkflow/framework/facades"
)

type QueueWorkCommand struct {
}

// Signature The name and signature of the console command.
func (receiver *QueueWorkCommand) Signature() string {
	return "queue:work"
}

// Description The console command description.
func (receiver *QueueWorkCommand) Description() string {
	return "Run background queue server"
}

// Extend The console command extend.
func (receiver *QueueWorkCommand) Extend() command.Extend {
	return command.Extend{
		Flags: []command.Flag{
			{
				Name:    "connection",
				Value:   "redis",
				Aliases: []string{"c"},
				Usage:   "connection driver for the queue server",
			},
			{
				Name:    "queue",
				Value:   "default",
				Aliases: []string{"q"},
				Usage:   "queue to be processed",
			},
			{
				Name:    "concurrent",
				Value:   "10",
				Aliases: []string{"p"},
				Usage:   "Concurrency to run queue server",
			},
		},
	}
}

// Handle Execute the console command.
func (receiver *QueueWorkCommand) Handle(ctx console.Context) error {
	concurrent := 10
	connection := ctx.Option("connection")
	queue := ctx.Option("queue")
	c := ctx.Option("concurrent")
	c1, err := strconv.Atoi(c)
	if err == nil {
		concurrent = c1
	}
	return facades.Queue.Worker(&queue2.Args{
		Connection: connection,
		Queue:      queue,
		Concurrent: concurrent,
	}).Run()
}
