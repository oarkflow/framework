package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/sujit-baniya/chi"
	"github.com/sujit-baniya/framework/view"
	"mime/multipart"
	"net"
	"strings"
	"time"

	contracthttp "github.com/sujit-baniya/framework/contracts/http"
	"net/http"
)

var trueClientIP = http.CanonicalHeaderKey("True-Client-IP")
var xForwardedFor = http.CanonicalHeaderKey("X-Forwarded-For")
var xRealIP = http.CanonicalHeaderKey("X-Real-IP")

type ChiConfig struct {
	NotFoundHandler         http.HandlerFunc
	MethodNotAllowedHandler http.HandlerFunc
	View                    *view.Engine
}

type ChiContext struct {
	Req    *http.Request
	Res    http.ResponseWriter
	next   http.Handler
	config ChiConfig
}

func NewChiContext(request *http.Request, response http.ResponseWriter, config ChiConfig, n ...http.Handler) contracthttp.Context {
	var next http.Handler
	if len(n) > 0 {
		next = n[0]
	}
	return &ChiContext{Req: request, Res: response, next: next, config: config}
}

func (c *ChiContext) Origin() any {
	return c.Req
}

func (c *ChiContext) Secure() bool {
	//TODO implement me
	panic("implement me")
}

func (c *ChiContext) Cookies(key string, defaultValue ...string) string {
	//TODO implement me
	panic("implement me")
}

func (c *ChiContext) Cookie(co *contracthttp.Cookie) {
	//TODO implement me
	panic("implement me")
}

func (c *ChiContext) StatusCode() int {
	return 200
}

func (c *ChiContext) Vary(key string, value ...string) {
	c.Append(key)
}

func (c *ChiContext) Append(field string, values ...string) {
	if len(values) == 0 {
		return
	}
	h := c.Res.Header().Get(field)
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

func (c *ChiContext) String(code int, format string, values ...interface{}) error {
	c.Res.WriteHeader(code)
	_, err := c.Res.Write([]byte(fmt.Sprintf(format, values...)))
	return err
}

func (c *ChiContext) Json(code int, obj interface{}) error {
	c.Res.WriteHeader(code)
	c.Res.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	_, err = c.Res.Write(jsonResp)
	return err
}

func (c *ChiContext) SendFile(filepath string, compress ...bool) error {
	//@TODO - Implement
	return nil
}

func (c *ChiContext) Download(filepath, filename string) error {
	//@TODO - Implement
	return nil
}

func (c *ChiContext) SetHeader(key, value string) contracthttp.Context {
	c.Res.Header().Set(key, value)
	return c
}

func (c *ChiContext) Render(name string, bind any, layouts ...string) error {
	return c.config.View.Render(c.Res, name, bind, layouts...)
}

func (c *ChiContext) WithValue(key string, value interface{}) {
	ctx := context.WithValue(c.Req.Context(), key, value)
	c.Req = c.Req.WithContext(ctx)
}

func (c *ChiContext) Deadline() (deadline time.Time, ok bool) {
	return c.Req.Context().Deadline()
}

func (c *ChiContext) Done() <-chan struct{} {
	return c.Req.Context().Done()
}

func (c *ChiContext) Err() error {
	return c.Req.Context().Err()
}

func (c *ChiContext) Value(key interface{}) interface{} {
	return c.Req.Context().Value(key)
}

func (c *ChiContext) Params(key string) string {
	return chi.URLParam(c.Req, key)
}

func (c *ChiContext) Query(key, defaultValue string) string {
	q := c.Req.URL.Query().Get(key)
	if q == "" {
		q = defaultValue
	}
	return q
}

func (c *ChiContext) Form(key, defaultValue string) string {
	q := c.Req.Form.Get(key)
	if q == "" {
		q = defaultValue
	}
	return q
}

func (c *ChiContext) Bind(obj interface{}) error {
	b := binding.Default(c.Req.Method, c.Req.Header.Get("Content-Type"))
	return b.Bind(c.Req, obj)
}

func (c *ChiContext) SaveFile(name string, dst string) error {
	//@TODO - Implement file save
	_, err := c.File(name)
	if err != nil {
		return err
	}

	return nil // f.Req.MultipartForm.File
}

func (c *ChiContext) File(name string) (*multipart.FileHeader, error) {
	_, fileHeader, err := c.Req.FormFile(name)
	return fileHeader, err
}

func (c *ChiContext) Header(key, defaultValue string) string {
	header := c.Req.Header.Get(key)
	if header != "" {
		return header
	}

	return defaultValue
}

func (c *ChiContext) Headers() http.Header {
	return c.Req.Header
}

func (c *ChiContext) Method() string {
	return c.Req.Method
}

func (c *ChiContext) Url() string {
	return c.Req.RequestURI
}

func (c *ChiContext) FullUrl() string {
	prefix := "https://"
	if c.Req.TLS == nil {
		prefix = "http://"
	}

	if c.Req.Host == "" {
		return ""
	}

	return prefix + c.Req.Host + c.Req.RequestURI
}

func (c *ChiContext) AbortWithStatus(code int) {
	c.Res.WriteHeader(code)
}

func (c *ChiContext) Next() error {
	if c.next != nil {
		c.next.ServeHTTP(c.Res, c.Req)
	}
	return nil
}

func (c *ChiContext) Path() string {
	return c.Req.URL.Path
}

func (c *ChiContext) EngineContext() any {
	return c
}

func (c *ChiContext) Ip() string {
	var ip string

	if tcip := c.Req.Header.Get(trueClientIP); tcip != "" {
		ip = tcip
	} else if xrip := c.Req.Header.Get(xRealIP); xrip != "" {
		ip = xrip
	} else if xff := c.Req.Header.Get(xForwardedFor); xff != "" {
		i := strings.Index(xff, ",")
		if i == -1 {
			i = len(xff)
		}
		ip = xff[:i]
	}
	if ip == "" || net.ParseIP(ip) == nil {
		return ""
	}
	return ip
}
