package http

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-chi/chi/v5"
	"net"
	"net/http"
	"strings"

	contracthttp "github.com/sujit-baniya/framework/contracts/http"
)

var trueClientIP = http.CanonicalHeaderKey("True-Client-IP")
var xForwardedFor = http.CanonicalHeaderKey("X-Forwarded-For")
var xRealIP = http.CanonicalHeaderKey("X-Real-IP")

type ChiRequest struct {
	request  *http.Request
	response http.ResponseWriter
	ctx      *chi.Context
}

func NewChiRequest(instance *http.Request, response http.ResponseWriter) contracthttp.Request {
	return &ChiRequest{request: instance, response: response}
}

func (r *ChiRequest) Params(key string) string {
	return chi.URLParam(r.request, key)
}

func (r *ChiRequest) Query(key, defaultValue string) string {
	q := r.request.URL.Query().Get(key)
	if q == "" {
		q = defaultValue
	}
	return q
}

func (r *ChiRequest) Form(key, defaultValue string) string {
	q := r.request.Form.Get(key)
	if q == "" {
		q = defaultValue
	}
	return q
}

func (r *ChiRequest) Bind(obj interface{}) error {
	b := binding.Default(r.request.Method, r.request.Header.Get("Content-Type"))
	return b.Bind(r.request, obj)
}

func (r *ChiRequest) File(name string) (contracthttp.File, error) {
	_, fileHeader, err := r.request.FormFile(name)
	if err != nil {
		return nil, err
	}

	return &ChiFile{request: r.request, file: fileHeader}, nil
}

func (r *ChiRequest) Header(key, defaultValue string) string {
	header := r.request.Header.Get(key)
	if header != "" {
		return header
	}

	return defaultValue
}

func (r *ChiRequest) Headers() http.Header {
	return r.request.Header
}

func (r *ChiRequest) Method() string {
	return r.request.Method
}

func (r *ChiRequest) Url() string {
	return r.request.RequestURI
}

func (r *ChiRequest) FullUrl() string {
	prefix := "https://"
	if r.request.TLS == nil {
		prefix = "http://"
	}

	if r.request.Host == "" {
		return ""
	}

	return prefix + r.request.Host + r.request.RequestURI
}

func (r *ChiRequest) AbortWithStatus(code int) {
	r.response.WriteHeader(code)
}

func (r *ChiRequest) Next() error {
	handler := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	}
	handler()
	return nil
}

func (r *ChiRequest) Path() string {
	return r.request.URL.Path
}

func (r *ChiRequest) Ip() string {
	var ip string

	if tcip := r.request.Header.Get(trueClientIP); tcip != "" {
		ip = tcip
	} else if xrip := r.request.Header.Get(xRealIP); xrip != "" {
		ip = xrip
	} else if xff := r.request.Header.Get(xForwardedFor); xff != "" {
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

func (r *ChiRequest) Origin() *http.Request {
	return &http.Request{}
}

func (r *ChiRequest) Response() contracthttp.Response {
	return NewChiResponse(r.response)
}
