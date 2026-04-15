// Package ratelimit provides a token-bucket rate limiter for controlling
// the throughput of log lines processed by greplive.
package ratelimit

import (
	"context"
	"time"
)

// Limiter controls how many lines per second are allowed through the pipeline.
type Limiter struct {
	tokens   chan struct{}
	rate     int
	stop     chan struct{}
}

// New creates a Limiter that allows up to linesPerSec lines per second.
// If linesPerSec is zero or negative, no rate limiting is applied.
func New(linesPerSec int) *Limiter {
	l := &Limiter{
		rate: linesPerSec,
		stop: make(chan struct{}),
	}
	if linesPerSec > 0 {
		l.tokens = make(chan struct{}, linesPerSec)
		go l.refill()
	}
	return l
}

// refill adds tokens to the bucket at the configured rate using a ticker.
func (l *Limiter) refill() {
	interval := time.Second / time.Duration(l.rate)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			select {
			case l.tokens <- struct{}{}:
			default:
				// bucket full, discard token
			}
		case <-l.stop:
			return
		}
	}
}

// Wait blocks until a token is available or the context is cancelled.
// Returns false if the context was cancelled before a token could be acquired.
func (l *Limiter) Wait(ctx context.Context) bool {
	if l.tokens == nil {
		return true
	}
	select {
	case <-l.tokens:
		return true
	case <-ctx.Done():
		return false
	}
}

// Stop shuts down the background refill goroutine.
func (l *Limiter) Stop() {
	if l.tokens != nil {
		close(l.stop)
	}
}

// Rate returns the configured lines-per-second limit (0 means unlimited).
func (l *Limiter) Rate() int {
	return l.rate
}
