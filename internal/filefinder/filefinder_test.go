package filefinder_test

import (
	"os"
	"path/filepath"
	"testing"

	"logslice/internal/filefinder"
)

func TestExpand_LiteralFile(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "app.log")
	mustWrite(t, f, "line")

	finder := filefinder.New()
	got, err := finder.Expand([]string{f})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 || got[0] != f {
		t.Errorf("expected [%s], got %v", f, got)
	}
}

func TestExpand_GlobPattern(t *testing.T) {
	tmp := t.TempDir()
	for _, name := range []string{"a.log", "b.log", "c.txt"} {
		mustWrite(t, filepath.Join(tmp, name), "")
	}

	finder := filefinder.New()
	got, err := finder.Expand([]string{filepath.Join(tmp, "*.log")})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("expected 2 files, got %d: %v", len(got), got)
	}
}

func TestExpand_DirectoryWithoutRecursive(t *testing.T) {
	tmp := t.TempDir()
	mustWrite(t, filepath.Join(tmp, "x.log"), "")

	finder := filefinder.New()
	_, err := finder.Expand([]string{tmp})
	if err == nil {
		t.Fatal("expected error for directory without recursive flag")
	}
}

func TestExpand_DirectoryRecursive(t *testing.T) {
	tmp := t.TempDir()
	sub := filepath.Join(tmp, "sub")
	if err := os.Mkdir(sub, 0755); err != nil {
		t.Fatal(err)
	}
	mustWrite(t, filepath.Join(tmp, "a.log"), "")
	mustWrite(t, filepath.Join(sub, "b.log"), "")

	finder := filefinder.New(filefinder.WithRecursive(true))
	got, err := finder.Expand([]string{tmp})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("expected 2 files, got %d: %v", len(got), got)
	}
}

func TestExpand_Deduplication(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "dup.log")
	mustWrite(t, f, "")

	finder := filefinder.New()
	got, err := finder.Expand([]string{f, f})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("expected 1 file after dedup, got %d", len(got))
	}
}

func TestExpand_MissingFile(t *testing.T) {
	finder := filefinder.New()
	_, err := finder.Expand([]string{"/nonexistent/path/file.log"})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func mustWrite(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("mustWrite: %v", err)
	}
}
