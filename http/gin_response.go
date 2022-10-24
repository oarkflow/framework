package http

import (
	"github.com/sujit-baniya/framework/view"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	httpcontract "github.com/sujit-baniya/framework/contracts/http"
)

type GinConfig struct {
	Mode        string `json:"mode"`
	ViewsLayout string `json:"views_layout"`
	Extension   string `json:"extension"`
	Path        string `json:"path"`
	View        *view.Engine
}

type GinResponse struct {
	instance *gin.Context
	config   GinConfig
	view     *view.Engine
}

func NewGinResponse(instance *gin.Context, config GinConfig, engine *view.Engine) httpcontract.Response {
	return &GinResponse{instance: instance, config: config, view: engine}
}

func (r *GinResponse) String(code int, format string, values ...any) error {
	r.instance.String(code, format, values...)
	return nil
}

func (r *GinResponse) Json(code int, obj any) error {
	r.instance.JSON(code, obj)
	return nil
}

func (r *GinResponse) Render(name string, bind any, layouts ...string) error {
	return r.view.Render(r.instance.Writer, name, bind, layouts...)
}

func (r *GinResponse) File(filepath string, compress ...bool) error {
	r.instance.File(filepath)
	return nil
}

func (r *GinResponse) Download(filepath, filename string) error {
	r.instance.FileAttachment(filepath, filename)
	return nil
}

func (r *GinResponse) Success() httpcontract.ResponseSuccess {
	return NewGinSuccess(r.instance)
}

func (r *GinResponse) Header(key, value string) httpcontract.Response {
	r.instance.Header(key, value)
	return r
}

func (r *GinResponse) StatusCode() int {
	return r.instance.Writer.Status()
}

func (r *GinResponse) Vary(field string, values ...string) {
	r.Append(field)
}

func (r *GinResponse) Append(field string, values ...string) {
	if len(values) == 0 {
		return
	}
	h := r.instance.GetHeader(field)
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

type GinSuccess struct {
	instance *gin.Context
}

func NewGinSuccess(instance *gin.Context) httpcontract.ResponseSuccess {
	return &GinSuccess{instance}
}

func (r *GinSuccess) String(format string, values ...any) error {
	r.instance.String(http.StatusOK, format, values...)
	return nil
}

func (r *GinSuccess) Json(obj any) error {
	r.instance.JSON(http.StatusOK, obj)
	return nil
}
