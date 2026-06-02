package goravelinertia

import (
	"sync"

	contractshttp "github.com/goravel/framework/contracts/http"
)

type InertiaManager struct {
	mu           sync.RWMutex
	adapter      *Adapter
	url          string
	version      string
	sharedFuncs  map[string]func(contractshttp.Context) any
}

var flashCtxKey = "goravel_inertia_flash"

func NewInertiaManager(adapter *Adapter, url string, version string) *InertiaManager {
	return &InertiaManager{
		adapter:     adapter,
		url:         url,
		version:     version,
		sharedFuncs: make(map[string]func(contractshttp.Context) any),
	}
}

func (m *InertiaManager) Render(ctx contractshttp.Context, component string, props map[string]any) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	evaluated := make(map[string]any, len(m.sharedFuncs))
	for key, fn := range m.sharedFuncs {
		evaluated[key] = fn(ctx)
	}

	for k, v := range props {
		evaluated[k] = v
	}

	if flashData, ok := ctx.Value(flashCtxKey).(map[string]any); ok && len(flashData) > 0 {
		evaluated["flash"] = flashData
	}

	w := m.adapter.Writer(ctx)
	r := m.adapter.Request(ctx)

	return m.adapter.Inertia().Render(w, r, component, evaluated)
}

func (m *InertiaManager) Share(key string, value any) {
	m.adapter.Inertia().Share(key, value)
}

func (m *InertiaManager) ShareFunc(key string, fn func(contractshttp.Context) any) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.sharedFuncs[key] = fn
}

func (m *InertiaManager) Flash(ctx contractshttp.Context, key string, value any) {
	existing, ok := ctx.Value(flashCtxKey).(map[string]any)
	if !ok || existing == nil {
		existing = make(map[string]any)
	}
	existing[key] = value
	ctx.WithValue(flashCtxKey, existing)
}

func (m *InertiaManager) Location(ctx contractshttp.Context, url string) error {
	w := m.adapter.Writer(ctx)
	r := m.adapter.Request(ctx)

	m.adapter.Inertia().Location(w, r, url)

	return nil
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
