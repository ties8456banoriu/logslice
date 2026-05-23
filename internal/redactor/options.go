package redactor

// Option is a functional option for configuring a Redactor.
type Option func(*config)

type config struct {
	rules map[string]string
}

// WithRule adds a single pattern→replacement rule to the Redactor.
// Calling WithRule multiple times accumulates rules.
func WithRule(pattern, replacement string) Option {
	return func(c *config) {
		if c.rules == nil {
			c.rules = make(map[string]string)
		}
		c.rules[pattern] = replacement
	}
}

// WithDefaultRules merges the built-in default rules into the Redactor.
func WithDefaultRules() Option {
	return func(c *config) {
		if c.rules == nil {
			c.rules = make(map[string]string)
		}
		for k, v := range DefaultRules() {
			c.rules[k] = v
		}
	}
}

// NewFromOptions constructs a Redactor using functional options.
// It is an alternative to New for callers that prefer the options pattern.
func NewFromOptions(opts ...Option) (*Redactor, error) {
	cfg := &config{rules: make(map[string]string)}
	for _, o := range opts {
		o(cfg)
	}
	return New(cfg.rules)
}
