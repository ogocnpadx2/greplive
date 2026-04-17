// Package contextline captures N lines before and after a match,
// similar to grep -A / -B / -C flags.
package contextline

// Collector holds a rolling pre-match buffer and emits lines with context.
type Collector struct {
	before int
	after  int
	pre    []string
	pending int // lines of post-context still to emit
}

// New creates a Collector. before is the number of lines to keep before a
// match; after is the number of lines to emit after a match.
func New(before, after int) *Collector {
	if before < 0 {
		before = 0
	}
	if after < 0 {
		after = 0
	}
	return &Collector{before: before, after: after}
}

// Feed accepts a line and whether it matched the filter. It returns the lines
// that should be emitted (may be empty, the line itself, or context lines).
func (c *Collector) Feed(line string, matched bool) []string {
	var out []string

	if matched {
		// Emit buffered pre-context lines.
		out = append(out, c.pre...)
		c.pre = c.pre[:0]
		out = append(out, line)
		c.pending = c.after
		return out
	}

	if c.pending > 0 {
		out = append(out, line)
		c.pending--
		return out
	}

	// Buffer for pre-context.
	if c.before > 0 {
		c.pre = append(c.pre, line)
		if len(c.pre) > c.before {
			c.pre = c.pre[len(c.pre)-c.before:]
		}
	}
	return nil
}

// Reset clears internal state.
func (c *Collector) Reset() {
	c.pre = c.pre[:0]
	c.pending = 0
}
