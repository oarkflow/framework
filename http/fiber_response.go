package http

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"net/http"

	httpcontract "github.com/sujit-baniya/framework/contracts/http"
)

type FiberResponse struct {
	instance *fiber.Ctx
}

func NewFiberResponse(instance *fiber.Ctx) httpcontract.Response {
	return &FiberResponse{instance: instance}
}

func (r *FiberResponse) String(code int, format string, values ...interface{}) error {
	return r.instance.Status(code).SendString(fmt.Sprintf(format, values...))
}

func (r *FiberResponse) Json(code int, obj interface{}) error {
	return r.instance.Status(code).JSON(obj)
}

func (r *FiberResponse) File(filepath string, compress ...bool) error {
	return r.instance.SendFile(filepath, compress...)
}

func (r *FiberResponse) Download(filepath, filename string) error {
	return r.instance.Download(filepath, filename)
}

func (r *FiberResponse) Success() httpcontract.ResponseSuccess {
	return NewFiberSuccess(r.instance)
}

func (r *FiberResponse) Header(key, value string) httpcontract.Response {
	r.instance.Set(key, value)
	return r
}

type FiberSuccess struct {
	instance *fiber.Ctx
}

func NewFiberSuccess(instance *fiber.Ctx) httpcontract.ResponseSuccess {
	return &FiberSuccess{instance}
}

func (r *FiberSuccess) String(format string, values ...interface{}) error {
	return r.instance.Status(http.StatusOK).SendString(fmt.Sprintf(format, values...))
}

func (r *FiberSuccess) Json(obj interface{}) error {
	return r.instance.Status(http.StatusOK).JSON(obj)
}
