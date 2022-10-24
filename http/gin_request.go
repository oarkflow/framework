package http

import (
	"github.com/sujit-baniya/framework/view"
	"net/http"

	"github.com/gin-gonic/gin"

	contracthttp "github.com/sujit-baniya/framework/contracts/http"
)

type GinRequest struct {
	instance *gin.Context
	config   GinConfig
	view     *view.Engine
}

func NewGinRequest(instance *gin.Context, config GinConfig, engine *view.Engine) contracthttp.Request {
	return &GinRequest{instance: instance, config: config, view: engine}
}

func (r *GinRequest) Origin() *http.Request {
	return r.instance.Request
}

func (r *GinRequest) Response() contracthttp.Response {
	return NewGinResponse(r.instance, r.config, r.view)
}
