// Package logparser provides parsers for structured log lines.
//
// Currently supported formats:
//
//   - JSON: parses newline-delimited JSON log entries, extracting a
//     configurable timestamp field. Suitable for logs produced by
//     structured loggers such as zerolog, zap, or logrus in JSON mode.
//
// Usage:
//
//	p := logparser.NewJSONParser("time", time.RFC3339)
//	entry, err := p.Parse(line)
//	if err != nil {
//	    // handle parse error (malformed JSON, missing/invalid timestamp)
//	}
//	// entry.Timestamp holds the parsed time.Time value
//	// entry.Fields holds all decoded JSON fields
//	// entry.Raw holds the original log line
package logparser
