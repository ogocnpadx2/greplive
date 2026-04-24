package reopen

// Config holds options for a managed reopen.Reader.
type Config struct {
	// Path is the file to open and monitor.
	Path string

	// AutoWatch, when true, starts WatchAndReopen automatically
	// when Build is called. The caller must supply a rotation
	// channel via RotateCh.
	AutoWatch bool
}

// DefaultConfig returns a Config with safe defaults.
func DefaultConfig() Config {
	return Config{
		AutoWatch: false,
	}
}

// Build validates the Config and returns a new Reader.
// It returns an error if Path is empty or the file cannot be opened.
func (c Config) Build() (*Reader, error) {
	if c.Path == "" {
		return nil, errEmptyPath
	}
	return New(c.Path)
}

import "errors"

var errEmptyPath = errors.New("reopen: path must not be empty")
