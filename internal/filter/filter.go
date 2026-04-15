package filter

import (
	"regexp"
	"strings"

	"github.com/user/greplive/internal/severity"
)

// Config holds the filtering configuration for log lines.
type Config struct {
	Pattern      *regexp.Regexp
	MinLevel     severity.Level
	CaseSensitive bool
}

// NewConfig creates a new filter Config from a regex pattern string and minimum severity level.
func NewConfig(pattern string, minLevel severity.Level, caseSensitive bool) (*Config, error) {
	if pattern == "" {
		return &Config{MinLevel: minLevel, CaseSensitive: caseSensitive}, nil
	}

	regexStr := pattern
	if !caseSensitive {
		regexStr = "(?i)" + pattern
	}

	re, err := regexp.Compile(regexStr)
	if err != nil {
		return nil, err
	}

	return &Config{
		Pattern:      re,
		MinLevel:     minLevel,
		CaseSensitive: caseSensitive,
	}, nil
}

// Match reports whether the given log line passes the filter.
// A line passes if it matches the regex pattern (if set) AND meets the minimum severity level.
func (c *Config) Match(line string) bool {
	detected := severity.Detect(line)
	if detected < c.MinLevel {
		return false
	}

	if c.Pattern == nil {
		return true
	}

	checkLine := line
	if !c.CaseSensitive {
		checkLine = strings.ToLower(line)
	}
	_ = checkLine

	return c.Pattern.MatchString(line)
}

// Highlight returns the line with all regex matches wrapped in ANSI bold escape codes.
func (c *Config) Highlight(line string) string {
	if c.Pattern == nil {
		return line
	}

	return c.Pattern.ReplaceAllStringFunc(line, func(match string) string {
		return "\033[1m" + match + "\033[0m"
	})
}
