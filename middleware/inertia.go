package middleware

import (
	stdhttp "net/http"

	contractshttp "github.com/goravel/framework/contracts/http"

	"github.com/eddyjj92/goravel-inertia/facades"
)

// Inertia returns the middleware that drives the Inertia protocol. Register it on
// the web route group so every Inertia request is handled:
//
//	facades.Route().Middleware(middleware.Inertia()).Group(func(r route.Router) { ... })
//
// Non-Inertia requests pass straight through. For Inertia GET requests whose asset
// version is stale it replies 409 Conflict + X-Inertia-Location so the client does a
// full reload. Per-request v3 props (Defer/Optional/Merge/...) need no setup here:
// they lazily initialise from the request context the first time one is set.
func Inertia() contractshttp.Middleware {
	return func(ctx contractshttp.Context) {
		inertia := facades.Inertia()

		if ctx.Request().Header("X-Inertia") != "" &&
			ctx.Request().Method() == stdhttp.MethodGet &&
			ctx.Request().Header("X-Inertia-Version") != inertia.Version() {
			ctx.Response().Header("X-Inertia-Location", ctx.Request().FullUrl())
			ctx.Request().Abort(stdhttp.StatusConflict)
			return
		}

		// Mirror session flash + validation errors into props for both the initial
		// HTML load and X-Inertia visits.
		inertia.ShareSession(ctx)

		ctx.Request().Next()
	}
}
