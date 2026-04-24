package skip_test

import (
	"testing"

	"greplive/internal/skip"
)

func TestNew_Disabled_WhenZero(t *testing.T) {
	s, err := skip.New(0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Enabled() {
		t.Error("expected Skipper to be disabled for n=0")
	}
}

func TestNew_Disabled_WhenNegative(t *testing.T) {
	s, err := skip.New(-5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Enabled() {
		t.Error("expected Skipper to be disabled for negative n")
	}
}

func TestNew_Enabled_WhenPositive(t *testing.T) {
	s, err := skip.New(3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !s.Enabled() {
		t.Error("expected Skipper to be enabled for n=3")
	}
}

func TestAllow_Disabled_AlwaysEmits(t *testing.T) {
	s, _ := skip.New(0)
	for i := 0; i < 5; i++ {
		if !s.Allow("line") {
			t.Errorf("disabled Skipper should always allow, failed on iteration %d", i)
		}
	}
}

func TestAllow_DropsFirstNLines(t *testing.T) {
	s, _ := skip.New(3)
	results := make([]bool, 6)
	for i := range results {
		results[i] = s.Allow("line")
	}
	expected := []bool{false, false, false, true, true, true}
	for i, got := range results {
		if got != expected[i] {
			t.Errorf("Allow()[%d] = %v, want %v", i, got, expected[i])
		}
	}
}

func TestRemaining_DecreasesWithAllowCalls(t *testing.T) {
	s, _ := skip.New(3)
	if r := s.Remaining(); r != 3 {
		t.Errorf("expected Remaining()=3, got %d", r)
	}
	s.Allow("line")
	if r := s.Remaining(); r != 2 {
		t.Errorf("expected Remaining()=2, got %d", r)
	}
	s.Allow("line")
	s.Allow("line")
	if r := s.Remaining(); r != 0 {
		t.Errorf("expected Remaining()=0 after exhaustion, got %d", r)
	}
}

func TestReset_RestartsDropping(t *testing.T) {
	s, _ := skip.New(2)
	s.Allow("a")
	s.Allow("b")
	if !s.Allow("c") {
		t.Fatal("expected third line to pass through")
	}
	s.Reset()
	if s.Allow("d") {
		t.Error("expected first line after Reset to be dropped")
	}
}

func TestRemaining_Disabled_ReturnsZero(t *testing.T) {
	s, _ := skip.New(0)
	if r := s.Remaining(); r != 0 {
		t.Errorf("expected Remaining()=0 for disabled skipper, got %d", r)
	}
}
