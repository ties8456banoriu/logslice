package logparser

import "time"

// options holds configuration for a Parser.
type options struct {
	timeField  string
	timeLayout string
}

// Option is a functional option for configuring a Parser.
type Option func(*options)

// defaultOptions returns options pre-filled with sensible defaults.
func defaultOptions() options {
	return options{
		timeField:  "time",
		timeLayout: time.RFC3339,
	}
}

// WithTimeField overrides the JSON key used to extract the log timestamp.
// An empty value is silently ignored.
func WithTimeField(field string) Option {
	return func(o *options) {
		if field != "" {
			o.timeField = field
		}
	}
}

// WithTimeLayout overrides the time.Parse layout used to parse the timestamp.
// An empty value is silently ignored.
func WithTimeLayout(layout string) Option {
	return func(o *options) {
		if layout != "" {
			o.timeLayout = layout
		}
	}
}
