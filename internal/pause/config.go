package pause

// Config holds options for constructing a Pauser.
type Config struct {
	// StartPaused controls whether the Pauser begins in the paused state.
	StartPaused bool
}

// Build constructs a Pauser from the Config.
func (c Config) Build() *Pauser {
	p := New()
	if c.StartPaused {
		p.Pause()
	}
	return p
}

// DefaultConfig returns a Config with sensible defaults (not paused).
func DefaultConfig() Config {
	return Config{StartPaused: false}
}
