package redactor_test

import (
	"strings"
	"testing"

	"logslice/internal/redactor"
)

func TestNew_ValidRules(t *testing.T) {
	r, err := redactor.New(map[string]string{`\d+`: "[NUM]"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r == nil {
		t.Fatal("expected non-nil Redactor")
	}
}

func TestNew_InvalidPattern(t *testing.T) {
	_, err := redactor.New(map[string]string{`[invalid`: "x"})
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
	var pe *redactor.PatternError
	if !isPatternError(err, &pe) {
		t.Fatalf("expected PatternError, got %T", err)
	}
	if !strings.Contains(pe.Pattern, "invalid") {
		t.Errorf("PatternError.Pattern = %q, want substring 'invalid'", pe.Pattern)
	}
}

func TestRedact_AppliesRules(t *testing.T) {
	r, _ := redactor.New(map[string]string{`\d+`: "[NUM]"})
	got := r.Redact("error code 404 at line 12")
	want := "error code [NUM] at line [NUM]"
	if got != want {
		t.Errorf("Redact() = %q, want %q", got, want)
	}
}

func TestRedact_NoMatch(t *testing.T) {
	r, _ := redactor.New(map[string]string{`\d+`: "[NUM]"})
	input := "no digits here"
	if got := r.Redact(input); got != input {
		t.Errorf("Redact() = %q, want %q", got, input)
	}
}

func TestRedactMap_StringValues(t *testing.T) {
	r, _ := redactor.New(map[string]string{`secret`: "[REDACTED]"})
	fields := map[string]any{
		"msg":   "my secret value",
		"count": 42,
		"flag":  true,
	}
	out := r.RedactMap(fields)
	if out["msg"] != "my [REDACTED] value" {
		t.Errorf("msg = %q, want 'my [REDACTED] value'", out["msg"])
	}
	if out["count"] != 42 {
		t.Errorf("count = %v, want 42", out["count"])
	}
	if out["flag"] != true {
		t.Errorf("flag = %v, want true", out["flag"])
	}
}

func TestRedactMap_DoesNotMutateInput(t *testing.T) {
	r, _ := redactor.New(map[string]string{`secret`: "[REDACTED]"})
	original := map[string]any{"msg": "my secret"}
	r.RedactMap(original)
	if original["msg"] != "my secret" {
		t.Error("RedactMap mutated the original map")
	}
}

func TestDefaultRules_NotEmpty(t *testing.T) {
	rules := redactor.DefaultRules()
	if len(rules) == 0 {
		t.Error("DefaultRules returned empty map")
	}
	_, err := redactor.New(rules)
	if err != nil {
		t.Fatalf("DefaultRules produced invalid patterns: %v", err)
	}
}

// isPatternError is a helper to avoid importing errors in the test.
func isPatternError(err error, target **redactor.PatternError) bool {
	if pe, ok := err.(*redactor.PatternError); ok {
		*target = pe
		return true
	}
	return false
}
