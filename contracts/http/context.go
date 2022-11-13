package http

import (
	"context"
	"net/http"
)

type HandlerFunc func(Context) error
type ErrorHandler = func(Context, error) error

func (h *HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

type Context interface {
	context.Context
	Context() context.Context
	WithValue(key string, value any)
	EngineContext() any
	Header(key, defaultValue string) string
	Headers() http.Header
	Method() string
	Path() string
	Secure() bool
	Url() string
	FullUrl() string
	Ip() string

	//Input Retrieve  an input item from the request: /users/{id}
	Params(key string) string
	// Query Retrieve a query string item form the request: /users?id=1
	Query(key, defaultValue string) string
	// Form Retrieve a form string item form the post: /users POST:id=1
	Form(key, defaultValue string) string
	Bind(obj any) error
	Status(code int) Context
	AbortWithStatus(code int)
	Next() error

	Cookies(key string, defaultValue ...string) string
	Cookie(co *Cookie)

	File
	Request
	Response
}
