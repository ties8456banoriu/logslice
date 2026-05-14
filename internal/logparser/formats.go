package logparser

import (
	"fmt"
	"strings"
)

// SupportedFormats lists the log formats this package can parse.
var SupportedFormats = []string{"json"}

// FormatError is returned when an unsupported log format is requested.
type FormatError struct {
	Format string
}

func (e *FormatError) Error() string {
	return fmt.Sprintf("unsupported log format: %q (supported: %s)",
		e.Format, strings.Join(SupportedFormats, ", "))
}

// ParserConfig holds configuration for creating a log parser.
type ParserConfig struct {
	// Format is the log format to parse (e.g. "json").
	Format string
	// TimeField is the name of the field containing the timestamp.
	TimeField string
	// TimeLayout is the Go time layout string used to parse timestamps.
	TimeLayout string
}

// NewParser creates a Parser for the given configuration.
// Returns a FormatError if the format is not supported.
func NewParser(cfg ParserConfig) (Parser, error) {
	switch strings.ToLower(cfg.Format) {
	case "json", "":
		return NewJSONParser(cfg.TimeField, cfg.TimeLayout)
	default:
		return nil, &FormatError{Format: cfg.Format}
	}
}
