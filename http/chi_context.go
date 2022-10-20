package http

import (
	"context"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-chi/chi/v5"
	"net"
	"strings"
	"time"

	contracthttp "github.com/sujit-baniya/framework/contracts/http"
	"net/http"
)

type ChiContext struct {
	request  *http.Request
	response http.ResponseWriter
}

func NewChiContext(request *http.Request, response http.ResponseWriter) contracthttp.Context {
	return &ChiContext{request: request, response: response}
}

func (c *ChiContext) Request() contracthttp.Request {
	return NewChiRequest(c.request, c.response)
}

func (c *ChiContext) Response() contracthttp.Response {
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

func (c *ChiContext) Params(key string) string {
	return chi.URLParam(c.request, key)
}

func (c *ChiContext) Query(key, defaultValue string) string {
	q := c.request.URL.Query().Get(key)
	if q == "" {
		q = defaultValue
	}
	return q
}

func (c *ChiContext) Form(key, defaultValue string) string {
	q := c.request.Form.Get(key)
	if q == "" {
		q = defaultValue
	}
	return q
}

func (c *ChiContext) Bind(obj interface{}) error {
	b := binding.Default(c.request.Method, c.request.Header.Get("Content-Type"))
	return b.Bind(c.request, obj)
}

func (c *ChiContext) File(name string) (contracthttp.File, error) {
	_, fileHeader, err := c.request.FormFile(name)
	if err != nil {
		return nil, err
	}

	return &ChiFile{request: c.request, file: fileHeader}, nil
}

func (c *ChiContext) Header(key, defaultValue string) string {
	header := c.request.Header.Get(key)
	if header != "" {
		return header
	}

	return defaultValue
}

func (c *ChiContext) Headers() http.Header {
	return c.request.Header
}

func (c *ChiContext) Method() string {
	return c.request.Method
}

func (c *ChiContext) Url() string {
	return c.request.RequestURI
}

func (c *ChiContext) FullUrl() string {
	prefix := "https://"
	if c.request.TLS == nil {
		prefix = "http://"
	}

	if c.request.Host == "" {
		return ""
	}

	return prefix + c.request.Host + c.request.RequestURI
}

func (c *ChiContext) AbortWithStatus(code int) {
	c.response.WriteHeader(code)
}

func (c *ChiContext) Next() error {
	return nil
}

func (c *ChiContext) Path() string {
	return c.request.URL.Path
}

func (c *ChiContext) Ip() string {
	var ip string

	if tcip := c.request.Header.Get(trueClientIP); tcip != "" {
		ip = tcip
	} else if xrip := c.request.Header.Get(xRealIP); xrip != "" {
		ip = xrip
	} else if xff := c.request.Header.Get(xForwardedFor); xff != "" {
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
