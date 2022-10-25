package http

import (
	"github.com/sujit-baniya/framework/view"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	contracthttp "github.com/sujit-baniya/framework/contracts/http"
)

type GinConfig struct {
	Mode        string `json:"mode"`
	ViewsLayout string `json:"views_layout"`
	Extension   string `json:"extension"`
	Path        string `json:"path"`
	View        *view.Engine
}

type GinContext struct {
	instance *gin.Context
	config   GinConfig
}

func NewGinContext(ctx *gin.Context, config GinConfig) contracthttp.Context {
	ct := &GinContext{instance: ctx, config: config}
	return ct
}

func (c *GinContext) Origin() *http.Request {
	return c.instance.Request
}

func (c *GinContext) String(code int, format string, values ...any) error {
	c.instance.String(code, format, values...)
	return nil
}

func (c *GinContext) Json(code int, obj any) error {
	c.instance.JSON(code, obj)
	return nil
}

func (c *GinContext) Render(name string, bind any, layouts ...string) error {
	return c.config.View.Render(c.instance.Writer, name, bind, layouts...)
}

func (c *GinContext) SendFile(filepath string, compress ...bool) error {
	c.instance.File(filepath)
	return nil
}

func (c *GinContext) Download(filepath, filename string) error {
	c.instance.FileAttachment(filepath, filename)
	return nil
}

func (c *GinContext) SetHeader(key, value string) contracthttp.Context {
	c.instance.Header(key, value)
	return c
}

func (c *GinContext) StatusCode() int {
	return c.instance.Writer.Status()
}

func (c *GinContext) Vary(field string, values ...string) {
	c.Append(field)
}

func (c *GinContext) Append(field string, values ...string) {
	if len(values) == 0 {
		return
	}
	h := c.instance.GetHeader(field)
	originalH := h
	for _, value := range values {
		if len(h) == 0 {
			h = value
		} else if h != value && !strings.HasPrefix(h, value+",") && !strings.HasSuffix(h, " "+value) &&
			!strings.Contains(h, " "+value+",") {
			h += ", " + value
		}
	}
	if originalH != h {
		c.SetHeader(field, h)
	}
}

func (c *GinContext) WithValue(key string, value any) {
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

func (c *GinContext) Value(key any) any {
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

func (c *GinContext) Bind(obj any) error {
	return c.instance.Bind(obj)
}

func (c *GinContext) SaveFile(name string, dst string) error {
	file, err := c.File(name)
	if err != nil {
		return err
	}
	return c.instance.SaveUploadedFile(file, dst)
}

func (c *GinContext) File(name string) (*multipart.FileHeader, error) {
	return c.instance.FormFile(name)
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

func (c *GinContext) Cookies(key string, defaultValue ...string) string {
	str, _ := c.instance.Cookie(key)
	return str
}

func (c *GinContext) Cookie(co *contracthttp.Cookie) {
	switch co.SameSite {
	case "Lax":
		c.instance.SetSameSite(http.SameSiteLaxMode)
		break
	case "None":
		c.instance.SetSameSite(http.SameSiteNoneMode)
		break
	case "Strict":
		c.instance.SetSameSite(http.SameSiteStrictMode)
		break
	default:
		c.instance.SetSameSite(http.SameSiteDefaultMode)
	}
	c.instance.SetCookie(co.Name, co.Value, co.MaxAge, co.Path, co.Domain, co.Secure, co.HTTPOnly)
}

func (c *GinContext) Path() string {
	return c.instance.Request.URL.Path
}

func (c *GinContext) EngineContext() any {
	return c.instance
}

func (c *GinContext) Secure() bool {
	if c.instance.Request.Proto == "https" {
		return true
	}
	return false
}

func (c *GinContext) Ip() string {
	return c.instance.ClientIP()
}
