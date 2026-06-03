package goravelinertia

import (
	"context"

	contractshttp "github.com/goravel/framework/contracts/http"
)

// contextKey namespaces values stored on the Goravel http.Context.
type contextKey string

// ctxKeyProps holds the accumulated context.Context that carries the Inertia v3
// per-request props (deferred, optional, merge, ...). petaki reads these from the
// *http.Request context, so the accumulated context is injected into the request
// just before Render runs.
const ctxKeyProps = contextKey("goravel-inertia.props")

// propsContext returns the accumulated v3 props context for this request, falling
// back to the underlying request context the first time a prop is set.
func (m *InertiaManager) propsContext(ctx contractshttp.Context) context.Context {
	if v := ctx.Value(ctxKeyProps); v != nil {
		if c, ok := v.(context.Context); ok {
			return c
		}
	}

	return m.adapter.Request(ctx).Context()
}

// storePropsContext persists the accumulated v3 props context back onto the
// Goravel http.Context so chained With* calls and Render see the same context.
func (m *InertiaManager) storePropsContext(ctx contractshttp.Context, c context.Context) {
	ctx.WithValue(ctxKeyProps, c)
}

// Defer registers a prop evaluated lazily by the client after the initial load.
func (m *InertiaManager) Defer(ctx contractshttp.Context, key string, fn func() any, group ...string) {
	c := m.adapter.Inertia().WithDeferredProp(m.propsContext(ctx), key, fn, group...)
	m.storePropsContext(ctx, c)
}
