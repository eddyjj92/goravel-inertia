package providers

import (
	"github.com/goravel/framework/contracts/foundation"
)

type InertiaServiceProvider struct {
}

func (p *InertiaServiceProvider) Register(app foundation.Application) {
}

func (p *InertiaServiceProvider) Boot(app foundation.Application) {
}

func (p *InertiaServiceProvider) Start(app foundation.Application) error {
	return nil
}
