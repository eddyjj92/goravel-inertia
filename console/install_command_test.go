package console

import (
	"os"
	"path/filepath"
	"testing"
)

// chdirTemp switches into a fresh temp dir for the duration of the test.
func chdirTemp(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	prev, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(prev) })
	return dir
}

func TestScaffoldCreatesAllFiles(t *testing.T) {
	dir := chdirTemp(t)

	created, skipped, err := scaffold(false)
	if err != nil {
		t.Fatalf("scaffold() error = %v", err)
	}
	if len(created) != len(fileMap) {
		t.Errorf("created %d files, want %d", len(created), len(fileMap))
	}
	if len(skipped) != 0 {
		t.Errorf("skipped = %v, want none on a clean dir", skipped)
	}

	for _, dst := range fileMap {
		full := filepath.Join(dir, dst)
		info, statErr := os.Stat(full)
		if statErr != nil {
			t.Errorf("expected %s to exist: %v", dst, statErr)
			continue
		}
		if info.Size() == 0 {
			t.Errorf("%s is empty", dst)
		}
	}
}

func TestScaffoldSkipsExistingWithoutForce(t *testing.T) {
	chdirTemp(t)

	if _, _, err := scaffold(false); err != nil {
		t.Fatalf("first scaffold() error = %v", err)
	}

	created, skipped, err := scaffold(false)
	if err != nil {
		t.Fatalf("second scaffold() error = %v", err)
	}
	if len(created) != 0 {
		t.Errorf("created = %v, want none on rerun without force", created)
	}
	if len(skipped) != len(fileMap) {
		t.Errorf("skipped %d files, want %d", len(skipped), len(fileMap))
	}
}

func TestScaffoldForceOverwrites(t *testing.T) {
	chdirTemp(t)

	if _, _, err := scaffold(false); err != nil {
		t.Fatalf("first scaffold() error = %v", err)
	}

	// Mutate a file, then force-overwrite and confirm it was restored.
	target := fileMap["stubs/vite.config.ts.stub"]
	if err := os.WriteFile(target, []byte("// mutated"), 0o644); err != nil {
		t.Fatal(err)
	}

	created, skipped, err := scaffold(true)
	if err != nil {
		t.Fatalf("force scaffold() error = %v", err)
	}
	if len(created) != len(fileMap) {
		t.Errorf("created %d files, want %d with force", len(created), len(fileMap))
	}
	if len(skipped) != 0 {
		t.Errorf("skipped = %v, want none with force", skipped)
	}

	data, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) == "// mutated" {
		t.Error("force did not overwrite the mutated file")
	}
}
