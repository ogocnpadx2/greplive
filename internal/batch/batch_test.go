package batch

import (
	"testing"
	"time"
)

func TestNew_DefaultInterval(t *testing.T) {
	b := New(10, 0)
	defer b.Stop()
	if b.interval != time.Second {
		t.Fatalf("expected 1s default interval, got %v", b.interval)
	}
}

func TestPush_BelowMaxSize_NotFlushed(t *testing.T) {
	b := New(5, time.Minute)
	defer b.Stop()
	b.Push("line1")
	b.Push("line2")
	select {
	case batch := <-b.Out():
		t.Fatalf("unexpected early flush: %v", batch)
	default:
	}
}

func TestPush_AtMaxSize_Flushed(t *testing.T) {
	b := New(3, time.Minute)
	defer b.Stop()
	b.Push("a")
	b.Push("b")
	b.Push("c")
	select {
	case batch := <-b.Out():
		if len(batch.Lines) != 3 {
			t.Fatalf("expected 3 lines, got %d", len(batch.Lines))
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("expected flush did not arrive")
	}
}

func TestStop_FlushesRemaining(t *testing.T) {
	b := New(100, time.Minute)
	b.Push("x")
	b.Push("y")
	b.Stop()
	var got []Batch
	for batch := range b.Out() {
		got = append(got, batch)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 batch, got %d", len(got))
	}
	if len(got[0].Lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(got[0].Lines))
	}
}

func TestTicker_FlushesOnInterval(t *testing.T) {
	b := New(100, 50*time.Millisecond)
	defer b.Stop()
	b.Push("tick")
	select {
	case batch := <-b.Out():
		if len(batch.Lines) != 1 || batch.Lines[0] != "tick" {
			t.Fatalf("unexpected batch: %v", batch)
		}
	case <-time.After(300 * time.Millisecond):
		t.Fatal("ticker did not fire")
	}
}

func TestPush_ZeroMaxSize_NeverSizeFlushes(t *testing.T) {
	b := New(0, time.Minute)
	defer b.Stop()
	for i := 0; i < 50; i++ {
		b.Push("line")
	}
	select {
	case batch := <-b.Out():
		t.Fatalf("unexpected flush with zero max size: %v", batch)
	default:
	}
}
