package auth

import (
	"errors"

	"github.com/oarkflow/frame"
	"github.com/oarkflow/frame/middlewares/server/session"
	"github.com/oarkflow/pkg/dto"

	"github.com/oarkflow/framework/contracts/auth"
)

type Session struct {
	guard string
	store *session.Store
}

func NewSession(guard string, store *session.Store) auth.Auth {
	return &Session{
		guard: guard,
		store: store,
	}
}

func (app *Session) Guard(name string) auth.Auth {
	return GetAuth(name)
}

// User need parse token first.
func (app *Session) User(ctx *frame.Context, user auth.User) error {
	s, err := session.Pick(ctx, app.store)
	if err != nil {
		return err
	}
	u := s.Get(ctx.AuthUserKey)
	if u == nil {
		user = nil
		return errors.New("not logged in")
	}
	switch v := u.(type) {
	case auth.User:
		user = v
	default:
		err := dto.Map(user, v)
		if err != nil {
			return err
		}
	}
	if _, ok := ctx.Get(ctx.AuthUserKey); !ok {
		ctx.Set(ctx.AuthUserKey, user)
	}
	return nil
}

func (app *Session) Parse(ctx *frame.Context, token string) error {
	return nil
}

func (app *Session) Data(ctx *frame.Context) (map[string]any, error) {
	s, err := session.Pick(ctx, app.store)
	if err != nil {
		return nil, err
	}
	data := make(map[string]any, len(s.Keys()))
	for _, key := range s.Keys() {
		data[key] = s.Get(key)
	}
	return data, nil
}

func (app *Session) Login(ctx *frame.Context, user auth.User, data ...map[string]any) (token string, err error) {
	s, err := session.Pick(ctx, app.store)
	if err != nil {
		return "", err
	}
	s.Set(ctx.AuthUserKey, user)
	if len(data) > 0 {
		for k, v := range data[0] {
			s.Set(k, v)
		}
	}
	s.Save()
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
	s, err := session.Pick(ctx, app.store)
	if err != nil {
		return err
	}
	return s.Destroy()
}
