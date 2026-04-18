package grep

// Config holds the user-facing configuration for the grep matcher.
type Config struct {
	Patterns []string
	Any      bool // OR logic when true, AND logic when false
}

// DefaultConfig returns a Config with AND logic and no patterns.
func DefaultConfig() Config {
	return Config{Any: false}
}

// Build constructs a Matcher from the config.
func (c Config) Build() (*Matcher, error) {
	return New(c.Patterns, c.Any)
}
