package scope

// Config holds the configuration for a Scope filter.
type Config struct {
	// Start is the regex pattern that begins the active window.
	Start string
	// End is the regex pattern that closes the active window.
	End string
}

// DefaultConfig returns a Config with both patterns empty (disabled).
func DefaultConfig() Config {
	return Config{}
}

// Build constructs a Scope from the Config.
func (c Config) Build() (*Scope, error) {
	return New(c.Start, c.End)
}
