package output_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"logslice/internal/output"
)

var sampleEntry = output.Entry{
	Timestamp: time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
	Message:   "user logged in",
	Fields:    map[string]any{"user_id": "42", "level": "info"},
}

func TestNew_ValidFormats(t *testing.T) {
	for _, f := range []output.Format{output.FormatJSON, output.FormatText, output.FormatCompact} {
		_, err := output.New(&bytes.Buffer{}, f)
		if err != nil {
			t.Errorf("New(%q) unexpected error: %v", f, err)
		}
	}
}

func TestNew_NilWriter(t *testing.T) {
	_, err := output.New(nil, output.FormatJSON)
	if err == nil {
		t.Fatal("expected error for nil writer")
	}
}

func TestNew_InvalidFormat(t *testing.T) {
	_, err := output.New(&bytes.Buffer{}, output.Format("xml"))
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestWrite_JSON(t *testing.T) {
	var buf bytes.Buffer
	w, _ := output.New(&buf, output.FormatJSON)
	if err := w.Write(sampleEntry); err != nil {
		t.Fatalf("Write error: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if got["message"] != "user logged in" {
		t.Errorf("unexpected message: %v", got["message"])
	}
	if got["user_id"] != "42" {
		t.Errorf("expected user_id field to be propagated")
	}
}

func TestWrite_Compact(t *testing.T) {
	var buf bytes.Buffer
	w, _ := output.New(&buf, output.FormatCompact)
	if err := w.Write(sampleEntry); err != nil {
		t.Fatalf("Write error: %v", err)
	}
	line := strings.TrimSpace(buf.String())
	if strings.Contains(line, "\n  ") {
		t.Error("compact format should not contain indentation")
	}
}

func TestWrite_Text(t *testing.T) {
	var buf bytes.Buffer
	w, _ := output.New(&buf, output.FormatText)
	if err := w.Write(sampleEntry); err != nil {
		t.Fatalf("Write error: %v", err)
	}
	if !strings.Contains(buf.String(), "user logged in") {
		t.Error("text output should contain the message")
	}
	if !strings.Contains(buf.String(), "2024-01-15") {
		t.Error("text output should contain the date")
	}
}
