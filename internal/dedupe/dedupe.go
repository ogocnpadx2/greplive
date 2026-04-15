// Package dedupe provides line deduplication for log streams.
// It suppresses consecutive identical lines and optionally reports
// how many times a repeated line was seen before a new one appears.
package dedupe

import "fmt"
import "sync"

// Deduper tracks consecutive duplicate lines and optionally emits
// a suppression summary when the repeated sequence ends.
type Deduper struct {
	mu      sync.Mutex
	last    string
	count   int
	summary bool
}

// New creates a new Deduper. When summary is true, Check will return a
// non-empty flush message whenever a run of duplicates ends.
func New(summary bool) *Deduper {
	return &Deduper{summary: summary}
}

// Check evaluates line against the previous line.
// It returns:
//
//	flush  – a summary string to emit before the new line (may be empty)
//	allow  – whether line itself should be forwarded downstream
func (d *Deduper) Check(line string) (flush string, allow bool) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if line == d.last {
		d.count++
		return "", false
	}

	if d.summary && d.count > 1 {
		flush = fmt.Sprintf("  [repeated %d times]", d.count)
	}

	d.last = line
	d.count = 1
	return flush, true
}

// Flush returns (and resets) any pending summary for the current run.
// Call this when the stream ends to surface a trailing duplicate run.
func (d *Deduper) Flush() string {
	d.mu.Lock()
	defer d.mu.Unlock()

	var s string
	if d.summary && d.count > 1 {
		s = fmt.Sprintf("  [repeated %d times]", d.count)
	}
	d.count = 0
	return s
}

// Reset clears all state, useful between log sources.
func (d *Deduper) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.last = ""
	d.count = 0
}
