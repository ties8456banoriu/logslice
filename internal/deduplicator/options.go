package deduplicator

// Option configures a Deduplicator.
type Option func(*config)

type config struct {
	capacity int
}

func defaultConfig() config {
	return config{
		capacity: 1024,
	}
}

// WithCapacity sets the maximum number of unique hashes to track.
// Values <= 0 are ignored and the default (1024) is retained.
func WithCapacity(n int) Option {
	return func(c *config) {
		if n > 0 {
			c.capacity = n
		}
	}
}

// NewFromOptions creates a Deduplicator configured by the provided options.
func NewFromOptions(opts ...Option) *Deduplicator {
	cfg := defaultConfig()
	for _, o := range opts {
		o(&cfg)
	}
	return New(cfg.capacity)
}
