package checkpoint_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/example/logslice/internal/checkpoint"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "checkpoint.json")
}

func TestNew_NonExistentFile(t *testing.T) {
	c, err := checkpoint.New(tempPath(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := c.Offset("/var/log/app.log"); got != 0 {
		t.Errorf("expected 0, got %d", got)
	}
}

func TestSetOffset_And_Offset(t *testing.T) {
	c, _ := checkpoint.New(tempPath(t))
	c.SetOffset("/var/log/app.log", 1024)
	if got := c.Offset("/var/log/app.log"); got != 1024 {
		t.Errorf("expected 1024, got %d", got)
	}
}

func TestSave_And_Reload(t *testing.T) {
	path := tempPath(t)
	c, err := checkpoint.New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	c.SetOffset("/logs/a.log", 512)
	c.SetOffset("/logs/b.log", 2048)
	if err := c.Save(); err != nil {
		t.Fatalf("Save: %v", err)
	}

	c2, err := checkpoint.New(path)
	if err != nil {
		t.Fatalf("reload New: %v", err)
	}
	if got := c2.Offset("/logs/a.log"); got != 512 {
		t.Errorf("/logs/a.log: expected 512, got %d", got)
	}
	if got := c2.Offset("/logs/b.log"); got != 2048 {
		t.Errorf("/logs/b.log: expected 2048, got %d", got)
	}
}

func TestSave_CreatesFileAtomically(t *testing.T) {
	path := tempPath(t)
	c, _ := checkpoint.New(path)
	c.SetOffset("x", 99)
	if err := c.Save(); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("checkpoint file missing: %v", err)
	}
	if _, err := os.Stat(path + ".tmp"); !os.IsNotExist(err) {
		t.Errorf("temp file should have been removed")
	}
}

func TestNew_CorruptFile_ReturnsError(t *testing.T) {
	path := tempPath(t)
	if err := os.WriteFile(path, []byte("not json{"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := checkpoint.New(path)
	if err == nil {
		t.Error("expected error for corrupt checkpoint file")
	}
}
