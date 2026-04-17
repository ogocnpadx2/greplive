package jsonparse

// Config holds declarative settings for building a Formatter.
type Config struct {
	// MessageKey is the JSON field treated as the primary log message.
	// Defaults to "message" when empty.
	MessageKey string

	// Pretty enables indented JSON output.
	Pretty bool
}

// Build constructs a Formatter from the Config.
func (c Config) Build() *Formatter {
	key := c.MessageKey
	if key == "" {
		key = "message"
	}
	return New(
		WithMessageKey(key),
		WithPretty(c.Pretty),
	)
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		MessageKey: "message",
		Pretty:     false,
	}
}
