package clamp

import (
	"testing"
	"time"
)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestNew_DisabledWhenEmptyPattern(t *testing.T) {
	c, err := New("", 5, time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Enabled() {
		t.Fatal("expected disabled")
	}
}

func TestNew_DisabledWhenZeroMax(t *testing.T) {
	c, err := New("error", 0, time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Enabled() {
		t.Fatal("expected disabled")
	}
}

func TestNew_InvalidPattern_ReturnsError(t *testing.T) {
	_, err := New("[invalid", 3, time.Minute)
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestNew_ValidPattern_Enabled(t *testing.T) {
	c, err := New("error", 3, time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !c.Enabled() {
		t.Fatal("expected enabled")
	}
}

func TestAllow_NonMatchingLine_AlwaysPasses(t *testing.T) {
	c, _ := newWithClock("error", 2, time.Minute, fixedClock(time.Now()))
	for i := 0; i < 10; i++ {
		if !c.Allow("info: all good") {
			t.Fatalf("non-matching line should always pass (iteration %d)", i)
		}
	}
}

func TestAllow_WithinMax_AllEmitted(t *testing.T) {
	now := time.Now()
	c, _ := newWithClock("error", 3, time.Minute, fixedClock(now))
	for i := 0; i < 3; i++ {
		if !c.Allow("error: something bad") {
			t.Fatalf("line %d should be allowed", i)
		}
	}
}

func TestAllow_ExceedsMax_Dropped(t *testing.T) {
	now := time.Now()
	c, _ := newWithClock("error", 3, time.Minute, fixedClock(now))
	for i := 0; i < 3; i++ {
		c.Allow("error: something bad")
	}
	if c.Allow("error: one more") {
		t.Fatal("4th match within window should be dropped")
	}
}

func TestAllow_WindowAdvances_ResetsCapacity(t *testing.T) {
	base := time.Now()
	clock := fixedClock(base)
	c, _ := newWithClock("error", 2, time.Second, clock)
	c.Allow("error: a")
	c.Allow("error: b")
	if c.Allow("error: c") {
		t.Fatal("3rd match should be dropped within window")
	}
	// advance clock beyond window
	c.clock = fixedClock(base.Add(2 * time.Second))
	if !c.Allow("error: d") {
		t.Fatal("after window advance, match should be allowed again")
	}
}

func TestAllow_Disabled_AlwaysTrue(t *testing.T) {
	c, _ := New("", 0, time.Minute)
	for i := 0; i < 20; i++ {
		if !c.Allow("error: flood") {
			t.Fatalf("disabled clamp should always allow (iteration %d)", i)
		}
	}
}
