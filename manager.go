package goravelinertia

import (
	"maps"
	stdhttp "net/http"
	"sync"

	contractshttp "github.com/goravel/framework/contracts/http"
)

// defaultFlashKeys are the session keys mirrored into props.flash when no
// flash_keys are configured.
var defaultFlashKeys = []string{"success", "error", "warning", "info", "message"}

type InertiaManager struct {
	mu          sync.RWMutex
	adapter     *Adapter
	url         string
	version     string
	flashKeys   []string
	sharedFuncs map[string]func(contractshttp.Context) any
}

func NewInertiaManager(adapter *Adapter, url string, version string, flashKeys ...string) *InertiaManager {
	if len(flashKeys) == 0 {
		flashKeys = defaultFlashKeys
	}

	return &InertiaManager{
		adapter:     adapter,
		url:         url,
		version:     version,
		flashKeys:   flashKeys,
		sharedFuncs: make(map[string]func(contractshttp.Context) any),
	}
}

func (m *InertiaManager) Render(ctx contractshttp.Context, component string, props map[string]any) contractshttp.Response {
	return newResponse(func() error {
		m.mu.RLock()
		defer m.mu.RUnlock()

		evaluated := make(map[string]any, len(m.sharedFuncs))
		for key, fn := range m.sharedFuncs {
			evaluated[key] = fn(ctx)
		}

		maps.Copy(evaluated, props)

		w := m.adapter.Writer(ctx)
		r := m.adapter.Request(ctx).WithContext(m.propsContext(ctx))

		return m.adapter.Inertia().Render(w, r, component, evaluated)
	})
}

func (m *InertiaManager) Share(key string, value any) {
	m.adapter.Inertia().Share(key, value)
}

func (m *InertiaManager) ShareFunc(key string, fn func(contractshttp.Context) any) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.sharedFuncs[key] = fn
}

// redirectStatus returns the HTTP status an Inertia redirect must use. Mutating
// methods (PUT/PATCH/DELETE) require 303 See Other so the browser follows the
// redirect with a GET instead of replaying the original method; everything else
// uses 302 Found. Mirrors the behaviour of inertiajs/inertia-laravel.
func redirectStatus(method string) int {
	switch method {
	case stdhttp.MethodPut, stdhttp.MethodPatch, stdhttp.MethodDelete:
		return stdhttp.StatusSeeOther
	default:
		return stdhttp.StatusFound
	}
}

// Redirect issues an Inertia-aware internal redirect, picking 303 for mutating
// requests so the client re-fetches the target with a GET.
func (m *InertiaManager) Redirect(ctx contractshttp.Context, url string) contractshttp.Response {
	return ctx.Response().Redirect(redirectStatus(ctx.Request().Method()), url)
}

func (m *InertiaManager) Location(ctx contractshttp.Context, url string) contractshttp.Response {
	return newResponse(func() error {
		w := m.adapter.Writer(ctx)
		r := m.adapter.Request(ctx)

		m.adapter.Inertia().Location(w, r, url)

		return nil
	})
}

func (m *InertiaManager) Version() string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.version
}

func (m *InertiaManager) URL() string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.url
}

func (m *InertiaManager) GetAdapter() *Adapter {
	return m.adapter
}
