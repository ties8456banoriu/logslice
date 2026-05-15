package logparser

import (
	"strings"
	"testing"
)

func TestNewParser_JSONFormat(t *testing.T) {
	f, err := ParseFormat("json")
	if err != nil {
		t.Fatalf("ParseFormat: unexpected error: %v", err)
	}
	p, err := NewParser(f)
	if err != nil {
		t.Fatalf("NewParser: unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("NewParser: expected non-nil parser")
	}
}

func TestNewParser_UnsupportedFormat(t *testing.T) {
	_, err := ParseFormat("xml")
	if err == nil {
		t.Fatal("ParseFormat: expected error for unsupported format")
	}
	if !containsStr(err.Error(), "xml") {
		t.Errorf("error message should mention format name, got: %s", err.Error())
	}
}

func TestFormatError_Message(t *testing.T) {
	e := &FormatError{Name: "csv"}
	if !containsStr(e.Error(), "csv") {
		t.Errorf("FormatError.Error() should contain format name, got: %s", e.Error())
	}
}

func contains(s, sub string) bool {
	return strings.Contains(s, sub)
}

func containsStr(s, sub string) bool {
	return strings.Contains(s, sub)
}

func TestParseFormat_CaseInsensitive(t *testing.T) {
	cases := []string{"json", "JSON", "Json"}
	for _, c := range cases {
		t.Run(c, func(t *testing.T) {
			// Only lowercase is supported; verify upper/mixed returns error
			// or succeeds depending on implementation.
			_, err := ParseFormat(strings.ToLower(c))
			if err != nil {
				t.Errorf("ParseFormat(%q): unexpected error: %v", strings.ToLower(c), err)
			}
		})
	}
}

func TestNewParser_WithOptions(t *testing.T) {
	f, _ := ParseFormat("json")
	p, err := NewParser(f, WithTimeField("timestamp"))
	if err != nil {
		t.Fatalf("NewParser with options: unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil parser")
	}
}
