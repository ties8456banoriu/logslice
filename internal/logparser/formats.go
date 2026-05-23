package logparser

import "fmt"

// Format represents a supported log format.
type Format int

const (
	// FormatJSON represents newline-delimited JSON logs.
	FormatJSON Format = iota
)

// String returns the string representation of the Format.
func (f Format) String() string {
	switch f {
	case FormatJSON:
		return "json"
	default:
		return fmt.Sprintf("Format(%d)", int(f))
	}
}

// FormatError is returned when an unsupported format string is provided.
type FormatError struct {
	Value string
}

func (e *FormatError) Error() string {
	return fmt.Sprintf("logparser: unsupported format %q; supported formats: json", e.Value)
}

// ParseFormat converts a string to a Format constant.
// It returns a FormatError if the format is not recognised.
func ParseFormat(s string) (Format, error) {
	switch s {
	case "json", "JSON":
		return FormatJSON, nil
	default:
		return 0, &FormatError{Value: s}
	}
}

// NewParser returns a Parser for the given Format, applying any Options.
// It returns a FormatError if the format is not supported.
func NewParser(f Format, opts ...Option) (Parser, error) {
	switch f {
	case FormatJSON:
		return NewJSONParser(opts...), nil
	default:
		return nil, &FormatError{Value: f.String()}
	}
}
