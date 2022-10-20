package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	contracthttp "github.com/sujit-baniya/framework/contracts/http"
)

type GinRequest struct {
	instance *gin.Context
}

func NewGinRequest(instance *gin.Context) contracthttp.Request {
	return &GinRequest{instance}
}

func (r *GinRequest) Origin() *http.Request {
	return r.instance.Request
}

func (r *GinRequest) Response() contracthttp.Response {
	return NewGinResponse(r.instance)
}
