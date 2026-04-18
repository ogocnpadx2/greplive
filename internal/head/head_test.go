package head

import "testing"

func TestNew_Disabled_WhenZero(t *testing.T) {
	l := New(0)
	if l.Enabled() {
		t.Fatal("expected disabled for max=0")
	}
}

func TestNew_Disabled_WhenNegative(t *testing.T) {
	l := New(-5)
	if l.Enabled() {
		t.Fatal("expected disabled for max=-5")
	}
}

func TestNew_Enabled_WhenPositive(t *testing.T) {
	l := New(3)
	if !l.Enabled() {
		t.Fatal("expected enabled for max=3")
	}
}

func TestAllow_Disabled_AlwaysEmits(t *testing.T) {
	l := New(0)
	for i := 0; i < 100; i++ {
		emit, done := l.Allow()
		if !emit || done {
			t.Fatalf("iteration %d: expected emit=true done=false", i)
		}
	}
}

func TestAllow_EmitsUpToMax(t *testing.T) {
	l := New(3)
	for i := 0; i < 3; i++ {
		emit, _ := l.Allow()
		if !emit {
			t.Fatalf("line %d should be emitted", i+1)
		}
	}
}

func TestAllow_DoneOnLastLine(t *testing.T) {
	l := New(2)
	l.Allow()
	_, done := l.Allow()
	if !done {
		t.Fatal("expected done=true on last allowed line")
	}
}

func TestAllow_BlocksAfterMax(t *testing.T) {
	l := New(2)
	l.Allow()
	l.Allow()
	emit, done := l.Allow()
	if emit || !done {
		t.Fatal("expected emit=false done=true after max")
	}
}

func TestReset_ResetsCounter(t *testing.T) {
	l := New(1)
	l.Allow()
	l.Reset()
	emit, _ := l.Allow()
	if !emit {
		t.Fatal("expected emit after reset")
	}
}

func TestRemaining_Disabled(t *testing.T) {
	l := New(0)
	if l.Remaining() != -1 {
		t.Fatal("expected -1 for disabled limiter")
	}
}

func TestRemaining_DecreasesWithUse(t *testing.T) {
	l := New(3)
	if l.Remaining() != 3 {
		t.Fatalf("expected 3, got %d", l.Remaining())
	}
	l.Allow()
	if l.Remaining() != 2 {
		t.Fatalf("expected 2, got %d", l.Remaining())
	}
}
