package providers

import (
	"time"

	contractshttp "github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/foundation"

	goravelinertia "github.com/eddyjj92/goravel-inertia"
	"github.com/eddyjj92/goravel-inertia/contracts"
	"github.com/eddyjj92/goravel-inertia/facades"
	petaki "github.com/petaki/inertia-go"
)

type InertiaServiceProvider struct {
}

func (p *InertiaServiceProvider) Register(app foundation.Application) {
	app.Singleton("goravel.inertia", func(app foundation.Application) (any, error) {
		config := app.MakeConfig()

		rootView := "app"
		version := ""
		ssr := false
		ssrURL := "http://127.0.0.1:13714/render"
		url := ""

		if config != nil {
			rootView = config.GetString("inertia.root_view", "app")
			version = config.GetString("inertia.version", "")
			ssr = config.GetBool("inertia.ssr", false)
			ssrURL = config.GetString("inertia.ssr_url", "http://127.0.0.1:13714/render")
			url = config.GetString("app.url", "")
		}

		inertia := petaki.New(url, rootView, version)

		if ssr {
			inertia.EnableSsr(ssrURL)
		}

		adapter := goravelinertia.NewAdapter(inertia)
		manager := goravelinertia.NewInertiaManager(adapter, url, version)

		return manager, nil
	})
}

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

	inertia.ShareFunc("timestamp", func(ctx contractshttp.Context) any {
		return time.Now().Format("2006-01-02 15:04:05")
	})

	facades.RegisterInertia(inertia)
}
