// Package batch groups consecutive log lines into fixed-size or time-bounded
// batches before passing them downstream.
package batch

import (
	"sync"
	"time"
)

// Batch holds a slice of lines collected within one window.
type Batch struct {
	Lines []string
}

// Batcher accumulates lines and flushes them as a Batch when either the
// maximum size is reached or the flush interval elapses.
type Batcher struct {
	mu       sync.Mutex
	buf      []string
	maxSize  int
	interval time.Duration
	out      chan Batch
	stop     chan struct{}
	wg       sync.WaitGroup
}

// New creates a Batcher that emits batches of up to maxSize lines or flushes
// every interval. maxSize <= 0 disables size-based flushing.
func New(maxSize int, interval time.Duration) *Batcher {
	if interval <= 0 {
		interval = time.Second
	}
	b := &Batcher{
		maxSize:  maxSize,
		interval: interval,
		out:      make(chan Batch, 16),
		stop:     make(chan struct{}),
	}
	b.wg.Add(1)
	go b.ticker()
	return b
}

// Push adds a line to the current batch. If the batch reaches maxSize it is
// flushed immediately.
func (b *Batcher) Push(line string) {
	b.mu.Lock()
	b.buf = append(b.buf, line)
	flush := b.maxSize > 0 && len(b.buf) >= b.maxSize
	b.mu.Unlock()
	if flush {
		b.flush()
	}
}

// Out returns the channel on which completed batches are delivered.
func (b *Batcher) Out() <-chan Batch { return b.out }

// Stop flushes any remaining lines and shuts down the background ticker.
func (b *Batcher) Stop() {
	close(b.stop)
	b.wg.Wait()
	b.flush()
	close(b.out)
}

func (b *Batcher) flush() {
	b.mu.Lock()
	if len(b.buf) == 0 {
		b.mu.Unlock()
		return
	}
	lines := make([]string, len(b.buf))
	copy(lines, b.buf)
	b.buf = b.buf[:0]
	b.mu.Unlock()
	b.out <- Batch{Lines: lines}
}

func (b *Batcher) ticker() {
	defer b.wg.Done()
	t := time.NewTicker(b.interval)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			b.flush()
		case <-b.stop:
			return
		}
	}
}
