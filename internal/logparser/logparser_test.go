package logparser

import (
	"testing"
	"time"
)

func TestNewJSONParser_Defaults(t *testing.T) {
	p := NewJSONParser("", "")
	if p.TimeField != "time" {
		t.Errorf("expected default TimeField \"time\", got %q", p.TimeField)
	}
	if p.TimeFormat != time.RFC3339 {
		t.Errorf("expected default TimeFormat RFC3339, got %q", p.TimeFormat)
	}
}

func TestJSONParser_Parse_Valid(t *testing.T) {
	p := NewJSONParser("time", time.RFC3339)
	line := `{"time":"2024-01-15T10:00:00Z","level":"info","msg":"hello"}`

	entry, err := p.Parse(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.Raw != line {
		t.Errorf("expected Raw to equal input line")
	}
	expected := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	if !entry.Timestamp.Equal(expected) {
		t.Errorf("expected timestamp %v, got %v", expected, entry.Timestamp)
	}
	if entry.Fields["msg"] != "hello" {
		t.Errorf("expected msg field \"hello\", got %v", entry.Fields["msg"])
	}
}

func TestJSONParser_Parse_EmptyLine(t *testing.T) {
	p := NewJSONParser("", "")
	_, err := p.Parse("")
	if err == nil {
		t.Error("expected error for empty line")
	}
}

func TestJSONParser_Parse_InvalidJSON(t *testing.T) {
	p := NewJSONParser("", "")
	_, err := p.Parse("not json at all")
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestJSONParser_Parse_MissingTimeField(t *testing.T) {
	p := NewJSONParser("timestamp", "")
	_, err := p.Parse(`{"level":"info","msg":"no ts"}`)
	if err == nil {
		t.Error("expected error for missing timestamp field")
	}
}

func TestJSONParser_Parse_InvalidTimeFormat(t *testing.T) {
	p := NewJSONParser("time", time.RFC3339)
	_, err := p.Parse(`{"time":"15/01/2024 10:00:00","msg":"bad format"}`)
	if err == nil {
		t.Error("expected error for unparseable timestamp")
	}
}

func TestJSONParser_Parse_CustomTimeField(t *testing.T) {
	p := NewJSONParser("@timestamp", time.RFC3339)
	line := `{"@timestamp":"2024-06-01T08:30:00Z","service":"api"}`
	entry, err := p.Parse(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := time.Date(2024, 6, 1, 8, 30, 0, 0, time.UTC)
	if !entry.Timestamp.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, entry.Timestamp)
	}
}
