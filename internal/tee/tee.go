// Package tee fans a stream of log lines out to multiple consumers.
package tee

import "sync"

// Tee multiplexes a single input channel to N output channels.
type Tee struct {
	inputs  []<-chan string
	outs    []chan string
	wg      sync.WaitGroup
}

// New creates a Tee that reads from src and writes to n output channels.
// Each output channel has the given buffer size.
func New(src <-chan string, n, bufSize int) *Tee {
	if n <= 0 {
		n = 1
	}
	outs := make([]chan string, n)
	for i := range outs {
		outs[i] = make(chan string, bufSize)
	}
	t := &Tee{outs: outs}
	t.wg.Add(1)
	go func() {
		defer t.wg.Done()
		defer func() {
			for _, o := range outs {
				close(o)
			}
		}()
		for line := range src {
			for _, o := range outs {
				o <- line
			}
		}
	}()
	return t
}

// Outputs returns the read-only output channels.
func (t *Tee) Outputs() []<-chan string {
	ro := make([]<-chan string, len(t.outs))
	for i, o := range t.outs {
		ro[i] = o
	}
	return ro
}

// Wait blocks until the fan-out goroutine has finished.
func (t *Tee) Wait() { t.wg.Wait() }
