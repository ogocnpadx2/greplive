package cli

import (
	"errors"
	"flag"
	"fmt"
	"time"
)

// Flags holds all parsed CLI options.
type Flags struct {
	Pattern       string
	Level         string
	Follow        bool
	Timestamp     bool
	StatsInterval time.Duration
	JSON          bool
	MaxRate       int
	Dedupe        bool
	MaxLen        int
	Before        int
	After         int
}

// ParseFlags parses os.Args using the provided FlagSet and returns Flags.
func ParseFlags(fs *flag.FlagSet, args []string) (Flags, error) {
	var f Flags
	var statsStr string
	var levelStr string

	fs.StringVar(&f.Pattern, "pattern", "", "regex filter pattern")
	fs.StringVar(&levelStr, "level", "any", "minimum severity level (debug|info|warn|error|fatal|any)")
	fs.BoolVar(&f.Follow, "follow", false, "tail file and follow new lines")
	fs.BoolVar(&f.Timestamp, "timestamp", false, "prefix output with timestamp")
	fs.StringVar(&statsStr, "stats", "0s", "stats reporting interval (0 to disable)")
	fs.BoolVar(&f.JSON, "json", false, "output lines as JSON")
	fs.IntVar(&f.MaxRate, "max-rate", 0, "max lines per second (0 = unlimited)")
	fs.BoolVar(&f.Dedupe, "dedupe", false, "suppress duplicate consecutive lines")
	fs.IntVar(&f.MaxLen, "max-len", 0, "truncate lines to N runes (0 = off)")
	fs.IntVar(&f.Before, "before", 0, "context lines before match")
	fs.IntVar(&f.After, "after", 0, "context lines after match")

	if err := fs.Parse(args); err != nil {
		return Flags{}, err
	}

	validLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true, "fatal": true, "any": true}
	if !validLevels[levelStr] {
		return Flags{}, fmt.Errorf("invalid level %q", levelStr)
	}
	f.Level = levelStr

	d, err := time.ParseDuration(statsStr)
	if err != nil {
		return Flags{}, errors.New("invalid stats interval: " + err.Error())
	}
	f.StatsInterval = d

	return f, nil
}
