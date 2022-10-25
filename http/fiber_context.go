package http

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"

	"github.com/gofiber/fiber/v2"
	contracthttp "github.com/sujit-baniya/framework/contracts/http"
)

type FiberContext struct {
	instance *fiber.Ctx
}

func (c *FiberContext) Origin() *http.Request {
	headers := make(map[string][]string)
	for header, value := range c.instance.GetReqHeaders() {
		headers[header] = []string{value}
	}
	parsedUrl, _ := url.Parse(c.instance.OriginalURL())
	return &http.Request{
		Method: c.instance.Method(),
		URL:    parsedUrl,
		Proto:  c.instance.Protocol(),
		Header: headers,
		Body:   io.NopCloser(bytes.NewReader(c.instance.Body())),
		GetBody: func() (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader(c.instance.Body())), nil
		},
		Host:       c.instance.Hostname(),
		RemoteAddr: c.instance.IP(),
		RequestURI: c.instance.Request().URI().String(),
	}
}

func NewFiberContext(ctx *fiber.Ctx) contracthttp.Context {
	return &FiberContext{ctx}
}

func (c *FiberContext) WithValue(key string, value any) {
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

func (c *FiberContext) Value(key any) any {
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

func (c *FiberContext) Bind(obj any) error {
	return nil
}

func (c *FiberContext) SaveFile(name string, dst string) error {
	file, err := c.File(name)
	if err != nil {
		return err
	}
	return c.instance.SaveFile(file, dst)
}

func (c *FiberContext) File(name string) (*multipart.FileHeader, error) {
	return c.instance.FormFile(name)
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

func (c *FiberContext) Cookies(key string, defaultValue ...string) string {
	return c.instance.Cookies(key, defaultValue...)
}

func (c *FiberContext) Cookie(co *contracthttp.Cookie) {
	c.instance.Cookie(&fiber.Cookie{
		Name:        co.Name,
		Value:       co.Value,
		Path:        co.Path,
		Domain:      co.Domain,
		MaxAge:      co.MaxAge,
		Expires:     co.Expires,
		Secure:      co.Secure,
		HTTPOnly:    co.HTTPOnly,
		SameSite:    co.SameSite,
		SessionOnly: co.SessionOnly,
	})
}

func (c *FiberContext) Path() string {
	return string(c.instance.Request().URI().Path())
}

func (c *FiberContext) EngineContext() any {
	return c.instance
}

func (c *FiberContext) Secure() bool {
	return c.instance.Secure()
}

func (c *FiberContext) Ip() string {
	return c.instance.IP()
}

func (c *FiberContext) String(code int, format string, values ...any) error {
	return c.instance.Status(code).SendString(fmt.Sprintf(format, values...))
}

func (c *FiberContext) Json(code int, obj any) error {
	return c.instance.Status(code).JSON(obj)
}

func (c *FiberContext) SendFile(filepath string, compress ...bool) error {
	return c.instance.SendFile(filepath, compress...)
}

func (c *FiberContext) Download(filepath, filename string) error {
	return c.instance.Download(filepath, filename)
}

func (c *FiberContext) StatusCode() int {
	return c.instance.Response().StatusCode()
}

func (c *FiberContext) Render(name string, bind any, layouts ...string) error {
	return c.instance.Render(name, bind, layouts...)
}

func (c *FiberContext) SetHeader(key, value string) contracthttp.Context {
	c.instance.Set(key, value)
	return c
}

func (c *FiberContext) Vary(key string, value ...string) {
	c.instance.Vary(key)
}
