package highlight

import "fmt"

// TermConfig holds configuration for a single highlight term.
type TermConfig struct {
	Pattern string
	Color   string
}

// Config holds a list of highlight term configurations.
type Config struct {
	Terms []TermConfig
}

// Build constructs a slice of Highlighter instances from the Config.
// Returns an error if any pattern is invalid.
func (c *Config) Build() ([]*Highlighter, error) {
	var highlighters []*Highlighter
	for _, term := range c.Terms {
		color := ColorByName(term.Color)
		h, err := New(term.Pattern, color)
		if err != nil {
			return nil, fmt.Errorf("highlight: invalid pattern %q: %w", term.Pattern, err)
		}
		highlighters = append(highlighters, h)
	}
	return highlighters, nil
}

// DefaultConfig returns a Config with common log term highlights pre-configured.
func DefaultConfig() *Config {
	return &Config{
		Terms: []TermConfig{
			{Pattern: `(?i)error`, Color: "red"},
			{Pattern: `(?i)warn(ing)?`, Color: "yellow"},
			{Pattern: `(?i)info`, Color: "cyan"},
		},
	}
}
