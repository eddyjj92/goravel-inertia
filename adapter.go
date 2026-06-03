package goravelinertia

import (
	stdhttp "net/http"

	contractshttp "github.com/goravel/framework/contracts/http"
	inertia "github.com/petaki/inertia-go"
)

type Adapter struct {
	inertia *inertia.Inertia
	csr     *inertia.Inertia
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

// CSR returns the SSR-disabled fallback engine, or nil when SSR is off. It is set
// once at boot and read-only afterwards.
func (a *Adapter) CSR() *inertia.Inertia {
	return a.csr
}

// SetCSR registers the SSR-disabled fallback engine used when an SSR render fails.
func (a *Adapter) SetCSR(i *inertia.Inertia) {
	a.csr = i
}
