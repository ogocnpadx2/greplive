package linecount

import (
	"testing"
	"time"
)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestNew_DefaultsBuckets(t *testing.T) {
	c := New(0, time.Second)
	if c.size != 1 {
		t.Fatalf("expected size 1, got %d", c.size)
	}
}

func TestInc_IncrementsRate(t *testing.T) {
	now := time.Now()
	c := New(10, time.Second)
	c.clock = fixedClock(now)
	c.Inc()
	c.Inc()
	c.Inc()
	if got := c.Rate(); got != 3 {
		t.Fatalf("expected 3, got %d", got)
	}
}

func TestRate_StaleEntriesNotCounted(t *testing.T) {
	now := time.Now()
	c := New(10, time.Second)
	c.clock = fixedClock(now)
	c.Inc()
	c.Inc()
	// advance past window
	c.clock = fixedClock(now.Add(2 * time.Second))
	if got := c.Rate(); got != 0 {
		t.Fatalf("expected 0 after window, got %d", got)
	}
}

func TestReset_ZerosCounters(t *testing.T) {
	now := time.Now()
	c := New(10, time.Second)
	c.clock = fixedClock(now)
	c.Inc()
	c.Inc()
	c.Reset()
	if got := c.Rate(); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
}

func TestRate_MultipleBuckets(t *testing.T) {
	base := time.Unix(1_000_000, 0)
	c := New(10, time.Second)
	bucketDur := c.bucketDuration()

	c.clock = fixedClock(base)
	c.Inc()

	c.clock = fixedClock(base.Add(bucketDur))
	c.Inc()

	// both within window
	c.clock = fixedClock(base.Add(bucketDur))
	if got := c.Rate(); got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}
}
