// Package output provides serialization of structured log entries to an
// io.Writer in multiple formats.
//
// Supported formats:
//
//   - FormatJSON    – pretty-printed JSON (one object per line).
//   - FormatCompact – single-line JSON (one object per line, no indentation).
//   - FormatText    – human-readable text using a Go text/template.
//
// Usage:
//
//	w, err := output.New(os.Stdout, output.FormatJSON)
//	if err != nil { ... }
//	w.Write(output.Entry{
//		Timestamp: time.Now(),
//		Message:   "hello",
//		Fields:    map[string]any{"level": "info"},
//	})
package output
