package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	httpcontract "github.com/sujit-baniya/framework/contracts/http"
)

type ChiResponse struct {
	instance http.ResponseWriter
}

func (r *ChiResponse) StatusCode() int {
	return 200
}

func (r *ChiResponse) Vary(key string, value ...string) {
	r.Append(key)
}

func (r *ChiResponse) Append(field string, values ...string) {
	if len(values) == 0 {
		return
	}
	h := r.instance.Header().Get(field)
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
		r.Header(field, h)
	}
}

func NewChiResponse(instance http.ResponseWriter) httpcontract.Response {
	return &ChiResponse{instance: instance}
}

func (r *ChiResponse) String(code int, format string, values ...interface{}) error {
	r.instance.WriteHeader(code)
	_, err := r.instance.Write([]byte(fmt.Sprintf(format, values...)))
	return err
}

func (r *ChiResponse) Json(code int, obj interface{}) error {
	r.instance.WriteHeader(code)
	r.instance.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	_, err = r.instance.Write(jsonResp)
	return err
}

func (r *ChiResponse) File(filepath string, compress ...bool) error {
	//@TODO - Implement
	return nil
}

func (r *ChiResponse) Download(filepath, filename string) error {
	//@TODO - Implement
	return nil
}

func (r *ChiResponse) Success() httpcontract.ResponseSuccess {
	return NewChiSuccess(r.instance)
}

func (r *ChiResponse) Header(key, value string) httpcontract.Response {
	r.instance.Header().Set(key, value)
	return r
}

type ChiSuccess struct {
	instance http.ResponseWriter
}

func NewChiSuccess(instance http.ResponseWriter) httpcontract.ResponseSuccess {
	return &ChiSuccess{instance}
}

func (r *ChiSuccess) String(format string, values ...interface{}) error {
	r.instance.WriteHeader(http.StatusOK)
	_, err := r.instance.Write([]byte(fmt.Sprintf(format, values...)))
	return err
}

func (r *ChiSuccess) Json(obj interface{}) error {
	r.instance.WriteHeader(http.StatusOK)
	r.instance.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	_, err = r.instance.Write(jsonResp)
	return err
}

func (r *ChiSuccess) Render(name string, bind any, layouts ...string) error {
	return nil
}

func (r *ChiResponse) Render(name string, bind any, layouts ...string) error {
	return nil
}
