package auth

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/oarkflow/frame"

	authContract "github.com/oarkflow/framework/contracts/auth"
	"github.com/oarkflow/framework/facades"
	supporttime "github.com/oarkflow/framework/support/time"
)

var (
	unit = time.Minute

	ErrorRefreshTimeExceeded = errors.New("refresh time exceeded")
	ErrorNoPrimaryKeyField   = errors.New("the primaryKey field was not found in the model, set primaryKey like orm.Model")
	ErrorEmptySecret         = errors.New("secret is required")
	ErrorTokenDisabled       = errors.New("token is disabled")
	ErrorParseTokenFirst     = errors.New("parse token first")
)

type Jwt struct {
	guard string
}

func NewJwt(guard string) authContract.Auth {
	return &Jwt{
		guard: guard,
	}
}

func (app *Jwt) Guard(name string) authContract.Auth {
	return NewJwt(name)
}

// User need parse token first.
func (app *Jwt) User(ctx *frame.Context, user authContract.User) error {
	val := ctx.Value(ctx.AuthUserKey)
	if val == nil {
		return ErrorParseTokenFirst
	}
	ctx.Set(ctx.AuthUserKey, val)
	return nil
}

func (app *Jwt) Parse(ctx *frame.Context, token string, user authContract.User) error {
	token = strings.ReplaceAll(token, "Bearer ", "")
	if app.tokenIsDisabled(token) {
		return ErrorTokenDisabled
	}
	var claims jwt.RegisteredClaims
	secret := facades.Config.GetString("jwt.secret")
	tokenClaims, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (any, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return err
	}
	subject, err := tokenClaims.Claims.GetSubject()
	if err != nil {
		return err
	}
	if err := facades.Orm.Query(facades.Config.GetString("database.default")).Find(user, subject).Error; err != nil {
		return err
	}
	ctx.Set(ctx.AuthUserKey, user)
	ctx.Set("token_claim", claims)
	ctx.Set("access_token", token)
	app.Refresh(ctx)
	return nil
}

func (app *Jwt) Data(ctx *frame.Context) (map[string]any, error) {
	return nil, nil
}

func (app *Jwt) Login(ctx *frame.Context, user authContract.User, data ...map[string]any) (token string, err error) {
	t := reflect.TypeOf(user).Elem()
	v := reflect.ValueOf(user).Elem()
	fmt.Println(user)
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

func (app *Jwt) LoginUsingID(ctx *frame.Context, id any) (token string, err error) {
	secret := facades.Config.GetString("jwt.secret")
	if secret == "" {
		return "", ErrorEmptySecret
	}

	nowTime := supporttime.Now()
	ttl := facades.Config.GetInt("jwt.ttl")
	expireTime := nowTime.Add(time.Duration(ttl) * unit)
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expireTime),
		IssuedAt:  jwt.NewNumericDate(nowTime),
		Subject:   fmt.Sprintf("%v", id),
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tokenClaims.SignedString([]byte(secret))
}

// Refresh need parse token first.
func (app *Jwt) Refresh(ctx *frame.Context) (token string, err error) {
	val := ctx.Value("token_claim")
	if val == nil {
		return "", ErrorParseTokenFirst
	}
	claim := val.(jwt.RegisteredClaims)

	nowTime := supporttime.Now()
	refreshTtl := facades.Config.GetInt("jwt.refresh_ttl")
	expireTime := claim.ExpiresAt.Add(time.Duration(refreshTtl) * unit)
	if nowTime.Unix() > expireTime.Unix() {
		return "", ErrorRefreshTimeExceeded
	}

	return app.LoginUsingID(ctx, claim.Subject)
}

func (app *Jwt) Logout(ctx *frame.Context) error {

	token, ok := ctx.Value(ctx.AuthUserKey).(string)
	if !ok || token == "" {
		return nil
	}

	if facades.Cache == nil {
		return errors.New("cache support is required")
	}

	if err := facades.Cache.Put(app.getDisabledCacheKey(token),
		true,
		time.Duration(facades.Config.GetInt("jwt.ttl"))*unit,
	); err != nil {
		return err
	}
	ctx.Set(ctx.AuthUserKey, nil)

	return nil
}

func (app *Jwt) tokenIsDisabled(token string) bool {
	return facades.Cache.GetBool(app.getDisabledCacheKey(token), false)
}

func (app *Jwt) getDisabledCacheKey(token string) string {
	return "paseto:disabled:" + token
}
