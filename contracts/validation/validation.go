package validation

import "github.com/oarkflow/frame"

type Option func(map[string]any)
type Validation interface {
	Make(ctx *frame.Context, data any, rules map[string]string, options ...Option) (Validator, error)
	AddRules([]Rule) error
	Rules() []Rule
}
type Validator interface {
	Bind(ptr any) error
	Errors() Errors
	Fails() bool
}
type Errors interface {
	One(key ...string) string
	Get(key string) map[string]string
	All() map[string]map[string]string
	Has(key string) bool
}
type Data interface {
	Get(key string) (val any, exist bool)
	Set(key string, val any) error
}
type Rule interface {
	Signature() string
	Passes(ctx *frame.Context, data Data, val any, options ...any) bool
	Message() string
}
