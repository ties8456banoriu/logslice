package logparser_test

import (
	"strings"
	"testing"

	"logslice/internal/logparser"
)

func TestNewParser_JSONFormat(t *testing.T) {
	p, err := logparser.NewParser(logparser.FormatJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil parser")
	}
}

func TestNewParser_UnsupportedFormat(t *testing.T) {
	_, err := logparser.NewParser(logparser.Format(99))
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
	if !containsStr(err.Error(), "unsupported format") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestFormatError_Message(t *testing.T) {
	_, err := logparser.ParseFormat("xml")
	if err == nil {
		t.Fatal("expected error")
	}
	if !containsStr(err.Error(), "xml") {
		t.Errorf("error should mention the bad value, got: %v", err)
	}
	if !containsStr(err.Error(), "json") {
		t.Errorf("error should list supported formats, got: %v", err)
	}
}

func contains(t *testing.T, haystack, needle string) bool {
	t.Helper()
	return strings.Contains(haystack, needle)
}

func containsStr(haystack, needle string) bool {
	return strings.Contains(haystack, needle)
}

func TestParseFormat_ValidInputs(t *testing.T) {
	cases := []struct {
		input string
		want  logparser.Format
	}{
		{"json", logparser.FormatJSON},
		{"JSON", logparser.FormatJSON},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got, err := logparser.ParseFormat(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestParseFormat_InvalidInput(t *testing.T) {
	for _, s := range []string{"", "csv", "text", "xml"} {
		t.Run(s, func(t *testing.T) {
			_, err := logparser.ParseFormat(s)
			if err == nil {
				t.Fatalf("expected error for %q", s)
			}
		})
	}
}
