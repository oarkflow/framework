package auth

import (
	"github.com/sujit-baniya/frame"
)

type Auth interface {
	Guard(name string) Auth
	Parse(ctx *frame.Context, token string) error
	User(ctx *frame.Context, user interface{}) error
	Login(ctx *frame.Context, user interface{}) (token string, err error)
	LoginUsingID(ctx *frame.Context, id interface{}) (token string, err error)
	Refresh(ctx *frame.Context) (token string, err error)
	Logout(ctx *frame.Context) error
}
