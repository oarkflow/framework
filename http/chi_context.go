package http

import (
	"context"
	"time"

	"github.com/sujit-baniya/framework/contracts/http"
	stdHttp "net/http"
)

type ChiContext struct {
	request  *stdHttp.Request
	response stdHttp.ResponseWriter
}

func NewChiContext(request *stdHttp.Request, response stdHttp.ResponseWriter) http.Context {
	return &ChiContext{request: request, response: response}
}

func (c *ChiContext) Request() http.Request {
	return NewChiRequest(c.request, c.response)
}

func (c *ChiContext) Response() http.Response {
	return NewChiResponse(c.response)
}

func (c *ChiContext) WithValue(key string, value interface{}) {
	ctx := context.WithValue(c.request.Context(), key, value)
	c.request = c.request.WithContext(ctx)
}

func (c *ChiContext) Deadline() (deadline time.Time, ok bool) {
	return c.request.Context().Deadline()
}

func (c *ChiContext) Done() <-chan struct{} {
	return c.request.Context().Done()
}

func (c *ChiContext) Err() error {
	return c.request.Context().Err()
}

func (c *ChiContext) Value(key interface{}) interface{} {
	return c.request.Context().Value(key)
}
