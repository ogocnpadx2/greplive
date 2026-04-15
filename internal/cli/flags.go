package cli

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/user/greplive/internal/severity"
)

// Config holds all values parsed from command-line flags.
type Config struct {
	Pattern       string
	MinLevel      severity.Level
	File          string
	ShowTimestamp bool
	ShowLevel     bool
	ShowStats     bool
	StatsInterval time.Duration
}

// ParseFlags parses args and returns a populated Config or an error.
func ParseFlags(args []string) (*Config, error) {
	fs := flag.NewFlagSet("greplive", flag.ContinueOnError)

	pattern := fs.String("pattern", "", "regex pattern to filter lines (empty = match all)")
	levelStr := fs.String("level", "debug", "minimum severity level: debug|info|warn|error")
	file := fs.String("file", "", "file to tail (omit to read from stdin)")
	timestamp := fs.Bool("timestamp", false, "prefix each line with a timestamp")
	showLevel := fs.Bool("level-prefix", false, "prefix each line with its severity level")
	showStats := fs.Bool("stats", false, "print running statistics to stderr")
	statsInterval := fs.Duration("stats-interval", 5*time.Second, "how often to print stats")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	lvl, err := severity.ParseLevel(strings.ToLower(*levelStr))
	if err != nil {
		return nil, fmt.Errorf("unknown level %q: %w", *levelStr, err)
	}

	return &Config{
		Pattern:       *pattern,
		MinLevel:      lvl,
		File:          *file,
		ShowTimestamp: *timestamp,
		ShowLevel:     *showLevel,
		ShowStats:     *showStats,
		StatsInterval: *statsInterval,
	}, nil
}
