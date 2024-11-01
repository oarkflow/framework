package auth

import (
	"time"

	"github.com/oarkflow/frame"
)

type AccessToken struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

type Auth interface {
	Guard(name string) Auth
	Parse(ctx *frame.Context, token string, user User) error
	User(ctx *frame.Context, user User) error
	Login(ctx *frame.Context, user User, data ...map[string]any) (token *AccessToken, err error)
	LoginUsingID(ctx *frame.Context, id interface{}) (token *AccessToken, err error)
	Refresh(ctx *frame.Context) (token *AccessToken, err error)
	Data(ctx *frame.Context) (map[string]any, error)
	Logout(ctx *frame.Context) error
}

type User interface {
	Authenticated() bool
}

type UserFn func() User
