// Package providers wires the Inertia manager into the Goravel service container.
package providers

import (
	"time"

	"github.com/goravel/framework/contracts/foundation"
	contractshttp "github.com/goravel/framework/contracts/http"

	petaki "github.com/petaki/inertia-go"

	goravelinertia "github.com/eddyjj92/goravel-inertia"
	"github.com/eddyjj92/goravel-inertia/contracts"
	"github.com/eddyjj92/goravel-inertia/facades"
)

// InertiaServiceProvider registers and boots the Inertia singleton and facade.
type InertiaServiceProvider struct {
}

// toStringSlice converts a config value (which may be []string or []any) into a
// []string, returning nil so the manager falls back to its default flash keys.
func toStringSlice(v any) []string {
	switch s := v.(type) {
	case []string:
		return s
	case []any:
		out := make([]string, 0, len(s))
		for _, item := range s {
			if str, ok := item.(string); ok {
				out = append(out, str)
			}
		}
		return out
	default:
		return nil
	}
}

// Register binds the Inertia manager as the "goravel.inertia" singleton.
func (p *InertiaServiceProvider) Register(app foundation.Application) {
	app.Singleton("goravel.inertia", func(app foundation.Application) (any, error) {
		config := app.MakeConfig()

		rootView := "app"
		version := ""
		ssr := false
		ssrURL := "http://127.0.0.1:13714/render"
		url := ""
		var flashKeys []string

		if config != nil {
			rootView = config.GetString("inertia.root_view", "app")
			version = config.GetString("inertia.version", "")
			ssr = config.GetBool("inertia.ssr", false)
			ssrURL = config.GetString("inertia.ssr_url", "http://127.0.0.1:13714/render")
			url = config.GetString("app.url", "")
			flashKeys = toStringSlice(config.Get("inertia.flash_keys"))
		}

		inertia := petaki.New(url, rootView, version)

		if ssr {
			inertia.EnableSsr(ssrURL)
		}

		adapter := goravelinertia.NewAdapter(inertia)
		manager := goravelinertia.NewInertiaManager(adapter, url, version, flashKeys...)

		return manager, nil
	})
}

// Boot resolves the manager, registers default shared props, and exposes the facade.
func (p *InertiaServiceProvider) Boot(app foundation.Application) {
	instance, err := app.Make("goravel.inertia")
	if err != nil {
		panic("Failed to resolve goravel.inertia: " + err.Error())
	}

	inertia := instance.(contracts.Inertia)

	config := app.MakeConfig()
	if config != nil {
		inertia.Share("appName", config.GetString("app.name"))
	}

	inertia.ShareFunc("timestamp", func(_ contractshttp.Context) any {
		return time.Now().Format("2006-01-02 15:04:05")
	})

	facades.RegisterInertia(inertia)
}
