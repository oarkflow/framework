package queue

type Task interface {
	Dispatch() error
	DispatchSync() error
	OnConnection(connection string) Task
	OnQueue(queue string) Task
}
