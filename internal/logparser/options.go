package logparser

// options holds configuration for parser construction.
type options struct {
	timeField  string
	timeLayout string
}

// defaultOptions returns sensible defaults for parser options.
func defaultOptions() options {
	return options{
		timeField:  "time",
		timeLayout: "2006-01-02T15:04:05Z07:00",
	}
}

// Option is a functional option for configuring a Parser.
type Option func(*options)

// WithTimeField sets the JSON key used to extract the log timestamp.
func WithTimeField(field string) Option {
	return func(o *options) {
		if field != "" {
			o.timeField = field
		}
	}
}

// WithTimeLayout sets the time.Parse layout used when parsing timestamps.
func WithTimeLayout(layout string) Option {
	return func(o *options) {
		if layout != "" {
			o.timeLayout = layout
		}
	}
}
