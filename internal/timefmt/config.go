package timefmt

// Config holds user-supplied timestamp format configuration.
type Config struct {
	// Name is one of the known format names (e.g. "rfc3339", "short").
	// Defaults to "rfc3339" when empty.
	Name string
	// UTC forces output in UTC regardless of the log line's timezone.
	UTC bool
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{Name: "rfc3339", UTC: false}
}

// Build resolves the Config into a Formatter.
func (c Config) Build() (Formatter, error) {
	name := c.Name
	if name == "" {
		name = "rfc3339"
	}
	fmt, ok := Parse(name)
	if !ok {
		return Formatter{}, &UnknownFormatError{Name: name}
	}
	return Formatter{fmt: fmt, utc: c.UTC}, nil
}

// UnknownFormatError is returned when the format name is not recognised.
type UnknownFormatError struct {
	Name string
}

func (e *UnknownFormatError) Error() string {
	return "timefmt: unknown format " + e.Name
}
