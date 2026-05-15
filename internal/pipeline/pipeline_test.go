package pipeline_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/pipeline"
)

func writeTempLog(t *testing.T, lines []string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "test-*.log")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	for _, l := range lines {
		f.WriteString(l + "\n")
	}
	f.Close()
	return f.Name()
}

func TestNew_NilWriter(t *testing.T) {
	_, err := pipeline.New(pipeline.Config{
		Patterns: []string{"*.log"},
		Writer:   nil,
	})
	if err == nil {
		t.Fatal("expected error for nil writer")
	}
}

func TestNew_NoPatterns(t *testing.T) {
	_, err := pipeline.New(pipeline.Config{
		Patterns: nil,
		Writer:   &bytes.Buffer{},
	})
	if err == nil {
		t.Fatal("expected error for empty patterns")
	}
}

func TestRun_FiltersByWindow(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	inWindow := now.Format(time.RFC3339)
	beforeWindow := now.Add(-2 * time.Hour).Format(time.RFC3339)

	lines := []string{
		`{"time":"` + beforeWindow + `","msg":"old"}`,
		`{"time":"` + inWindow + `","msg":"current"}`,
	}
	logFile := writeTempLog(t, lines)

	from := now.Add(-1 * time.Hour).Format(time.RFC3339)
	to := now.Add(1 * time.Hour).Format(time.RFC3339)

	var buf bytes.Buffer
	p, err := pipeline.New(pipeline.Config{
		Patterns: []string{logFile},
		From:     from,
		To:       to,
		Format:   "json",
		Writer:   &buf,
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	n, err := p.Run()
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if n != 1 {
		t.Fatalf("expected 1 entry, got %d", n)
	}

	var entry map[string]interface{}
	if err := json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &entry); err != nil {
		t.Fatalf("unmarshal output: %v", err)
	}
	if entry["msg"] != "current" {
		t.Errorf("unexpected msg: %v", entry["msg"])
	}
}

func TestRun_InvalidTimeWindow(t *testing.T) {
	var buf bytes.Buffer
	p, err := pipeline.New(pipeline.Config{
		Patterns: []string{"*.log"},
		From:     "not-a-time",
		To:       "also-not",
		Format:   "json",
		Writer:   &buf,
	})
	if err != nil {
		t.Fatalf("New should not fail here: %v", err)
	}
	_, err = p.Run()
	if err == nil {
		t.Fatal("expected error for invalid time window")
	}
}

func TestRun_GlobPattern(t *testing.T) {
	dir := t.TempDir()
	now := time.Now().UTC().Truncate(time.Second)
	ts := now.Format(time.RFC3339)

	for i := 0; i < 3; i++ {
		path := filepath.Join(dir, "app"+string(rune('a'+i))+".log")
		line := `{"time":"` + ts + `","msg":"entry"}` + "\n"
		os.WriteFile(path, []byte(line), 0644)
	}

	var buf bytes.Buffer
	p, err := pipeline.New(pipeline.Config{
		Patterns: []string{filepath.Join(dir, "*.log")},
		From:     now.Add(-time.Minute).Format(time.RFC3339),
		To:       now.Add(time.Minute).Format(time.RFC3339),
		Format:   "json",
		Writer:   &buf,
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	n, err := p.Run()
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if n != 3 {
		t.Errorf("expected 3 entries, got %d", n)
	}
}
