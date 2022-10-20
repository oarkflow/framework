package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	httpcontract "github.com/sujit-baniya/framework/contracts/http"
)

type GinResponse struct {
	instance *gin.Context
}

func NewGinResponse(instance *gin.Context) httpcontract.Response {
	return &GinResponse{instance: instance}
}

func (r *GinResponse) String(code int, format string, values ...interface{}) error {
	r.instance.String(code, format, values...)
	return nil
}

func (r *GinResponse) Json(code int, obj interface{}) error {
	r.instance.JSON(code, obj)
	return nil
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

type GinSuccess struct {
	instance *gin.Context
}

func NewGinSuccess(instance *gin.Context) httpcontract.ResponseSuccess {
	return &GinSuccess{instance}
}

func (r *GinSuccess) String(format string, values ...interface{}) error {
	r.instance.String(http.StatusOK, format, values...)
	return nil
}

func (r *GinSuccess) Json(obj interface{}) error {
	r.instance.JSON(http.StatusOK, obj)
	return nil
}
