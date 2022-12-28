package log

import (
	"context"
	"time"
)

const (
	StackDriver  = "stack"
	SingleDriver = "single"
	DailyDriver  = "daily"
	CustomDriver = "custom"
)

const (
	PanicLevel Level = iota
	FatalLevel
	ErrorLevel
	WarningLevel
	InfoLevel
	DebugLevel
)

type Log interface {
	WithContext(ctx context.Context) Writer
	Writer
}

type Writer interface {
	Debug(args ...any)
	Debugf(format string, args ...any)
	Info(args ...any)
	Infof(format string, args ...any)
	Warning(args ...any)
	Warningf(format string, args ...any)
	Error(args ...any)
	Errorf(format string, args ...any)
	Fatal(args ...any)
	Fatalf(format string, args ...any)
	Panic(args ...any)
	Panicf(format string, args ...any)
}

type Logger interface {
	// Handle pass channel config path here
	Handle(channel string) (Hook, error)
}

type Hook interface {
	// Levels monitoring level
	Levels() []Level
	// Fire execute logic when trigger
	Fire(Entry) error
}

type Entry interface {
	Context() context.Context
	Level() Level
	Time() time.Time
	Message() string
	// DEPRECATED: use Level()
	GetLevel() Level
	// DEPRECATED: use Time()
	GetTime() time.Time
	// DEPRECATED: use Message()
	GetMessage() string
}
