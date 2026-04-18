package alert

import (
	"fmt"
	"io"
	"time"
)

// Config holds the declarative configuration for an Alert.
type Config struct {
	// Pattern is the regex to match against each line.
	Pattern string
	// Threshold is the number of matches within Window that triggers an alert.
	Threshold int
	// Window is the rolling time window for counting matches.
	Window time.Duration
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Threshold: 10,
		Window:    time.Minute,
	}
}

// Build constructs an Alert from the Config, writing alerts to w.
func (c Config) Build(w io.Writer) (*Alert, error) {
	if c.Threshold < 0 {
		return nil, fmt.Errorf("alert: threshold must be non-negative")
	}
	window := c.Window
	if window <= 0 {
		window = time.Minute
	}
	return New(c.Pattern, c.Threshold, window, w)
}
