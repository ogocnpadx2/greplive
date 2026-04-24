// Package reopen provides a file-backed io.Reader that transparently
// reopens the underlying file when a rotation signal is received.
package reopen

import (
	"context"
	"io"
	"os"
	"sync"
)

// Reader wraps an *os.File and reopens it on demand.
type Reader struct {
	mu   sync.Mutex
	path string
	f    *os.File
}

// New opens path and returns a Reader ready for use.
func New(path string) (*Reader, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return &Reader{path: path, f: f}, nil
}

// Reopen closes the current file handle and opens a fresh one,
// seeking to the beginning. It is safe to call concurrently.
func (r *Reader) Reopen() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.f != nil {
		_ = r.f.Close()
	}
	f, err := os.Open(r.path)
	if err != nil {
		return err
	}
	r.f = f
	return nil
}

// Read implements io.Reader.
func (r *Reader) Read(p []byte) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.f == nil {
		return 0, io.ErrClosedPipe
	}
	return r.f.Read(p)
}

// Close closes the underlying file.
func (r *Reader) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.f == nil {
		return nil
	}
	err := r.f.Close()
	r.f = nil
	return err
}

// WatchAndReopen listens on rotated and calls Reopen each time a
// signal arrives. It returns when ctx is cancelled.
func (r *Reader) WatchAndReopen(ctx context.Context, rotated <-chan struct{}) {
	for {
		select {
		case <-ctx.Done():
			return
		case _, ok := <-rotated:
			if !ok {
				return
			}
			_ = r.Reopen()
		}
	}
}
