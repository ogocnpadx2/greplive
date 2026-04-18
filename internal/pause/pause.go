// Package pause provides a toggle that allows the pipeline to be
// paused and resumed without dropping the underlying source.
package pause

import "sync"

// Pauser holds the paused/running state for a pipeline stage.
type Pauser struct {
	mu     sync.Mutex
	cond   *sync.Cond
	paused bool
}

// New returns a new Pauser in the running state.
func New() *Pauser {
	p := &Pauser{}
	p.cond = sync.NewCond(&p.mu)
	return p
}

// Pause suspends callers of Wait until Resume is called.
func (p *Pauser) Pause() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.paused = true
}

// Resume unblocks any callers waiting in Wait.
func (p *Pauser) Resume() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.paused = false
	p.cond.Broadcast()
}

// Toggle flips the current paused state.
func (p *Pauser) Toggle() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.paused = !p.paused
	if !p.paused {
		p.cond.Broadcast()
	}
}

// Paused reports whether the Pauser is currently paused.
func (p *Pauser) Paused() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.paused
}

// Wait blocks the caller while the Pauser is in the paused state.
func (p *Pauser) Wait() {
	p.mu.Lock()
	defer p.mu.Unlock()
	for p.paused {
		p.cond.Wait()
	}
}
