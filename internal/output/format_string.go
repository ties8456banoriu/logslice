package output

import "fmt"

// ParseFormat converts a string to a Format constant.
// It returns an error if the value is not recognised.
func ParseFormat(s string) (Format, error) {
	switch Format(s) {
	case FormatJSON:
		return FormatJSON, nil
	case FormatText:
		return FormatText, nil
	case FormatCompact:
		return FormatCompact, nil
	}
	return "", fmt.Errorf("output: unknown format %q; valid values: json, text, compact", s)
}

// String implements fmt.Stringer.
func (f Format) String() string { return string(f) }

// Formats returns all supported Format values.
func Formats() []Format {
	return []Format{FormatJSON, FormatText, FormatCompact}
}
