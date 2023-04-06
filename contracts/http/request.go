package http

import (
	"github.com/oarkflow/frame"

	"github.com/oarkflow/framework/contracts/validation"
)

type FormRequest interface {
	Authorize(ctx *frame.Context) error
	Rules() map[string]string
	Messages() map[string]string
	Attributes() map[string]string
	PrepareForValidation(data validation.Data)
}
