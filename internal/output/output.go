package output

import (
	"encoding/json"
	"fmt"
	"io"
	"text/template"
	"time"
)

// Format defines the output serialization format.
type Format string

const (
	FormatJSON    Format = "json"
	FormatText    Format = "text"
	FormatCompact Format = "compact"
)

// Entry represents a single log record to be written.
type Entry struct {
	Timestamp time.Time
	Message   string
	Fields    map[string]any
}

// Writer serializes log entries to an io.Writer in the requested format.
type Writer struct {
	format Format
	w      io.Writer
	tmpl   *template.Template
}

const defaultTextTemplate = `{{.Timestamp.Format "2006-01-02T15:04:05Z07:00"}} {{.Message}}\n`

// New creates a Writer for the given format and destination.
func New(w io.Writer, format Format) (*Writer, error) {
	if w == nil {
		return nil, fmt.Errorf("output: writer must not be nil")
	}
	switch format {
	case FormatJSON, FormatText, FormatCompact:
		// valid
	default:
		return nil, fmt.Errorf("output: unsupported format %q", format)
	}
	var tmpl *template.Template
	if format == FormatText {
		var err error
		tmpl, err = template.New("entry").Parse(defaultTextTemplate)
		if err != nil {
			return nil, fmt.Errorf("output: failed to parse template: %w", err)
		}
	}
	return &Writer{format: format, w: w, tmpl: tmpl}, nil
}

// Write serializes a single Entry to the underlying writer.
func (wr *Writer) Write(e Entry) error {
	switch wr.format {
	case FormatJSON:
		return wr.writeJSON(e, true)
	case FormatCompact:
		return wr.writeJSON(e, false)
	case FormatText:
		return wr.tmpl.Execute(wr.w, e)
	}
	return nil
}

func (wr *Writer) writeJSON(e Entry, indent bool) error {
	payload := map[string]any{
		"timestamp": e.Timestamp.Format(time.RFC3339Nano),
		"message":   e.Message,
	}
	for k, v := range e.Fields {
		if _, exists := payload[k]; !exists {
			payload[k] = v
		}
	}
	var b []byte
	var err error
	if indent {
		b, err = json.MarshalIndent(payload, "", "  ")
	} else {
		b, err = json.Marshal(payload)
	}
	if err != nil {
		return fmt.Errorf("output: marshal error: %w", err)
	}
	_, err = fmt.Fprintf(wr.w, "%s\n", b)
	return err
}
