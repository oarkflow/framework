package auth

import (
	"github.com/oarkflow/frame"
)

type Auth interface {
	Guard(name string) Auth
	Parse(ctx *frame.Context, token string, user User) error
	User(ctx *frame.Context, user User) error
	Login(ctx *frame.Context, user User, data ...map[string]any) (token string, err error)
	LoginUsingID(ctx *frame.Context, id interface{}) (token string, err error)
	Refresh(ctx *frame.Context) (token string, err error)

	Data(ctx *frame.Context) (map[string]any, error)
	Logout(ctx *frame.Context) error
}

type User interface {
	Authenticated() bool
}
