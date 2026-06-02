package goravelinertia

import (
	stdhttp "net/http"

	contractshttp "github.com/goravel/framework/contracts/http"
	inertia "github.com/petaki/inertia-go"
)

type Adapter struct {
	inertia *inertia.Inertia
}

func NewAdapter(inertia *inertia.Inertia) *Adapter {
	return &Adapter{inertia: inertia}
}

func (a *Adapter) Writer(ctx contractshttp.Context) stdhttp.ResponseWriter {
	return ctx.Response().Writer()
}

func (a *Adapter) Request(ctx contractshttp.Context) *stdhttp.Request {
	return ctx.Request().Origin()
}

func (a *Adapter) Inertia() *inertia.Inertia {
	return a.inertia
}
