package http

import (
	"github.com/gofiber/fiber/v2"
	"net/http"

	contracthttp "github.com/sujit-baniya/framework/contracts/http"
)

type FiberRequest struct {
	instance *fiber.Ctx
}

func NewFiberRequest(instance *fiber.Ctx) contracthttp.Request {
	return &FiberRequest{instance}
}

func (r *FiberRequest) Params(key string) string {
	return r.instance.Params(key)
}

func (r *FiberRequest) Query(key, defaultValue string) string {
	return r.instance.Query(key, defaultValue)
}

func (r *FiberRequest) Form(key, defaultValue string) string {
	return r.instance.FormValue(key, defaultValue)
}

func (r *FiberRequest) Bind(obj interface{}) error {
	return nil
}

func (r *FiberRequest) File(name string) (contracthttp.File, error) {
	file, err := r.instance.FormFile(name)
	if err != nil {
		return nil, err
	}

	return &FiberFile{instance: r.instance, file: file}, nil
}

func (r *FiberRequest) Header(key, defaultValue string) string {
	header := r.instance.Get(key)
	if header != "" {
		return header
	}

	return defaultValue
}

func (r *FiberRequest) Headers() http.Header {
	mp := make(map[string][]string)
	headers := r.instance.GetReqHeaders()
	for key, header := range headers {
		mp[key] = []string{header}
	}
	return mp
}

func (r *FiberRequest) Method() string {
	return r.instance.Method()
}

func (r *FiberRequest) Url() string {
	return r.instance.OriginalURL()
}

func (r *FiberRequest) FullUrl() string {
	prefix := "https://"
	if !r.instance.Secure() {
		prefix = "http://"
	}

	if r.instance.Hostname() == "" {
		return ""
	}

	return prefix + string(r.instance.Request().Host()) + string(r.instance.Request().RequestURI())
}

func (r *FiberRequest) AbortWithStatus(code int) {
	r.instance.Status(code)
}

func (r *FiberRequest) Next() error {
	return r.instance.Next()
}

func (r *FiberRequest) Path() string {
	return string(r.instance.Request().URI().Path())
}

func (r *FiberRequest) Ip() string {
	return r.instance.IP()
}

func (r *FiberRequest) Origin() *http.Request {
	return &http.Request{}
}

func (r *FiberRequest) Response() contracthttp.Response {
	return NewFiberResponse(r.instance)
}
