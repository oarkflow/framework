package http

import (
	"context"
)

type HandlerFunc func(Context) error

type Context interface {
	context.Context
	WithValue(key string, value interface{})
	Request() Request
	Response() Response
}
