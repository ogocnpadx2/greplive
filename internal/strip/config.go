package strip

// Config holds the flags used to construct a Stripper.
type Config struct {
	// ANSI controls whether ANSI escape sequences are removed.
	ANSI bool
	// Control controls whether raw control characters are removed.
	Control bool
}

// DefaultConfig returns a Config with ANSI stripping enabled and control
// character stripping disabled — the most common production default.
func DefaultConfig() Config {
	return Config{
		ANSI:    true,
		Control: false,
	}
}

// Build constructs a Stripper from the Config.
func (c Config) Build() *Stripper {
	return New(c.ANSI, c.Control)
}
