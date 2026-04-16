// Package multiline provides support for collapsing multi-line log entries
// (e.g. stack traces) into a single logical line.
package multiline

import (
	"regexp"
	"strings"
)

// Collector accumulates lines that belong to the same logical event.
type Collector struct {
	start   *regexp.Regexp
	continue_ *regexp.Regexp
	max     int
	buf     []string
}

// New creates a Collector.
// startPattern marks the first line of a new event.
// continuePattern matches continuation lines (e.g. stack frames).
// maxLines caps how many lines are merged (0 = unlimited).
func New(startPattern, continuePattern string, maxLines int) (*Collector, error) {
	start, err := regexp.Compile(startPattern)
	if err != nil {
		return nil, err
	}
	var cont *regexp.Regexp
	if continuePattern != "" {
		cont, err = regexp.Compile(continuePattern)
		if err != nil {
			return nil, err
		}
	}
	return &Collector{start: start, continue_: cont, max: maxLines}, nil
}

// Push feeds a line to the collector.
// It returns a completed event string and true when the buffer is flushed,
// or "", false when the line has been buffered.
func (c *Collector) Push(line string) (string, bool) {
	if c.start.MatchString(line) {
		event := c.flush()
		c.buf = append(c.buf, line)
		if event != "" {
			return event, true
		}
		return "", false
	}
	if c.isContinuation(line) && len(c.buf) > 0 {
		if c.max == 0 || len(c.buf) < c.max {
			c.buf = append(c.buf, line)
			return "", false
		}
	}
	// Not a continuation and not a start — flush then emit standalone.
	event := c.flush()
	if event != "" {
		c.buf = append(c.buf, line)
		return event, true
	}
	return line, true
}

// Flush returns any buffered lines as a single event.
func (c *Collector) Flush() string { return c.flush() }

func (c *Collector) flush() string {
	if len(c.buf) == 0 {
		return ""
	}
	event := strings.Join(c.buf, "\n")
	c.buf = c.buf[:0]
	return event
}

func (c *Collector) isContinuation(line string) bool {
	if c.continue_ != nil {
		return c.continue_.MatchString(line)
	}
	// Default: indented lines are continuations.
	return len(line) > 0 && (line[0] == '\t' || line[0] == ' ')
}
