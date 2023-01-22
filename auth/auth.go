package auth

import (
	"errors"
	"github.com/sujit-baniya/frame"
	"reflect"
	"strings"
	"time"

	contractauth "github.com/sujit-baniya/framework/contracts/auth"
	"github.com/sujit-baniya/framework/facades"
	supporttime "github.com/sujit-baniya/framework/support/time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/cast"
)

const ctxKey = "GoravelAuth"

var (
	unit = time.Minute

	ErrorRefreshTimeExceeded = errors.New("refresh time exceeded")
	ErrorTokenExpired        = errors.New("token expired")
	ErrorNoPrimaryKeyField   = errors.New("the primaryKey field was not found in the model, set primaryKey like orm.Model")
	ErrorEmptySecret         = errors.New("secret is required")
	ErrorTokenDisabled       = errors.New("token is disabled")
	ErrorParseTokenFirst     = errors.New("parse token first")
	ErrorInvalidClaims       = errors.New("invalid claims")
	ErrorInvalidToken        = errors.New("invalid token")
)

type Claims struct {
	Key string `json:"key"`
	jwt.RegisteredClaims
}

type Guard struct {
	Claims *Claims
	Token  string
}

type Guards map[string]*Guard

type Auth struct {
	guard string
}

func NewAuth(guard string) contractauth.Auth {
	return &Auth{
		guard: guard,
	}
}

func (app *Auth) Guard(name string) contractauth.Auth {
	return NewAuth(name)
}

// User need parse token first.
func (app *Auth) User(ctx *frame.Context, user contractauth.User) error {
	auth, ok := ctx.Value(ctxKey).(Guards)
	if !ok || auth[app.guard] == nil {
		return ErrorParseTokenFirst
	}
	if auth[app.guard].Claims == nil {
		return ErrorParseTokenFirst
	}
	if auth[app.guard].Token == "" {
		return ErrorTokenExpired
	}
	if err := facades.Orm.Query().Find(user, auth[app.guard].Claims.Key); err != nil {
		return err
	}

	return nil
}

func (app *Auth) Parse(ctx *frame.Context, token string) error {
	token = strings.ReplaceAll(token, "Bearer ", "")
	if tokenIsDisabled(token) {
		return ErrorTokenDisabled
	}

	jwtSecret := facades.Config.GetString("jwt.secret")
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (any, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		if strings.Contains(err.Error(), jwt.ErrTokenExpired.Error()) && tokenClaims != nil {
			claims, ok := tokenClaims.Claims.(*Claims)
			if !ok {
				return ErrorInvalidClaims
			}

			app.makeAuthContext(ctx, claims, "")

			return ErrorTokenExpired
		} else {
			return err
		}
	}
	if tokenClaims == nil || !tokenClaims.Valid {
		return ErrorInvalidToken
	}

	claims, ok := tokenClaims.Claims.(*Claims)
	if !ok {
		return ErrorInvalidClaims
	}

	app.makeAuthContext(ctx, claims, token)

	return nil
}

func (app *Auth) Login(ctx *frame.Context, user contractauth.User) (token string, err error) {
	t := reflect.TypeOf(user).Elem()
	v := reflect.ValueOf(user).Elem()
	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).Name == "Model" {
			if v.Field(i).Type().Kind() == reflect.Struct {
				structField := v.Field(i).Type()
				for j := 0; j < structField.NumField(); j++ {
					if structField.Field(j).Tag.Get("gorm") == "primaryKey" {
						ctx.Set(ctx.AuthUserKey, user)
						return app.LoginUsingID(ctx, v.Field(i).Field(j).Interface())
					}
				}
			}
		}
		if t.Field(i).Tag.Get("gorm") == "primaryKey" {
			ctx.Set(ctx.AuthUserKey, user)
			return app.LoginUsingID(ctx, v.Field(i).Interface())
		}
	}

	return "", ErrorNoPrimaryKeyField
}

func (app *Auth) LoginUsingID(ctx *frame.Context, id any) (token string, err error) {
	jwtSecret := facades.Config.GetString("jwt.secret")
	if jwtSecret == "" {
		return "", ErrorEmptySecret
	}

	nowTime := supporttime.Now()
	ttl := facades.Config.GetInt("jwt.ttl")
	expireTime := nowTime.Add(time.Duration(ttl) * unit)
	claims := Claims{
		cast.ToString(id),
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(nowTime),
			Subject:   app.guard,
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = tokenClaims.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	app.makeAuthContext(ctx, &claims, token)

	return
}

// Refresh need parse token first.
func (app *Auth) Refresh(ctx *frame.Context) (token string, err error) {
	auth, ok := ctx.Value(ctxKey).(Guards)
	if !ok || auth[app.guard] == nil {
		return "", ErrorParseTokenFirst
	}
	if auth[app.guard].Claims == nil {
		return "", ErrorParseTokenFirst
	}

	nowTime := supporttime.Now()
	refreshTtl := facades.Config.GetInt("jwt.refresh_ttl")
	expireTime := auth[app.guard].Claims.ExpiresAt.Add(time.Duration(refreshTtl) * unit)
	if nowTime.Unix() > expireTime.Unix() {
		return "", ErrorRefreshTimeExceeded
	}

	return app.LoginUsingID(ctx, auth[app.guard].Claims.Key)
}

func (app *Auth) Logout(ctx *frame.Context) error {
	auth, ok := ctx.Value(ctxKey).(Guards)
	if !ok || auth[app.guard] == nil || auth[app.guard].Token == "" {
		return nil
	}

	if facades.Cache == nil {
		return errors.New("cache support is required")
	}

	if err := facades.Cache.Put(getDisabledCacheKey(auth[app.guard].Token),
		true,
		time.Duration(facades.Config.GetInt("jwt.ttl"))*unit,
	); err != nil {
		return err
	}

	delete(auth, app.guard)
	ctx.Set(ctxKey, auth)

	return nil
}

func (app *Auth) makeAuthContext(ctx *frame.Context, claims *Claims, token string) {
	ctx.Set(ctxKey, Guards{
		app.guard: {claims, token},
	})
}

func tokenIsDisabled(token string) bool {
	return facades.Cache.GetBool(getDisabledCacheKey(token), false)
}

func getDisabledCacheKey(token string) string {
	return "jwt:disabled:" + token
}