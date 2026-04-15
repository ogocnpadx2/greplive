// Package input provides utilities for reading log lines from various sources.
package input

import (
	"context"
	"io"
	"os"
	"time"
)

// TailOptions configures the behavior of the Tail function.
type TailOptions struct {
	// PollInterval is how often to check for new data when the reader is exhausted.
	// Defaults to 250ms if zero.
	PollInterval time.Duration
}

// defaultPollInterval is used when TailOptions.PollInterval is not set.
const defaultPollInterval = 250 * time.Millisecond

// TailFile opens the named file and streams new lines as they are appended,
// similar to `tail -f`. It sends each non-empty line to the returned channel.
// The channel is closed when ctx is cancelled or a read error occurs.
func TailFile(ctx context.Context, path string, opts TailOptions) (<-chan string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	if opts.PollInterval == 0 {
		opts.PollInterval = defaultPollInterval
	}

	ch := make(chan string, 64)
	go func() {
		defer close(ch)
		defer f.Close()
		tailStream(ctx, f, opts.PollInterval, ch)
	}()

	return ch, nil
}

// tailStream reads lines from r, polling when EOF is reached, and sends them
// to ch until ctx is cancelled.
func tailStream(ctx context.Context, r io.Reader, poll time.Duration, ch chan<- string) {
	reader := New(r)
	ticker := time.NewTicker(poll)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		lineCh := reader.Lines(ctx)
		gotLine := false
		for line := range lineCh {
			gotLine = true
			select {
			case ch <- line:
			case <-ctx.Done():
				return
			}
		}

		if !gotLine {
			// No new data; wait before polling again.
			select {
			case <-ticker.C:
			case <-ctx.Done():
				return
			}
		}
	}
}
