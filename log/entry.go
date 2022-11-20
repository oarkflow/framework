package log

import (
	"context"
	"github.com/sujit-baniya/framework/contracts/log"
	"time"
)

type Entry struct {
	ctx     context.Context
	level   log.Level
	time    time.Time
	message string
}

func (r *Entry) Context() context.Context {
	return r.ctx
}

func (r *Entry) Level() log.Level {
	return r.level
}

func (r *Entry) Time() time.Time {
	return r.time
}

func (r *Entry) Message() string {
	return r.message
}

// DEPRECATED: use Level()
func (r *Entry) GetLevel() log.Level {
	return r.Level()
}

// DEPRECATED: use Level()
func (r *Entry) GetTime() time.Time {
	return r.Time()
}

// DEPRECATED: use Level()
func (r *Entry) GetMessage() string {
	return r.Message()
}
