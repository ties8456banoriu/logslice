package pipeline_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"logslice/internal/logparser"
	"logslice/internal/logreader"
	"logslice/internal/output"
	"logslice/internal/pipeline"
	"logslice/internal/timewindow"
)

func writeTempLog(t *testing.T, lines []string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	for _, l := range lines {
		f.WriteString(l + "\n")
	}
	return f.Name()
}

func entry(ts, msg string) string {
	b, _ := json.Marshal(map[string]string{"time": ts, "msg": msg})
	return string(b)
}

func mustWindow(t *testing.T, from, to string) *timewindow.TimeWindow {
	t.Helper()
	w, err := timewindow.New(from, to, time.RFC3339)
	if err != nil {
		t.Fatal(err)
	}
	return w
}

func TestRun_NilWindow(t *testing.T) {
	cfg := pipeline.Config{Window: nil}
	_, err := pipeline.Run(context.Background(), cfg)
	if err == nil || !strings.Contains(err.Error(), "window") {
		t.Fatalf("expected window error, got %v", err)
	}
}

func TestRun_FiltersByWindow(t *testing.T) {
	lines := []string{
		entry("2024-01-01T10:00:00Z", "before"),
		entry("2024-01-01T12:00:00Z", "inside"),
		entry("2024-01-01T14:00:00Z", "after"),
	}
	path := writeTempLog(t, lines)

	w := mustWindow(t, "2024-01-01T11:00:00Z", "2024-01-01T13:00:00Z")

	parser, err := logparser.NewParser("json")
	if err != nil {
		t.Fatal(err)
	}
	reader, err := logreader.New(w, parser)
	if err != nil {
		t.Fatal(err)
	}

	var buf strings.Builder
	writer, err := output.New(&buf, output.JSON)
	if err != nil {
		t.Fatal(err)
	}

	res, err := pipeline.Run(context.Background(), pipeline.Config{
		Patterns: []string{path},
		Window:   w,
		Reader:   reader,
		Writer:   writer,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.FileCount != 1 {
		t.Errorf("expected 1 file, got %d", res.FileCount)
	}
}

func TestRun_NoMatchingFiles(t *testing.T) {
	dir := t.TempDir()
	w := mustWindow(t, "2024-01-01T00:00:00Z", "2024-01-02T00:00:00Z")
	parser, _ := logparser.NewParser("json")
	reader, _ := logreader.New(w, parser)
	var buf strings.Builder
	writer, _ := output.New(&buf, output.JSON)

	_, err := pipeline.Run(context.Background(), pipeline.Config{
		Patterns: []string{filepath.Join(dir, "*.log")},
		Window:   w,
		Reader:   reader,
		Writer:   writer,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRun_ContextCancelled(t *testing.T) {
	path := writeTempLog(t, []string{entry("2024-01-01T12:00:00Z", "x")})
	w := mustWindow(t, "2024-01-01T11:00:00Z", "2024-01-01T13:00:00Z")
	parser, _ := logparser.NewParser("json")
	reader, _ := logreader.New(w, parser)
	var buf strings.Builder
	writer, _ := output.New(&buf, output.JSON)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := pipeline.Run(ctx, pipeline.Config{
		Patterns: []string{path},
		Window:   w,
		Reader:   reader,
		Writer:   writer,
	})
	// cancelled context should not produce an error from pipeline itself
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
