// Package banner prints a startup header summarising the active greplive
// configuration so operators can confirm flags at a glance.
package banner

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// Config holds the values that will be rendered in the banner.
type Config struct {
	Pattern   string
	Level     string
	InputFile string
	RateLimit int
	Dedupe    bool
	Truncate  int
}

// Printer writes a banner to an io.Writer.
type Printer struct {
	w io.Writer
}

// New returns a Printer that writes to w. If w is nil, os.Stderr is used.
func New(w io.Writer) *Printer {
	if w == nil {
		w = os.Stderr
	}
	return &Printer{w: w}
}

// Print renders the banner for cfg.
func (p *Printer) Print(cfg Config) {
	fmt.Fprintln(p.w, strings.Repeat("─", 50))
	fmt.Fprintln(p.w, "  greplive – real-time log filter")
	fmt.Fprintln(p.w, strings.Repeat("─", 50))
	pattern := cfg.Pattern
	if pattern == "" {
		pattern = "(none)"
	}
	fmt.Fprintf(p.w, "  pattern   : %s\n", pattern)
	level := cfg.Level
	if level == "" {
		level = "all"
	}
	fmt.Fprintf(p.w, "  level     : %s\n", level)
	source := cfg.InputFile
	if source == "" {
		source = "stdin"
	}
	fmt.Fprintf(p.w, "  source    : %s\n", source)
	if cfg.RateLimit > 0 {
		fmt.Fprintf(p.w, "  rate-limit: %d lines/s\n", cfg.RateLimit)
	}
	if cfg.Dedupe {
		fmt.Fprintln(p.w, "  dedupe    : enabled")
	}
	if cfg.Truncate > 0 {
		fmt.Fprintf(p.w, "  truncate  : %d chars\n", cfg.Truncate)
	}
	fmt.Fprintln(p.w, strings.Repeat("─", 50))
}
