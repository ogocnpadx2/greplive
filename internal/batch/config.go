package batch

import (
	"fmt"
	"time"
)

// Config holds user-facing parameters for constructing a Batcher.
type Config struct {
	// MaxSize is the maximum number of lines per batch. 0 disables size-based
	// flushing.
	MaxSize int
	// Interval is the maximum time to wait before flushing an incomplete batch.
	// Defaults to 1 second when zero.
	Interval time.Duration
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		MaxSize:  100,
		Interval: time.Second,
	}
}

// Build validates the Config and returns a ready-to-use Batcher.
func (c Config) Build() (*Batcher, error) {
	if c.MaxSize < 0 {
		return nil, fmt.Errorf("batch: MaxSize must be >= 0, got %d", c.MaxSize)
	}
	return New(c.MaxSize, c.Interval), nil
}
