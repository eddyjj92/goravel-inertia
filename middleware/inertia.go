package middleware

import (
	stdhttp "net/http"

	contractshttp "github.com/goravel/framework/contracts/http"

	"github.com/eddyjj92/goravel-inertia/facades"
)

func Inertia() contractshttp.Middleware {
	return func(ctx contractshttp.Context) {
		if ctx.Request().Header("X-Inertia") == "" {
			ctx.Request().Next()
			return
		}

		version := facades.Inertia().Version()
		if ctx.Request().Method() == stdhttp.MethodGet && ctx.Request().Header("X-Inertia-Version") != version {
			ctx.Response().Header("X-Inertia-Location", ctx.Request().FullUrl())
			ctx.Request().Abort(stdhttp.StatusConflict)
			return
		}

		ctx.Request().Next()
	}
}
