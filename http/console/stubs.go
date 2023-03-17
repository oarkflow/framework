package console

type Stubs struct {
}

func (r Stubs) Request() string {
	return `package requests

import (
	"github.com/oarkflow/frame"
	"github.com/oarkflow/framework/contracts/validation"
)

type DummyRequest struct {
	DummyField
}

func (r *DummyRequest) Authorize(ctx *frame.Context) error {
	return nil
}

func (r *DummyRequest) Rules() map[string]string {
	return map[string]string{}
}

func (r *DummyRequest) Messages() map[string]string {
	return map[string]string{}
}

func (r *DummyRequest) Attributes() map[string]string {
	return map[string]string{}
}

func (r *DummyRequest) PrepareForValidation(data validation.Data) {

}
`
}
