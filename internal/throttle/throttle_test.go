package throttle_test

import (
	"testing"
	"time"

	"greplive/internal/throttle"
)

func TestEnabled_ZeroMaxLines(t *testing.T) {
	th := throttle.New(0, time.Second)
	if th.Enabled() {
		t.Fatal("expected throttle to be disabled when maxLines=0")
	}
}

func TestEnabled_PositiveMaxLines(t *testing.T) {
	th := throttle.New(5, time.Second)
	if !th.Enabled() {
		t.Fatal("expected throttle to be enabled")
	}
}

func TestAllow_ZeroRate_AlwaysAllows(t *testing.T) {
	th := throttle.New(0, time.Second)
	for i := 0; i < 1000; i++ {
		if !th.Allow() {
			t.Fatal("disabled throttle should always allow")
		}
	}
}

func TestAllow_WithinLimit(t *testing.T) {
	th := throttle.New(5, time.Second)
	for i := 0; i < 5; i++ {
		if !th.Allow() {
			t.Fatalf("line %d should be allowed (within limit)", i)
		}
	}
}

func TestAllow_ExceedsLimit_Drops(t *testing.T) {
	th := throttle.New(3, time.Second)
	allowed := 0
	for i := 0; i < 10; i++ {
		if th.Allow() {
			allowed++
		}
	}
	if allowed != 3 {
		t.Fatalf("expected 3 allowed, got %d", allowed)
	}
}

func TestAllow_AfterWindowExpires_Resets(t *testing.T) {
	th := throttle.New(2, 50*time.Millisecond)
	if !th.Allow() {
		t.Fatal("first allow should succeed")
	}
	if !th.Allow() {
		t.Fatal("second allow should succeed")
	}
	if th.Allow() {
		t.Fatal("third allow should be dropped")
	}
	time.Sleep(60 * time.Millisecond)
	if !th.Allow() {
		t.Fatal("allow after window expiry should succeed")
	}
}

func TestReset_ClearsCounter(t *testing.T) {
	th := throttle.New(2, time.Second)
	th.Allow()
	th.Allow()
	th.Reset()
	if !th.Allow() {
		t.Fatal("allow after Reset should succeed")
	}
}
