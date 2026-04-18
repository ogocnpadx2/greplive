package repeat

import (
	"testing"
	"time"
)

var fixed = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestEnabled_ZeroMax(t *testing.T) {
	r := New(0, time.Minute)
	if r.Enabled() {
		t.Fatal("expected disabled")
	}
}

func TestEnabled_PositiveMax(t *testing.T) {
	r := New(3, time.Minute)
	if !r.Enabled() {
		t.Fatal("expected enabled")
	}
}

func TestPush_Disabled_PassThrough(t *testing.T) {
	r := New(0, time.Minute)
	for i := 0; i < 10; i++ {
		out, emit := r.Push("hello")
		if !emit || out != "hello" {
			t.Fatalf("expected pass-through at i=%d", i)
		}
	}
}

func TestPush_WithinMax_AllEmitted(t *testing.T) {
	r := New(3, time.Minute).WithClock(fixedClock(fixed))
	for i := 0; i < 3; i++ {
		out, emit := r.Push("msg")
		if !emit || out != "msg" {
			t.Fatalf("expected emit at i=%d", i)
		}
	}
}

func TestPush_ExceedsMax_SummaryThenDrop(t *testing.T) {
	r := New(2, time.Minute).WithClock(fixedClock(fixed))
	r.Push("msg")
	r.Push("msg")

	out, emit := r.Push("msg") // 3rd — summary
	if !emit {
		t.Fatal("expected summary emitted")
	}
	if out == "msg" {
		t.Fatal("expected summary, got original line")
	}

	_, emit = r.Push("msg") // 4th — dropped
	if emit {
		t.Fatal("expected drop")
	}
}

func TestPush_LineChange_Resets(t *testing.T) {
	r := New(2, time.Minute).WithClock(fixedClock(fixed))
	r.Push("a")
	r.Push("a")
	r.Push("a") // suppressed

	out, emit := r.Push("b")
	if !emit || out != "b" {
		t.Fatal("expected new line to reset and emit")
	}
}

func TestPush_WindowExpiry_Resets(t *testing.T) {
	now := fixed
	r := New(1, time.Second).WithClock(func() time.Time { return now })
	r.Push("x")
	r.Push("x") // summary

	now = fixed.Add(2 * time.Second) // advance past window
	out, emit := r.Push("x")
	if !emit || out != "x" {
		t.Fatal("expected reset after window expiry")
	}
}

func TestFlush_ReturnsSummary(t *testing.T) {
	r := New(1, time.Minute).WithClock(fixedClock(fixed))
	r.Push("hi")
	r.Push("hi") // summary emitted inline
	r.Push("hi") // dropped

	s, ok := r.Flush()
	if !ok || s == "" {
		t.Fatal("expected flush summary")
	}
}

func TestFlush_NoSummaryWhenNotSuppressed(t *testing.T) {
	r := New(3, time.Minute).WithClock(fixedClock(fixed))
	r.Push("hi")
	_, ok := r.Flush()
	if ok {
		t.Fatal("expected no flush summary")
	}
}
