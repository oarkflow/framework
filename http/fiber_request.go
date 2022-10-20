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

func (r *FiberRequest) Origin() *http.Request {
	return &http.Request{}
}

func (r *FiberRequest) Response() contracthttp.Response {
	return NewFiberResponse(r.instance)
}
