package http

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sujit-baniya/framework/contracts/http"
)

type FiberContext struct {
	instance *fiber.Ctx
}

func NewFiberContext(ctx *fiber.Ctx) http.Context {
	return &FiberContext{ctx}
}

func (c *FiberContext) Request() http.Request {
	return NewFiberRequest(c.instance)
}

func (c *FiberContext) Response() http.Response {
	return NewFiberResponse(c.instance)
}

func (c *FiberContext) WithValue(key string, value interface{}) {
	c.instance.Context().SetUserValue(key, value)
}

func (c *FiberContext) Deadline() (deadline time.Time, ok bool) {
	return c.instance.Context().Deadline()
}

func (c *FiberContext) Done() <-chan struct{} {
	return c.instance.Context().Done()
}

func (c *FiberContext) Err() error {
	return c.instance.Context().Err()
}

func (c *FiberContext) Value(key interface{}) interface{} {
	return c.instance.Context().Value(key)
}
