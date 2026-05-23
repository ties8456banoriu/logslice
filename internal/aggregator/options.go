package aggregator

import "time"

// Config holds configuration for an Aggregator.
type Config struct {
	Window time.Duration
}

func defaultConfig() Config {
	return Config{
		Window: time.Minute,
	}
}

// Option is a functional option for configuring an Aggregator.
type Option func(*Config)

// WithWindow sets the bucket window duration.
// Values <= 0 are ignored and the default (1 minute) is used.
func WithWindow(d time.Duration) Option {
	return func(c *Config) {
		if d > 0 {
			c.Window = d
		}
	}
}

// NewFromOptions constructs an Aggregator using functional options.
func NewFromOptions(opts ...Option) *Aggregator {
	cfg := defaultConfig()
	for _, o := range opts {
		o(&cfg)
	}
	return New(cfg.Window)
}
