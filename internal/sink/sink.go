// Package sink provides writers that direct processed log lines to one or
// more output destinations such as files, stderr, or arbitrary io.Writers.
package sink

import (
	"fmt"
	"io"
	"os"
	"sync"
)

// Sink writes lines to an underlying io.Writer.
type Sink struct {
	mu sync.Mutex
	w  io.Writer
}

// New returns a Sink that writes to w. If w is nil, os.Stdout is used.
func New(w io.Writer) *Sink {
	if w == nil {
		w = os.Stdout
	}
	return &Sink{w: w}
}

// NewFile opens the file at path for appending (creating it if necessary) and
// returns a Sink that writes to it, along with a close function the caller
// must invoke when finished.
func NewFile(path string) (*Sink, func() error, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, nil, fmt.Errorf("sink: open %q: %w", path, err)
	}
	s := New(f)
	return s, f.Close, nil
}

// Write sends line to the underlying writer followed by a newline.
// It is safe for concurrent use.
func (s *Sink) Write(line string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, err := fmt.Fprintln(s.w, line)
	if err != nil {
		return fmt.Errorf("sink: write: %w", err)
	}
	return nil
}

// WriteAll sends each line in lines to the underlying writer in order.
func (s *Sink) WriteAll(lines []string) error {
	for _, l := range lines {
		if err := s.Write(l); err != nil {
			return err
		}
	}
	return nil
}

// Multi fans a single Write call out to multiple sinks. The first error
// encountered is returned; remaining sinks are still attempted.
func Multi(sinks ...*Sink) *multiSink {
	return &multiSink{sinks: sinks}
}

type multiSink struct {
	sinks []*Sink
}

// Write sends line to every contained Sink.
func (m *multiSink) Write(line string) error {
	var first error
	for _, s := range m.sinks {
		if err := s.Write(line); err != nil && first == nil {
			first = err
		}
	}
	return first
}
