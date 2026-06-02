package goravelinertia

import (
	"testing"

	petaki "github.com/petaki/inertia-go"
)

func newTestManager(url, version string) *InertiaManager {
	adapter := NewAdapter(petaki.New(url, "app", version))
	return NewInertiaManager(adapter, url, version)
}

func TestManagerVersionAndURL(t *testing.T) {
	m := newTestManager("http://localhost", "v1")

	if got := m.Version(); got != "v1" {
		t.Errorf("Version() = %q, want %q", got, "v1")
	}
	if got := m.URL(); got != "http://localhost" {
		t.Errorf("URL() = %q, want %q", got, "http://localhost")
	}
}

func TestManagerAdapterWired(t *testing.T) {
	m := newTestManager("http://localhost", "v1")

	if m.GetAdapter() == nil {
		t.Fatal("GetAdapter() = nil, want adapter")
	}
	if m.GetAdapter().Inertia() == nil {
		t.Fatal("adapter.Inertia() = nil, want petaki instance")
	}
}
