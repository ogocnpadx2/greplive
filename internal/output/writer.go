// Package output handles formatted, color-coded writing of filtered log lines
// to the terminal.
package output

import (
	"fmt"
	"io"
	"os"
	"time"

	"greplive/internal/severity"
)

// Writer wraps an io.Writer and provides formatted log line output.
type Writer struct {
	out       io.Writer
	showTime  bool
	showLevel bool
}

// Option configures a Writer.
type Option func(*Writer)

// WithTimestamp enables prepending a timestamp to each output line.
func WithTimestamp(enabled bool) Option {
	return func(w *Writer) {
		w.showTime = enabled
	}
}

// WithLevel enables prepending the detected severity level to each output line.
func WithLevel(enabled bool) Option {
	return func(w *Writer) {
		w.showLevel = enabled
	}
}

// New creates a new Writer writing to out with the given options.
// If out is nil, os.Stdout is used.
func New(out io.Writer, opts ...Option) *Writer {
	if out == nil {
		out = os.Stdout
	}
	w := &Writer{out: out}
	for _, o := range opts {
		o(w)
	}
	return w
}

// WriteLine writes a single log line, applying color based on severity and
// optionally prepending a timestamp and level label.
func (w *Writer) WriteLine(line string) {
	lvl := severity.Detect(line)
	colored := severity.Colorize(line, lvl)

	prefix := ""
	if w.showTime {
		prefix += fmt.Sprintf("[%s] ", time.Now().Format("15:04:05"))
	}
	if w.showLevel {
		prefix += fmt.Sprintf("[%-5s] ", lvl.String())
	}

	fmt.Fprintf(w.out, "%s%s\n", prefix, colored)
}
