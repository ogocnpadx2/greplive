// Package fanout distributes each input line to multiple independent output
// channels. It is useful when a single source of log lines must be consumed
// by several downstream stages simultaneously (e.g. writing to a file while
// also streaming to stdout).
package fanout

import "sync"

// Fanout reads lines from a single input channel and writes each line to every
// registered output channel. All outputs receive every line; slow consumers
// are not dropped — callers should buffer outputs appropriately.
type Fanout struct {
	mu      sync.RWMutex
	outputs []chan<- string
}

// New creates a Fanout that will forward lines from src to all channels
// added via Add. The forwarding goroutine runs until src is closed, after
// which all registered output channels are closed.
func New(src <-chan string) *Fanout {
	f := &Fanout{}
	go f.run(src)
	return f
}

// Add registers an output channel. The channel will receive every line that
// arrives on the source after the call returns. Add is safe to call
// concurrently, but channels added after the source has already been closed
// will be closed immediately.
func (f *Fanout) Add(ch chan<- string) {
	f.mu.Lock()
	f.outputs = append(f.outputs, ch)
	f.mu.Unlock()
}

// run is the internal forwarding loop.
func (f *Fanout) run(src <-chan string) {
	for line := range src {
		f.mu.RLock()
		for _, out := range f.outputs {
			out <- line
		}
		f.mu.RUnlock()
	}

	// Source exhausted — close all registered outputs.
	f.mu.Lock()
	for _, out := range f.outputs {
		close(out)
	}
	f.outputs = nil
	f.mu.Unlock()
}

// Pipe is a convenience helper that creates a buffered output channel,
// registers it with the Fanout, and returns it to the caller.
func (f *Fanout) Pipe(bufSize int) <-chan string {
	if bufSize < 0 {
		bufSize = 0
	}
	ch := make(chan string, bufSize)
	f.Add(ch)
	return ch
}
