package pipeline

import (
	"fmt"
	"strings"
)

// Validate checks that the Config fields are consistent and returns a
// descriptive error if any required field is missing or obviously wrong.
func (c *Config) Validate() error {
	var errs []string

	if len(c.Patterns) == 0 {
		errs = append(errs, "at least one file pattern is required")
	}
	if c.From == "" {
		errs = append(errs, "'from' timestamp is required")
	}
	if c.To == "" {
		errs = append(errs, "'to' timestamp is required")
	}
	if c.Writer == nil {
		errs = append(errs, "writer must not be nil")
	}
	if c.Format == "" {
		c.Format = "json" // apply default silently
	}

	if len(errs) > 0 {
		return fmt.Errorf("pipeline config: %s", strings.Join(errs, "; "))
	}
	return nil
}
