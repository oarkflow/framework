package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	contracthttp "github.com/sujit-baniya/framework/contracts/http"
)

type GinContext struct {
	instance *gin.Context
}

func NewGinContext(ctx *gin.Context) contracthttp.Context {
	return &GinContext{ctx}
}

func (c *GinContext) Request() contracthttp.Request {
	return NewGinRequest(c.instance)
}

func (c *GinContext) Response() contracthttp.Response {
	return NewGinResponse(c.instance)
}

func (c *GinContext) WithValue(key string, value interface{}) {
	c.instance.Set(key, value)
}

func (c *GinContext) Deadline() (deadline time.Time, ok bool) {
	return c.instance.Deadline()
}

func (c *GinContext) Done() <-chan struct{} {
	return c.instance.Done()
}

func (c *GinContext) Err() error {
	return c.instance.Err()
}

func (c *GinContext) Value(key interface{}) interface{} {
	return c.instance.Value(key)
}

func (c *GinContext) Params(key string) string {
	return c.instance.Param(key)
}

func (c *GinContext) Query(key, defaultValue string) string {
	return c.instance.DefaultQuery(key, defaultValue)
}

func (c *GinContext) Form(key, defaultValue string) string {
	return c.instance.DefaultPostForm(key, defaultValue)
}

func (c *GinContext) Bind(obj interface{}) error {
	return c.instance.ShouldBind(obj)
}

func (c *GinContext) File(name string) (contracthttp.File, error) {
	file, err := c.instance.FormFile(name)
	if err != nil {
		return nil, err
	}

	return &GinFile{instance: c.instance, file: file}, nil
}

func (c *GinContext) Header(key, defaultValue string) string {
	header := c.instance.GetHeader(key)
	if header != "" {
		return header
	}

	return defaultValue
}

func (c *GinContext) Headers() http.Header {
	return c.instance.Request.Header
}

func (c *GinContext) Method() string {
	return c.instance.Request.Method
}

func (c *GinContext) Url() string {
	return c.instance.Request.RequestURI
}

func (c *GinContext) FullUrl() string {
	prefix := "https://"
	if c.instance.Request.TLS == nil {
		prefix = "http://"
	}

	if c.instance.Request.Host == "" {
		return ""
	}

	return prefix + c.instance.Request.Host + c.instance.Request.RequestURI
}

func (c *GinContext) AbortWithStatus(code int) {
	c.instance.AbortWithStatus(code)
}

func (c *GinContext) Next() error {
	c.instance.Next()
	return nil
}

func (c *GinContext) Path() string {
	return c.instance.Request.URL.Path
}

func (c *GinContext) Ip() string {
	return c.instance.ClientIP()
}
