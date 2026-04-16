package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/user/greplive/internal/severity"
)

// Flags holds all parsed command-line options.
type Flags struct {
	Pattern       string
	Level         severity.Level
	Follow        bool
	File          string
	Timestamp     bool
	Rate          int
	Dedupe        bool
	MaxLineLen    int
	CheckpointFile string
}

// ParseFlags parses os.Args and returns a populated Flags or an error.
func ParseFlags(args []string) (*Flags, error) {
	fs := flag.NewFlagSet("greplive", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	var (
		pattern    = fs.String("pattern", "", "regex pattern to filter lines")
		levelStr   = fs.String("level", "any", "minimum severity level (debug|info|warn|error|fatal|any)")
		follow     = fs.Bool("follow", false, "tail the file and follow new lines")
		file       = fs.String("file", "", "file to read (defaults to stdin)")
		timestamp  = fs.Bool("timestamp", false, "prefix each line with current timestamp")
		rate       = fs.Int("rate", 0, "max lines per second (0 = unlimited)")
		dedupe     = fs.Bool("dedupe", false, "suppress consecutive duplicate lines")
		maxLen     = fs.Int("max-line-len", 0, "truncate lines longer than this (0 = disabled)")
		checkpoint = fs.String("checkpoint", "", "path to checkpoint file for resuming tailed files")
	)

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil, err
		}
		return nil, fmt.Errorf("flag parse: %w", err)
	}

	lvl, err := severity.ParseLevel(*levelStr)
	if err != nil {
		return nil, fmt.Errorf("invalid level %q: %w", *levelStr, err)
	}

	return &Flags{
		Pattern:        *pattern,
		Level:          lvl,
		Follow:         *follow,
		File:           *file,
		Timestamp:      *timestamp,
		Rate:           *rate,
		Dedupe:         *dedupe,
		MaxLineLen:     *maxLen,
		CheckpointFile: *checkpoint,
	}, nil
}
