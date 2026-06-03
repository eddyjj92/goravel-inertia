package console

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeGoMod(t *testing.T, module string) {
	t.Helper()
	body := "module " + module + "\n\ngo 1.23\n"
	if err := os.WriteFile("go.mod", []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestModuleName(t *testing.T) {
	chdirTemp(t)
	writeGoMod(t, "example.com/myapp")

	got, err := moduleName()
	if err != nil {
		t.Fatalf("moduleName() error = %v", err)
	}
	if got != "example.com/myapp" {
		t.Errorf("moduleName() = %q, want example.com/myapp", got)
	}
}

func TestInstallWebRoutesCreatesAndTemplatesModule(t *testing.T) {
	chdirTemp(t)

	status, err := installWebRoutes("example.com/myapp", false)
	if err != nil {
		t.Fatalf("installWebRoutes() error = %v", err)
	}
	if !strings.HasPrefix(status, "created") {
		t.Errorf("status = %q, want created…", status)
	}

	data, err := os.ReadFile(webRoutesPath)
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)
	if strings.Contains(content, modulePlaceholder) {
		t.Error("module placeholder was not replaced")
	}
	if !strings.Contains(content, "example.com/myapp/app/http/controllers") {
		t.Errorf("controllers import not templated with module, got:\n%s", content)
	}
}

func TestInstallWebRoutesReplacesDefault(t *testing.T) {
	chdirTemp(t)
	if err := os.MkdirAll("routes", 0o755); err != nil {
		t.Fatal(err)
	}
	// A default skeleton web.go references welcome.tmpl.
	def := "package routes\n\nfunc Web() { /* welcome.tmpl */ }\n"
	if err := os.WriteFile(webRoutesPath, []byte(def), 0o644); err != nil {
		t.Fatal(err)
	}

	status, err := installWebRoutes("m/app", false)
	if err != nil {
		t.Fatalf("installWebRoutes() error = %v", err)
	}
	if !strings.HasPrefix(status, "replaced") {
		t.Errorf("status = %q, want replaced…", status)
	}

	bak, err := os.ReadFile(webRoutesPath + ".bak")
	if err != nil {
		t.Fatalf("expected backup: %v", err)
	}
	if string(bak) != def {
		t.Error("backup does not match original web.go")
	}

	cur, _ := os.ReadFile(webRoutesPath)
	if !strings.Contains(string(cur), "NewHomeController") {
		t.Error("replaced web.go missing Inertia routes")
	}
}

func TestInstallWebRoutesSkipsCustomized(t *testing.T) {
	chdirTemp(t)
	if err := os.MkdirAll("routes", 0o755); err != nil {
		t.Fatal(err)
	}
	custom := "package routes\n\nfunc Web() { /* my custom routes */ }\n"
	if err := os.WriteFile(webRoutesPath, []byte(custom), 0o644); err != nil {
		t.Fatal(err)
	}

	status, err := installWebRoutes("m/app", false)
	if err != nil {
		t.Fatalf("installWebRoutes() error = %v", err)
	}
	if !strings.HasPrefix(status, "skipped") {
		t.Errorf("status = %q, want skipped…", status)
	}

	cur, _ := os.ReadFile(webRoutesPath)
	if string(cur) != custom {
		t.Error("customized web.go should be left untouched")
	}
	if _, err := os.Stat("routes/web.inertia.go.txt"); err != nil {
		t.Errorf("expected generated alternative file: %v", err)
	}
}

func TestInstallWebRoutesForceOverwritesCustomized(t *testing.T) {
	chdirTemp(t)
	if err := os.MkdirAll("routes", 0o755); err != nil {
		t.Fatal(err)
	}
	custom := "package routes\n\nfunc Web() { /* mine */ }\n"
	if err := os.WriteFile(webRoutesPath, []byte(custom), 0o644); err != nil {
		t.Fatal(err)
	}

	status, err := installWebRoutes("m/app", true)
	if err != nil {
		t.Fatalf("installWebRoutes() error = %v", err)
	}
	if !strings.HasPrefix(status, "replaced") {
		t.Errorf("status = %q, want replaced… with force", status)
	}
}

func TestRemoveWelcome(t *testing.T) {
	chdirTemp(t)

	// Nothing to remove → false, no error.
	if removed, err := removeWelcome(); err != nil || removed {
		t.Errorf("removeWelcome() on empty = (%v, %v), want (false, nil)", removed, err)
	}

	if err := os.MkdirAll(filepath.Dir(welcomePath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(welcomePath, []byte("<html></html>"), 0o644); err != nil {
		t.Fatal(err)
	}

	removed, err := removeWelcome()
	if err != nil || !removed {
		t.Fatalf("removeWelcome() = (%v, %v), want (true, nil)", removed, err)
	}
	if _, err := os.Stat(welcomePath); err == nil {
		t.Error("welcome.tmpl should be gone")
	}
}
