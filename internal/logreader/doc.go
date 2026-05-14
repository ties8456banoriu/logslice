// Package logreader implements line-by-line reading of structured log files
// with time-window filtering.
//
// A Reader wraps a [logparser.JSONParser] and a [timewindow.TimeWindow] to
// provide a simple API for extracting log entries that fall within a specific
// time range.
//
// Basic usage:
//
//	parser, _ := logparser.NewJSONParser(logparser.Options{})
//	win, _   := timewindow.New("2024-01-15T10:00:00Z", "2024-01-15T11:00:00Z")
//	r := logreader.New(parser, win)
//
//	entries, err := r.ReadFile("/var/log/app.log")
//
// Lines that cannot be parsed are silently skipped so that a single
// malformed line does not abort processing of an entire file.
package logreader
