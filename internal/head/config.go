package head

const defaultLines = 10

// Config holds configuration for the head limiter.
type Config struct {
	Lines int // maximum number of lines to emit; 0 disables limiting
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{Lines: 0} // disabled by default
}

// Build constructs a Limiter from the Config.
// If Lines is zero or negative the returned limiter is disabled.
func (c Config) Build() (*Limiter, error) {
	lines := c.Lines
	if lines < 0 {
		lines = 0
	}
	return New(lines), nil
}

// BuildWithDefault constructs a Limiter, substituting defaultLines when
// Lines is zero.
func (c Config) BuildWithDefault() (*Limiter, error) {
	lines := c.Lines
	if lines <= 0 {
		lines = defaultLines
	}
	return New(lines), nil
}
