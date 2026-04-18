// Package linecount provides a simple per-second line rate tracker.
package linecount

import (
	"sync"
	"time"
)

// Counter tracks lines seen within a sliding window of buckets.
type Counter struct {
	mu       sync.Mutex
	buckets  []int64
	times    []time.Time
	size     int
	window   time.Duration
	clock    func() time.Time
}

// New creates a Counter with the given number of buckets spanning the window duration.
func New(buckets int, window time.Duration) *Counter {
	if buckets < 1 {
		buckets = 1
	}
	return &Counter{
		buckets: make([]int64, buckets),
		times:   make([]time.Time, buckets),
		size:    buckets,
		window:  window,
		clock:   time.Now,
	}
}

// Inc records one line at the current time.
func (c *Counter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := c.clock()
	idx := c.bucketIndex(now)
	if c.times[idx].IsZero() || now.Sub(c.times[idx]) >= c.bucketDuration() {
		c.buckets[idx] = 0
		c.times[idx] = now
	}
	c.buckets[idx]++
}

// Rate returns the total count across all valid buckets.
func (c *Counter) Rate() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := c.clock()
	var total int64
	for i := 0; i < c.size; i++ {
		if !c.times[i].IsZero() && now.Sub(c.times[i]) < c.window {
			total += c.buckets[i]
		}
	}
	return total
}

// Reset zeroes all buckets.
func (c *Counter) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for i := range c.buckets {
		c.buckets[i] = 0
		c.times[i] = time.Time{}
	}
}

func (c *Counter) bucketDuration() time.Duration {
	return c.window / time.Duration(c.size)
}

func (c *Counter) bucketIndex(t time.Time) int {
	return int(t.UnixNano()/int64(c.bucketDuration())) % c.size
}
