package timestamp_test

import (
	"strings"
	"testing"
	"time"

	"greplive/internal/timestamp"
)

var fixed = time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

func fixedClock() time.Time { return fixed }

func TestNew_Disabled_WhenEmptyFormat(t *testing.T) {
	s := timestamp.New("")
	if s.Enabled() {
		t.Fatal("expected disabled for empty format")
	}
}

func TestNew_Enabled_WhenFormatProvided(t *testing.T) {
	s := timestamp.New(time.RFC3339)
	if !s.Enabled() {
		t.Fatal("expected enabled")
	}
}

func TestApply_Disabled_Unchanged(t *testing.T) {
	s := timestamp.New("")
	got := s.Apply("hello")
	if got != "hello" {
		t.Fatalf("expected unchanged, got %q", got)
	}
}

func TestApply_PrependsTimestamp(t *testing.T) {
	s := timestamp.New(time.RFC3339, timestamp.WithClock(fixedClock))
	got := s.Apply("hello")
	want := fixed.Format(time.RFC3339) + " hello"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestApply_CustomFormat(t *testing.T) {
	format := "2006-01-02"
	s := timestamp.New(format, timestamp.WithClock(fixedClock))
	got := s.Apply("msg")
	if !strings.HasPrefix(got, "2024-01-15 ") {
		t.Fatalf("unexpected prefix in %q", got)
	}
}

func TestApplyAll_MultipleStampers(t *testing.T) {
	// Only the first enabled stamper should fire; second is disabled.
	s1 := timestamp.New("15:04:05", timestamp.WithClock(fixedClock))
	s2 := timestamp.New("")
	got := timestamp.ApplyAll("line", []*timestamp.Stamper{s1, s2})
	if !strings.HasPrefix(got, "12:00:00 ") {
		t.Fatalf("unexpected result %q", got)
	}
	if strings.Count(got, "12:00:00") != 1 {
		t.Fatalf("timestamp applied more than once: %q", got)
	}
}

func TestApplyAll_NoStampers_Unchanged(t *testing.T) {
	got := timestamp.ApplyAll("line", nil)
	if got != "line" {
		t.Fatalf("expected unchanged, got %q", got)
	}
}
