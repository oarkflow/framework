package http

import (
	"github.com/go-chi/chi/v5"
	contracthttp "github.com/sujit-baniya/framework/contracts/http"
	"net/http"
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

func (r *ChiRequest) Origin() *http.Request {
	return &http.Request{}
}

func (r *ChiRequest) Response() contracthttp.Response {
	return NewChiResponse(r.response)
}
