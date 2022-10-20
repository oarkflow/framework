package middleware

import (
	"github.com/sujit-baniya/framework/contracts/http"
)

func Cors() http.HandlerFunc {
	return func(ctx http.Context) error {
		method := ctx.Method()
		origin := ctx.Header("Origin", "")
		if origin != "" {
			ctx.Response().Header("Access-Control-Allow-Origin", "*")
			ctx.Response().Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			ctx.Response().Header("Access-Control-Allow-Headers", "*")
			ctx.Response().Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Authorization")
			ctx.Response().Header("Access-Control-Max-Age", "172800")
			ctx.Response().Header("Access-Control-Allow-Credentials", "true")
		}

		if method == "OPTIONS" {
			ctx.AbortWithStatus(204)
			return nil
		}

		return ctx.Next()
	}
}
