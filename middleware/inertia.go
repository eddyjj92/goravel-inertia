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

		inertia := facades.Inertia()
		version := inertia.Version()
		if ctx.Request().Method() == "GET" && ctx.Request().Header("X-Inertia-Version") != version {
			w := ctx.Response().Writer()
			w.Header().Set("X-Inertia-Location", ctx.Request().FullUrl())
			w.WriteHeader(stdhttp.StatusConflict)
			ctx.Request().Abort(stdhttp.StatusConflict)
			return
		}

		ctx.Request().Next()
	}
}
