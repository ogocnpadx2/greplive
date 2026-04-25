// Package splice merges multiple line-producing channels into a single
// output channel, preserving order on a best-effort (first-ready) basis.
package splice

import (
	"context"
	"sync"
)

// Splitter fans in N input channels into one output channel.
type Splitter struct {
	out chan string
}

// New creates a Splitter that merges all srcs into a single channel.
// The output channel is closed once every source channel is drained and
// the provided context is cancelled or all inputs are exhausted.
func New(ctx context.Context, srcs ...<-chan string) *Splitter {
	out := make(chan string, 64)
	var wg sync.WaitGroup

	for _, src := range srcs {
		wg.Add(1)
		go func(ch <-chan string) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case line, ok := <-ch:
					if !ok {
						return
					}
					select {
					case out <- line:
					case <-ctx.Done():
						return
					}
				}
			}
		}(src)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return &Splitter{out: out}
}

// Out returns the merged output channel.
func (s *Splitter) Out() <-chan string {
	return s.out
}
