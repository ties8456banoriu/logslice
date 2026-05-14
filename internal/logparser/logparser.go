package logparser

import (
	"encoding/json"
	"fmt"
	"time"
)

// Entry represents a single parsed log entry with a timestamp and raw data.
type Entry struct {
	Timestamp time.Time
	Raw       string
	Fields    map[string]interface{}
}

// Parser defines the interface for log line parsers.
type Parser interface {
	Parse(line string) (*Entry, error)
}

// JSONParser parses JSON-structured log lines.
type JSONParser struct {
	TimeField  string
	TimeFormat string
}

// NewJSONParser creates a JSONParser with the given timestamp field name and format.
// If timeFormat is empty, time.RFC3339 is used.
func NewJSONParser(timeField, timeFormat string) *JSONParser {
	if timeFormat == "" {
		timeFormat = time.RFC3339
	}
	if timeField == "" {
		timeField = "time"
	}
	return &JSONParser{TimeField: timeField, TimeFormat: timeFormat}
}

// Parse parses a single JSON log line into an Entry.
func (p *JSONParser) Parse(line string) (*Entry, error) {
	if line == "" {
		return nil, fmt.Errorf("empty line")
	}

	var fields map[string]interface{}
	if err := json.Unmarshal([]byte(line), &fields); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	raw, ok := fields[p.TimeField]
	if !ok {
		return nil, fmt.Errorf("timestamp field %q not found", p.TimeField)
	}

	timeStr, ok := raw.(string)
	if !ok {
		return nil, fmt.Errorf("timestamp field %q is not a string", p.TimeField)
	}

	t, err := time.Parse(p.TimeFormat, timeStr)
	if err != nil {
		return nil, fmt.Errorf("cannot parse timestamp %q with format %q: %w", timeStr, p.TimeFormat, err)
	}

	return &Entry{
		Timestamp: t,
		Raw:       line,
		Fields:    fields,
	}, nil
}
