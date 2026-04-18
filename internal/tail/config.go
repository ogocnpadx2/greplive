package tail

// Config holds options for the tail buffer.
type Config struct {
	// Lines is the number of trailing lines to retain.
	// Defaults to 10 when zero.
	Lines int
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{Lines: 10}
}

// Build constructs a Buffer from the Config.
func (c Config) Build() *Buffer {
	n := c.Lines
	if n <= 0 {
		n = DefaultConfig().Lines
	}
	return New(n)
}
