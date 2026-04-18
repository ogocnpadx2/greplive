// Package tail provides a bounded ring-buffer of the last N lines seen,
// suitable for displaying trailing context when a match occurs.
package tail

import "sync"

// Buffer holds the last N lines written to it.
type Buffer struct {
	mu   sync.Mutex
	lines []string
	cap  int
	pos  int
	full bool
}

// New returns a Buffer that retains the last n lines.
// If n <= 0 it is clamped to 1.
func New(n int) *Buffer {
	if n <= 0 {
		n = 1
	}
	return &Buffer{lines: make([]string, n), cap: n}
}

// Push adds a line to the buffer, evicting the oldest when full.
func (b *Buffer) Push(line string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.lines[b.pos] = line
	b.pos = (b.pos + 1) % b.cap
	if !b.full && b.pos == 0 {
		b.full = true
	}
}

// Lines returns the buffered lines in insertion order (oldest first).
func (b *Buffer) Lines() []string {
	b.mu.Lock()
	defer b.mu.Unlock()
	if !b.full {
		out := make([]string, b.pos)
		copy(out, b.lines[:b.pos])
		return out
	}
	out := make([]string, b.cap)
	copy(out, b.lines[b.pos:])
	copy(out[b.cap-b.pos:], b.lines[:b.pos])
	return out
}

// Len returns the number of lines currently held.
func (b *Buffer) Len() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.full {
		return b.cap
	}
	return b.pos
}

// Reset clears all buffered lines.
func (b *Buffer) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()
	for i := range b.lines {
		b.lines[i] = ""
	}
	b.pos = 0
	b.full = false
}
