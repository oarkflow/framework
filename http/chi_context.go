package http

import (
	"context"
	"github.com/gin-gonic/gin/binding"
	"github.com/sujit-baniya/chi"
	"net"
	"strings"
	"time"

	contracthttp "github.com/sujit-baniya/framework/contracts/http"
	"net/http"
)

type ChiContext struct {
	Req  *http.Request
	Res  http.ResponseWriter
	next http.Handler
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

func NewChiContext(request *http.Request, response http.ResponseWriter, n ...http.Handler) contracthttp.Context {
	var next http.Handler
	if len(n) > 0 {
		next = n[0]
	}
	return &ChiContext{Req: request, Res: response, next: next}
}

func (c *ChiContext) Request() contracthttp.Request {
	return NewChiRequest(c.Req, c.Res)
}

func (c *ChiContext) Response() contracthttp.Response {
	return NewChiResponse(c.Res)
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

func (c *ChiContext) File(name string) (contracthttp.File, error) {
	_, fileHeader, err := c.Req.FormFile(name)
	if err != nil {
		return nil, err
	}

	return &ChiFile{request: c.Req, file: fileHeader}, nil
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
