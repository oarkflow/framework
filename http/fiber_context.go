package http

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	contracthttp "github.com/sujit-baniya/framework/contracts/http"
)

type FiberContext struct {
	instance *fiber.Ctx
}

func NewFiberContext(ctx *fiber.Ctx) contracthttp.Context {
	return &FiberContext{ctx}
}

func (c *FiberContext) Request() contracthttp.Request {
	return NewFiberRequest(c.instance)
}

func (c *FiberContext) Response() contracthttp.Response {
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

func (c *FiberContext) Params(key string) string {
	return c.instance.Params(key)
}

func (c *FiberContext) Query(key, defaultValue string) string {
	return c.instance.Query(key, defaultValue)
}

func (c *FiberContext) Form(key, defaultValue string) string {
	return c.instance.FormValue(key, defaultValue)
}

func (c *FiberContext) Bind(obj interface{}) error {
	return nil
}

func (c *FiberContext) File(name string) (contracthttp.File, error) {
	file, err := c.instance.FormFile(name)
	if err != nil {
		return nil, err
	}

	return &FiberFile{instance: c.instance, file: file}, nil
}

func (c *FiberContext) Header(key, defaultValue string) string {
	header := c.instance.Get(key)
	if header != "" {
		return header
	}

	return defaultValue
}

func (c *FiberContext) Headers() http.Header {
	mp := make(map[string][]string)
	headers := c.instance.GetReqHeaders()
	for key, header := range headers {
		mp[key] = []string{header}
	}
	return mp
}

func (c *FiberContext) Method() string {
	return c.instance.Method()
}

func (c *FiberContext) Url() string {
	return c.instance.OriginalURL()
}

func (c *FiberContext) FullUrl() string {
	prefix := "https://"
	if !c.instance.Secure() {
		prefix = "http://"
	}

	if c.instance.Hostname() == "" {
		return ""
	}

	return prefix + string(c.instance.Request().Host()) + string(c.instance.Request().RequestURI())
}

func (c *FiberContext) AbortWithStatus(code int) {
	c.instance.Status(code)
}

func (c *FiberContext) Next() error {
	return c.instance.Next()
}

func (c *FiberContext) Path() string {
	return string(c.instance.Request().URI().Path())
}

func (c *FiberContext) Ip() string {
	return c.instance.IP()
}
