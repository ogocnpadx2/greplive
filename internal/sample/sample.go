// Package sample provides probabilistic line sampling for high-volume log streams.
package sample

import (
	"math/rand"
	"sync"
)

// Sampler drops lines probabilistically, keeping approximately Rate fraction of
// all lines. A Rate of 1.0 keeps every line; 0.0 drops all lines.
type Sampler struct {
	mu   sync.Mutex
	rate float64
	rng  *rand.Rand
}

// New returns a Sampler with the given keep rate clamped to [0.0, 1.0].
// If rate >= 1.0 the sampler is effectively disabled (all lines pass).
func New(rate float64, seed int64) *Sampler {
	if rate < 0 {
		rate = 0
	}
	if rate > 1 {
		rate = 1
	}
	return &Sampler{
		rate: rate,
		rng:  rand.New(rand.NewSource(seed)),
	}
}

// Rate returns the configured keep rate.
func (s *Sampler) Rate() float64 {
	return s.rate
}

// Enabled reports whether sampling is active (rate < 1.0).
func (s *Sampler) Enabled() bool {
	return s.rate < 1.0
}

// Keep returns true if the line should be forwarded downstream.
func (s *Sampler) Keep() bool {
	if !s.Enabled() {
		return true
	}
	s.mu.Lock()
	v := s.rng.Float64()
	s.mu.Unlock()
	return v < s.rate
}
