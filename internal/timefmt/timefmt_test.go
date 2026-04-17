package timefmt_test

import (
	"testing"
	"time"

	"greplive/internal/timefmt"
)

var epoch = time.Date(2024, 6, 1, 12, 30, 45, 123000000, time.UTC)

func TestFormat_RFC3339(t *testing.T) {
	got := timefmt.RFC3339.Format(epoch)
	want := "2024-06-01T12:30:45Z"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestFormat_RFC3339Ms(t *testing.T) {
	got := timefmt.RFC3339Ms.Format(epoch)
	want := "2024-06-01T12:30:45.123Z"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestFormat_Short(t *testing.T) {
	got := timefmt.Short.Format(epoch)
	if got != "12:30:45" {
		t.Fatalf("unexpected %q", got)
	}
}

func TestFormat_ShortMs(t *testing.T) {
	got := timefmt.ShortMs.Format(epoch)
	if got != "12:30:45.123" {
		t.Fatalf("unexpected %q", got)
	}
}

func TestParse_Known(t *testing.T) {
	for _, name := range timefmt.Names() {
		f, ok := timefmt.Parse(name)
		if !ok {
			t.Errorf("expected %q to be found", name)
		}
		if f.Name() != name {
			t.Errorf("name mismatch: got %q want %q", f.Name(), name)
		}
	}
}

func TestParse_Unknown(t *testing.T) {
	f, ok := timefmt.Parse("bogus")
	if ok {
		t.Fatal("expected ok=false")
	}
	if f.Name() != "rfc3339" {
		t.Fatalf("expected default rfc3339, got %q", f.Name())
	}
}

func TestNames_NonEmpty(t *testing.T) {
	if len(timefmt.Names()) == 0 {
		t.Fatal("expected at least one format")
	}
}

func TestFormat_Layout(t *testing.T) {
	if timefmt.RFC3339.Layout() == "" {
		t.Fatal("layout should not be empty")
	}
}
