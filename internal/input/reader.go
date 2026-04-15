// Package input provides utilities for reading log lines from
// various sources such as stdin, files, or arbitrary io.Reader streams.
package input

import (
	"bufio"
	"context"
	"io"
)

// LineReader streams lines from an io.Reader and sends them over a channel.
type LineReader struct {
	src io.Reader
}

// New creates a new LineReader that reads from src.
func New(src io.Reader) *LineReader {
	return &LineReader{src: src}
}

// Lines reads lines from the underlying source and sends each non-empty line
// to the returned channel. The channel is closed when the source is exhausted
// or the context is cancelled. Any read error other than io.EOF is sent to
// the errs channel before both channels are closed.
func (r *LineReader) Lines(ctx context.Context) (<-chan string, <-chan error) {
	lines := make(chan string)
	errs := make(chan error, 1)

	go func() {
		defer close(lines)
		defer close(errs)

		scanner := bufio.NewScanner(r.src)
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				return
			default:
			}

			line := scanner.Text()
			if line == "" {
				continue
			}

			select {
			case lines <- line:
			case <-ctx.Done():
				return
			}
		}

		if err := scanner.Err(); err != nil {
			select {
			case errs <- err:
			case <-ctx.Done():
			}
		}
	}()

	return lines, errs
}
