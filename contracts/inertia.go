package contracts

import "github.com/goravel/framework/contracts/http"

type Inertia interface {
	Render(ctx http.Context, component string, props map[string]any) http.Response
	Share(key string, value any)
	ShareFunc(key string, fn func(ctx http.Context) any)
	Location(ctx http.Context, url string) error
	Version() string
	URL() string
}
