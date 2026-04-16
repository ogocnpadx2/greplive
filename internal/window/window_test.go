package window

import (
	"testing"
	"time"
)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestNew_InitialisesEmpty(t *testing.T) {
	w := New(5 * time.Second)
	if got := w.Count(); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestAdd_IncrementsCount(t *testing.T) {
	w := New(10 * time.Second)
	w.Add()
	w.Add()
	if got := w.Count(); got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}
}

func TestEviction_StaleEntriesRemoved(t *testing.T) {
	base := time.Now()
	w := New(5 * time.Second)

	// Add two events at t=0
	w.now = fixedClock(base)
	w.Add()
	w.Add()

	// Advance clock beyond window
	w.now = fixedClock(base.Add(6 * time.Second))
	if got := w.Count(); got != 0 {
		t.Fatalf("expected 0 after eviction, got %d", got)
	}
}

func TestEviction_KeepsRecentEntries(t *testing.T) {
	base := time.Now()
	w := New(10 * time.Second)

	w.now = fixedClock(base)
	w.Add()

	w.now = fixedClock(base.Add(8 * time.Second))
	w.Add()

	// Advance 9s: first event is within 10s window from new now (base+9)
	w.now = fixedClock(base.Add(9 * time.Second))
	if got := w.Count(); got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}
}

func TestRate_ZeroDuration(t *testing.T) {
	w := New(0)
	w.Add()
	if r := w.Rate(); r != 0 {
		t.Fatalf("expected 0 rate for zero duration, got %f", r)
	}
}

func TestRate_Calculation(t *testing.T) {
	base := time.Now()
	w := New(2 * time.Second)
	w.now = fixedClock(base)
	w.Add()
	w.Add()

	// 2 events / 2 seconds = 1.0
	if r := w.Rate(); r != 1.0 {
		t.Fatalf("expected rate 1.0, got %f", r)
	}
}

func TestReset_ClearsAllEntries(t *testing.T) {
	w := New(10 * time.Second)
	w.Add()
	w.Add()
	w.Reset()
	if got := w.Count(); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
}

func TestZeroDuration_NoEviction(t *testing.T) {
	base := time.Now()
	w := New(0)
	w.now = fixedClock(base)
	w.Add()
	w.now = fixedClock(base.Add(1 * time.Hour))
	w.Add()
	if got := w.Count(); got != 2 {
		t.Fatalf("expected 2 with no eviction, got %d", got)
	}
}
