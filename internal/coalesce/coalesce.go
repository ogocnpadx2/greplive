// Package coalesce merges bursts of identical log lines into a single
// representative line with a repeat count annotation.
package coalesce

import (
	"fmt"
	"sync"
	"time"
)

// Coalescer buffers repeated lines and emits a summary when the burst ends.
type Coalescer struct {
	mu       sync.Mutex
	window   time.Duration
	last     string
	count    int
	timer    *time.Timer
	emit     func(string)
}

// New returns a Coalescer that groups repeated lines arriving within window.
// emit is called with the final (possibly annotated) line.
func New(window time.Duration, emit func(string)) *Coalescer {
	if window <= 0 {
		window = 200 * time.Millisecond
	}
	return &Coalescer{window: window, emit: emit}
}

// Push submits a line for coalescing.
func (c *Coalescer) Push(line string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if line == c.last {
		c.count++
		if c.timer != nil {
			c.timer.Reset(c.window)
		}
		return
	}

	// Different line — flush previous burst first.
	c.flushLocked()

	c.last = line
	c.count = 1
	c.timer = time.AfterFunc(c.window, func() {
		c.mu.Lock()
		defer c.mu.Unlock()
		c.flushLocked()
	})
}

// Flush forces any buffered line to be emitted immediately.
func (c *Coalescer) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.flushLocked()
}

func (c *Coalescer) flushLocked() {
	if c.last == "" {
		return
	}
	if c.timer != nil {
		c.timer.Stop()
		c.timer = nil
	}
	line := c.last
	if c.count > 1 {
		line = fmt.Sprintf("%s  [x%d]", c.last, c.count)
	}
	c.last = ""
	c.count = 0
	c.emit(line)
}
