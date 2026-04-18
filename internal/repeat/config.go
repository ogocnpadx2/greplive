package repeat

import (
	"fmt"
	"time"
)

// Config holds the user-facing configuration for the repeat suppressor.
type Config struct {
	// Max is the maximum number of identical consecutive lines allowed within
	// Window before further occurrences are suppressed. 0 disables the feature.
	Max int
	// Window is the duration over which repetitions are counted.
	// Defaults to 1 minute when zero.
	Window time.Duration
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Max:    0,
		Window: time.Minute,
	}
}

// Build validates the Config and returns a ready-to-use Repeater.
func (c Config) Build() (*Repeater, error) {
	if c.Max < 0 {
		return nil, fmt.Errorf("repeat: max must be >= 0, got %d", c.Max)
	}
	win := c.Window
	if win <= 0 {
		win = time.Minute
	}
	return New(c.Max, win), nil
}
