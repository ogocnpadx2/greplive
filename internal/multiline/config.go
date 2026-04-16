package multiline

// Config holds configuration for the multiline collector.
type Config struct {
	// StartPattern is a regex that identifies the first line of a new event.
	StartPattern string `json:"start_pattern"`
	// ContinuePattern is an optional regex for continuation lines.
	// When empty, indented lines (tab or space prefix) are treated as continuations.
	ContinuePattern string `json:"continue_pattern"`
	// MaxLines is the maximum number of lines to merge. 0 means unlimited.
	MaxLines int `json:"max_lines"`
}

// DefaultConfig returns a Config suited for typical Java/Go stack traces.
func DefaultConfig() Config {
	return Config{
		StartPattern:    `^(ERROR|WARN|INFO|DEBUG|FATAL)`,
		ContinuePattern: "",
		MaxLines:        50,
	}
}

// Build constructs a Collector from the Config.
func (cfg Config) Build() (*Collector, error) {
	return New(cfg.StartPattern, cfg.ContinuePattern, cfg.MaxLines)
}
