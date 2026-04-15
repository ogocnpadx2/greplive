// Package buffer provides a fixed-size ring buffer for storing recent log lines.
// It is useful for capturing context around matched lines (e.g. --before-context
// and --after-context style behaviour).
package buffer

import "sync"

// Ring is a thread-safe circular buffer that retains the last N lines.
type Ring struct {
	mu   sync.Mutex
	data []string
	size int
	head int
	count int
}

// New creates a new Ring buffer that holds at most size lines.
// If size is zero or negative the buffer is effectively disabled (capacity 0).
func New(size int) *Ring {
	if size < 0 {
		size = 0
	}
	return &Ring{
		data: make([]string, size),
		size: size,
	}
}

// Push adds a line to the ring buffer, overwriting the oldest entry when full.
func (r *Ring) Push(line string) {
	if r.size == 0 {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[r.head] = line
	r.head = (r.head + 1) % r.size
	if r.count < r.size {
		r.count++
	}
}

// Snapshot returns the buffered lines in insertion order (oldest first).
// The returned slice is a copy; mutations do not affect the ring.
func (r *Ring) Snapshot() []string {
	if r.size == 0 {
		return nil
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]string, r.count)
	start := (r.head - r.count + r.size) % r.size
	for i := 0; i < r.count; i++ {
		out[i] = r.data[(start+i)%r.size]
	}
	return out
}

// Len returns the number of lines currently stored.
func (r *Ring) Len() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.count
}

// Cap returns the maximum number of lines the buffer can hold.
func (r *Ring) Cap() int {
	return r.size
}

// Reset clears all stored lines without reallocating.
func (r *Ring) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.head = 0
	r.count = 0
}
