package stats_test

import (
	"testing"
	"time"

	"greplive/internal/stats"
)

func TestNew_InitialisesCounters(t *testing.T) {
	c := stats.New()
	if c.LinesRead.Load() != 0 {
		t.Fatalf("expected 0 lines read, got %d", c.LinesRead.Load())
	}
	if c.StartTime.IsZero() {
		t.Fatal("expected non-zero start time")
	}
}

func TestIncrRead(t *testing.T) {
	c := stats.New()
	c.IncrRead()
	c.IncrRead()
	if got := c.LinesRead.Load(); got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}
}

func TestIncrMatched(t *testing.T) {
	c := stats.New()
	c.IncrMatched()
	if got := c.LinesMatched.Load(); got != 1 {
		t.Fatalf("expected 1, got %d", got)
	}
}

func TestIncrDropped(t *testing.T) {
	c := stats.New()
	c.IncrDropped()
	c.IncrDropped()
	c.IncrDropped()
	if got := c.LinesDropped.Load(); got != 3 {
		t.Fatalf("expected 3, got %d", got)
	}
}

func TestIncrSeverity(t *testing.T) {
	c := stats.New()
	c.IncrSeverity("ERROR")
	c.IncrSeverity("ERROR")
	c.IncrSeverity("WARN")

	counts := c.SeverityCounts()
	if counts["ERROR"] != 2 {
		t.Fatalf("expected ERROR=2, got %d", counts["ERROR"])
	}
	if counts["WARN"] != 1 {
		t.Fatalf("expected WARN=1, got %d", counts["WARN"])
	}
}

func TestSeverityCounts_ReturnsCopy(t *testing.T) {
	c := stats.New()
	c.IncrSeverity("INFO")
	copy1 := c.SeverityCounts()
	copy1["INFO"] = 999
	copy2 := c.SeverityCounts()
	if copy2["INFO"] != 1 {
		t.Fatal("SeverityCounts should return an independent copy")
	}
}

func TestElapsed(t *testing.T) {
	c := stats.New()
	time.Sleep(10 * time.Millisecond)
	if c.Elapsed() < 10*time.Millisecond {
		t.Fatal("expected elapsed >= 10ms")
	}
}

func TestSnapshot(t *testing.T) {
	c := stats.New()
	c.IncrRead()
	c.IncrRead()
	c.IncrMatched()
	c.IncrDropped()
	c.IncrSeverity("DEBUG")

	s := c.Snapshot()
	if s.LinesRead != 2 {
		t.Fatalf("snapshot LinesRead: expected 2, got %d", s.LinesRead)
	}
	if s.LinesMatched != 1 {
		t.Fatalf("snapshot LinesMatched: expected 1, got %d", s.LinesMatched)
	}
	if s.LinesDropped != 1 {
		t.Fatalf("snapshot LinesDropped: expected 1, got %d", s.LinesDropped)
	}
	if s.SeverityCounts["DEBUG"] != 1 {
		t.Fatalf("snapshot severity DEBUG: expected 1, got %d", s.SeverityCounts["DEBUG"])
	}
	if s.Elapsed <= 0 {
		t.Fatal("snapshot Elapsed should be positive")
	}
}
