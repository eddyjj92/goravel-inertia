package middleware

import contractshttp "github.com/goravel/framework/contracts/http"

func Inertia() contractshttp.Middleware {
	return func(ctx contractshttp.Context) {
		ctx.Request().Next()
	}
}
