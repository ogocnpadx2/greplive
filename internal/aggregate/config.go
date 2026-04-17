package aggregate

import (
	"fmt"
	"regexp"
)

// Config holds options for building an Aggregator.
type Config struct {
	// Pattern is the regex that identifies lines to aggregate.
	Pattern string
}

// Build validates the config and returns a ready Aggregator.
// An empty Pattern returns an Aggregator that never aggregates.
func (c Config) Build() (*Aggregator, error) {
	if c.Pattern == "" {
		return New(nil), nil
	}
	re, err := regexp.Compile(c.Pattern)
	if err != nil {
		return nil, fmt.Errorf("aggregate: invalid pattern %q: %w", c.Pattern, err)
	}
	return New(re), nil
}
