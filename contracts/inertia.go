package contracts

import "github.com/goravel/framework/contracts/http"

type Inertia interface {
	Render(ctx http.Context, component string, props map[string]any) error
	Share(key string, value any)
	ShareFunc(key string, fn func(ctx http.Context) any)
	Flash(ctx http.Context, key string, value any)
	Location(ctx http.Context, url string) error
	Version() string
	URL() string
}
