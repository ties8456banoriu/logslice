package logreader_test

import (
	"strings"
	"testing"
	"time"

	"github.com/user/logslice/internal/logparser"
	"github.com/user/logslice/internal/logreader"
	"github.com/user/logslice/internal/timewindow"
)

func mustWindow(t *testing.T, from, to string) *timewindow.TimeWindow {
	t.Helper()
	w, err := timewindow.New(from, to)
	if err != nil {
		t.Fatalf("timewindow.New: %v", err)
	}
	return w
}

func mustParser(t *testing.T) *logparser.JSONParser {
	t.Helper()
	p, err := logparser.NewJSONParser(logparser.Options{})
	if err != nil {
		t.Fatalf("NewJSONParser: %v", err)
	}
	return p
}

const sampleLogs = `{"time":"2024-01-15T10:00:00Z","level":"info","msg":"startup"}
{"time":"2024-01-15T10:05:00Z","level":"warn","msg":"high memory"}
{"time":"2024-01-15T10:10:00Z","level":"error","msg":"disk full"}
{"time":"2024-01-15T10:20:00Z","level":"info","msg":"recovered"}
not valid json
`

func TestRead_FiltersByWindow(t *testing.T) {
	w := mustWindow(t, "2024-01-15T10:04:00Z", "2024-01-15T10:11:00Z")
	r := logreader.New(mustParser(t), w)

	entries, err := r.Read(strings.NewReader(sampleLogs))
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Message != "high memory" {
		t.Errorf("entry[0].Message = %q, want %q", entries[0].Message, "high memory")
	}
	if entries[1].Message != "disk full" {
		t.Errorf("entry[1].Message = %q, want %q", entries[1].Message, "disk full")
	}
}

func TestRead_EmptySource(t *testing.T) {
	w := mustWindow(t, "2024-01-15T10:00:00Z", "2024-01-15T11:00:00Z")
	r := logreader.New(mustParser(t), w)

	entries, err := r.Read(strings.NewReader(""))
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestRead_NoMatchInWindow(t *testing.T) {
	w := mustWindow(t, "2024-01-15T12:00:00Z", "2024-01-15T13:00:00Z")
	r := logreader.New(mustParser(t), w)

	entries, err := r.Read(strings.NewReader(sampleLogs))
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestRead_SkipsInvalidLines(t *testing.T) {
	w := mustWindow(t, "2024-01-15T09:00:00Z", "2024-01-15T11:00:00Z")
	r := logreader.New(mustParser(t), w)

	// All valid entries should be returned; invalid JSON lines are skipped.
	entries, err := r.Read(strings.NewReader(sampleLogs))
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if len(entries) != 4 {
		t.Errorf("expected 4 entries, got %d", len(entries))
	}
}

func TestRead_BoundaryInclusion(t *testing.T) {
	// timewindow.Contains should include the boundary timestamps.
	start := "2024-01-15T10:00:00Z"
	end := "2024-01-15T10:20:00Z"
	w := mustWindow(t, start, end)
	r := logreader.New(mustParser(t), w)

	entries, err := r.Read(strings.NewReader(sampleLogs))
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	_ = time.Now() // ensure time import is used
	if len(entries) != 4 {
		t.Errorf("expected 4 entries (boundaries inclusive), got %d", len(entries))
	}
}
