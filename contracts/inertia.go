package contracts

import "github.com/goravel/framework/contracts/http"

type Inertia interface {
	Render(ctx http.Context, component string, props map[string]any) http.Response
	Share(key string, value any)
	ShareFunc(key string, fn func(ctx http.Context) any)
	Defer(ctx http.Context, key string, fn func() any, group ...string)
	Location(ctx http.Context, url string) http.Response
	Version() string
	URL() string
}
