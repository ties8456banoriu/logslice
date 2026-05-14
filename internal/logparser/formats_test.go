package logparser

import (
	"testing"
)

func TestNewParser_JSONFormat(t *testing.T) {
	cases := []struct {
		name   string
		format string
	}{
		{"explicit json", "json"},
		{"uppercase JSON", "JSON"},
		{"empty defaults to json", ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p, err := NewParser(ParserConfig{Format: tc.format})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if p == nil {
				t.Fatal("expected non-nil parser")
			}
		})
	}
}

func TestNewParser_UnsupportedFormat(t *testing.T) {
	_, err := NewParser(ParserConfig{Format: "csv"})
	if err == nil {
		t.Fatal("expected error for unsupported format, got nil")
	}
	fe, ok := err.(*FormatError)
	if !ok {
		t.Fatalf("expected *FormatError, got %T", err)
	}
	if fe.Format != "csv" {
		t.Errorf("expected Format=\"csv\", got %q", fe.Format)
	}
	if fe.Error() == "" {
		t.Error("expected non-empty error message")
	}
}

func TestFormatError_Message(t *testing.T) {
	fe := &FormatError{Format: "xml"}
	msg := fe.Error()
	if msg == "" {
		t.Fatal("expected non-empty error message")
	}
	for _, supported := range SupportedFormats {
		if !contains(msg, supported) {
			t.Errorf("expected error message to mention %q, got: %s", supported, msg)
		}
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsStr(s, sub))
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
