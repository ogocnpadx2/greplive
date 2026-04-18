// Package offset tracks byte offsets within a stream, enabling resume-from-position
// behaviour when combined with the checkpoint package.
package offset

import "sync/atomic"

// Tracker records how many bytes have been consumed from a source.
type Tracker struct {
	bytes atomic.Int64
}

// New returns a new Tracker starting at the given offset.
func New(initial int64) *Tracker {
	t := &Tracker{}
	t.bytes.Store(initial)
	return t
}

// Add increments the tracked offset by n bytes and returns the new total.
func (t *Tracker) Add(n int) int64 {
	return t.bytes.Add(int64(n))
}

// Get returns the current byte offset.
func (t *Tracker) Get() int64 {
	return t.bytes.Load()
}

// Reset sets the offset back to zero.
func (t *Tracker) Reset() {
	t.bytes.Store(0)
}

// Snapshot returns a point-in-time copy of the current offset.
func (t *Tracker) Snapshot() int64 {
	return t.bytes.Load()
}
