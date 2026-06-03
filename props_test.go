package goravelinertia

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	contractshttp "github.com/goravel/framework/contracts/http"
)

// fakeContext implements just enough of contractshttp.Context for the manager.
// context.Context is not embedded anonymously because the interface also exposes
// a Context() method, which would collide with the embedded field name.
type fakeContext struct {
	base context.Context
	vals map[any]any
	req  *http.Request
	w    http.ResponseWriter
}

func newFakeContext(req *http.Request, w http.ResponseWriter) *fakeContext {
	return &fakeContext{base: context.Background(), vals: map[any]any{}, req: req, w: w}
}

func (f *fakeContext) Deadline() (time.Time, bool) { return f.base.Deadline() }
func (f *fakeContext) Done() <-chan struct{}       { return f.base.Done() }
func (f *fakeContext) Err() error                  { return f.base.Err() }
func (f *fakeContext) Value(key any) any           { return f.vals[key] }
func (f *fakeContext) Context() context.Context    { return f.base }
func (f *fakeContext) WithContext(c context.Context) { f.base = c }
func (f *fakeContext) WithValue(k, v any)             { f.vals[k] = v }
func (f *fakeContext) Request() contractshttp.ContextRequest   { return &fakeRequest{r: f.req} }
func (f *fakeContext) Response() contractshttp.ContextResponse { return &fakeResponse{w: f.w} }

type fakeRequest struct {
	contractshttp.ContextRequest
	r *http.Request
}

func (f *fakeRequest) Origin() *http.Request { return f.r }

type fakeResponse struct {
	contractshttp.ContextResponse
	w http.ResponseWriter
}

func (f *fakeResponse) Writer() http.ResponseWriter { return f.w }

func inertiaJSON(t *testing.T, m *InertiaManager, ctx contractshttp.Context, component string, props map[string]any) map[string]any {
	t.Helper()

	resp := m.Render(ctx, component, props)
	if err := resp.Render(); err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	rec, ok := ctx.Response().Writer().(*httptest.ResponseRecorder)
	if !ok {
		t.Fatal("writer is not a ResponseRecorder")
	}

	var page map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &page); err != nil {
		t.Fatalf("unmarshal page: %v (body=%q)", err, rec.Body.String())
	}

	return page
}

func newInertiaCtx() *fakeContext {
	req := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
	req.Header.Set("X-Inertia", "true") // force the JSON branch, no template needed
	return newFakeContext(req, httptest.NewRecorder())
}

func TestDeferThreadsIntoRenderedJSON(t *testing.T) {
	m := newTestManager("", "")
	ctx := newInertiaCtx()

	m.Defer(ctx, "stats", func() any { return map[string]any{"x": 1} }, "metrics")

	page := inertiaJSON(t, m, ctx, "Dashboard", map[string]any{"a": 1})

	deferred, ok := page["deferredProps"].(map[string]any)
	if !ok {
		t.Fatalf("deferredProps missing or wrong type: %#v", page["deferredProps"])
	}

	group, ok := deferred["metrics"].([]any)
	if !ok || len(group) != 1 || group[0] != "stats" {
		t.Fatalf("deferredProps[metrics] = %#v, want [stats]", deferred["metrics"])
	}

	// deferred prop must NOT be evaluated into props on the initial load.
	props, _ := page["props"].(map[string]any)
	if _, present := props["stats"]; present {
		t.Errorf("deferred prop 'stats' should be absent from initial props, got %#v", props)
	}
}

func TestRenderWithoutDeferIsClean(t *testing.T) {
	m := newTestManager("", "")
	ctx := newInertiaCtx()

	page := inertiaJSON(t, m, ctx, "Dashboard", map[string]any{"a": float64(1)})

	if _, present := page["deferredProps"]; present {
		t.Errorf("deferredProps should be absent when nothing deferred, got %#v", page["deferredProps"])
	}
	props, _ := page["props"].(map[string]any)
	if props["a"] != float64(1) {
		t.Errorf("props[a] = %#v, want 1", props["a"])
	}
}
