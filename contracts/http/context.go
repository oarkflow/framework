package http

import (
	"context"
	"net/http"
)

type HandlerFunc func(Context) error

type Context interface {
	context.Context
	WithValue(key string, value interface{})
	Request() Request
	Response() Response
	Header(key, defaultValue string) string
	Headers() http.Header
	Method() string
	Path() string
	Url() string
	FullUrl() string
	Ip() string

	//Input Retrieve  an input item from the request: /users/{id}
	Params(key string) string
	// Query Retrieve a query string item form the request: /users?id=1
	Query(key, defaultValue string) string
	// Form Retrieve a form string item form the post: /users POST:id=1
	Form(key, defaultValue string) string
	Bind(obj interface{}) error
	File(name string) (File, error)

	AbortWithStatus(code int)
	Next() error
}
