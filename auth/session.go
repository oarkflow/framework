package auth

import (
	"errors"
	"github.com/sujit-baniya/frame"
	"github.com/sujit-baniya/frame/middlewares/server/sessions"
	contractauth "github.com/sujit-baniya/framework/contracts/auth"
	"github.com/sujit-baniya/framework/facades"
)

type Session struct {
	guard string
}

func NewSession(guard string) contractauth.Auth {
	if facades.SessionAuth == nil {
		facades.SessionAuth = &Session{
			guard: guard,
		}
	}
	return facades.SessionAuth
}

// User need parse token first.
func (app *Session) User(ctx *frame.Context, user contractauth.User) error {
	session := sessions.Default(ctx)
	u := session.Get(ctx.AuthUserKey)
	if u == nil {
		user = nil
		return errors.New("Not logged in")
	}
	user = u.(contractauth.User)
	if _, ok := ctx.Get(ctx.AuthUserKey); !ok {
		ctx.Set(ctx.AuthUserKey, user)
	}
	return nil
}

func (app *Session) Parse(ctx *frame.Context, token string) error {
	return nil
}

func (app *Session) Login(ctx *frame.Context, user contractauth.User) (token string, err error) {
	session := sessions.Default(ctx)
	session.Set(ctx.AuthUserKey, user)
	err = session.Save()
	ctx.Set(ctx.AuthUserKey, user)
	return
}

func (app *Session) LoginUsingID(ctx *frame.Context, id any) (token string, err error) {
	return
}

// Refresh need parse token first.
func (app *Session) Refresh(ctx *frame.Context) (token string, err error) {
	return "", nil
}

func (app *Session) Logout(ctx *frame.Context) error {
	session := sessions.Default(ctx)
	session.Clear()
	return session.Save()
}
