package logparser

import "fmt"

// Format represents a supported log format.
type Format int

const (
	// FormatJSON represents newline-delimited JSON log format.
	FormatJSON Format = iota
)

// FormatError is returned when an unsupported format is requested.
type FormatError struct {
	Name string
}

func (e *FormatError) Error() string {
	return fmt.Sprintf("logparser: unsupported format %q", e.Name)
}

// formatNames maps string names to Format constants.
var formatNames = map[string]Format{
	"json": FormatJSON,
}

// ParseFormat converts a string name to a Format constant.
// It returns a FormatError if the name is not recognised.
func ParseFormat(name string) (Format, error) {
	f, ok := formatNames[name]
	if !ok {
		return 0, &FormatError{Name: name}
	}
	return f, nil
}

// NewParser constructs a Parser for the given Format.
// It returns a FormatError for unrecognised formats.
func NewParser(f Format, opts ...Option) (Parser, error) {
	switch f {
	case FormatJSON:
		return NewJSONParser(opts...), nil
	default:
		return nil, &FormatError{Name: fmt.Sprintf("%d", f)}
	}
}
