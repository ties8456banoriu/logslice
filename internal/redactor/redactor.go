package redactor

import (
	"regexp"
	"strings"
)

// Rule describes a single redaction rule: a compiled pattern and its replacement.
type Rule struct {
	Pattern     *regexp.Regexp
	Replacement string
}

// Redactor applies a set of redaction rules to log line values.
type Redactor struct {
	rules []Rule
}

// New creates a Redactor from the provided rules.
// Returns an error if any pattern fails to compile.
func New(rawRules map[string]string) (*Redactor, error) {
	rules := make([]Rule, 0, len(rawRules))
	for pattern, replacement := range rawRules {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, &PatternError{Pattern: pattern, Err: err}
		}
		rules = append(rules, Rule{Pattern: re, Replacement: replacement})
	}
	return &Redactor{rules: rules}, nil
}

// Redact applies all rules to the given string and returns the sanitised result.
func (r *Redactor) Redact(value string) string {
	for _, rule := range r.rules {
		value = rule.Pattern.ReplaceAllString(value, rule.Replacement)
	}
	return value
}

// RedactMap returns a shallow copy of fields with values redacted.
// Only string values are processed; other types are passed through unchanged.
func (r *Redactor) RedactMap(fields map[string]any) map[string]any {
	out := make(map[string]any, len(fields))
	for k, v := range fields {
		if s, ok := v.(string); ok {
			out[k] = r.Redact(s)
		} else {
			out[k] = v
		}
	}
	return out
}

// PatternError is returned when a redaction pattern cannot be compiled.
type PatternError struct {
	Pattern string
	Err     error
}

func (e *PatternError) Error() string {
	return "redactor: invalid pattern " + e.Pattern + ": " + e.Err.Error()
}

// Unwrap satisfies the errors.Unwrap interface.
func (e *PatternError) Unwrap() error { return e.Err }

// DefaultRules returns a set of commonly useful redaction rules.
func DefaultRules() map[string]string {
	return map[string]string{
		`(?i)(password|passwd|secret|token)=[^\s&]+`: strings.Repeat("*", 8),
		`\b\d{4}[- ]?\d{4}[- ]?\d{4}[- ]?\d{4}\b`:  "[CARD]",
	}
}
