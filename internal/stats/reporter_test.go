package stats_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"greplive/internal/stats"
)

func TestReporter_Print_ContainsFields(t *testing.T) {
	c := stats.New()
	c.IncrRead()
	c.IncrRead()
	c.IncrMatched()

	var buf bytes.Buffer
	r := stats.NewReporter(c, time.Hour, &buf)
	r.Print()

	out := buf.String()
	for _, want := range []string{"[stats]", "read=", "matched=", "dropped=", "elapsed="} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output %q", want, out)
		}
	}
}

func TestReporter_Print_CorrectCounts(t *testing.T) {
	c := stats.New()
	for i := 0; i < 5; i++ {
		c.IncrRead()
	}
	for i := 0; i < 3; i++ {
		c.IncrMatched()
	}
	c.IncrDropped()

	var buf bytes.Buffer
	r := stats.NewReporter(c, time.Hour, &buf)
	r.Print()

	out := buf.String()
	if !strings.Contains(out, "read=5") {
		t.Errorf("expected read=5 in %q", out)
	}
	if !strings.Contains(out, "matched=3") {
		t.Errorf("expected matched=3 in %q", out)
	}
	if !strings.Contains(out, "dropped=1") {
		t.Errorf("expected dropped=1 in %q", out)
	}
}

func TestReporter_NilWriterDefaultsToStderr(t *testing.T) {
	c := stats.New()
	// Should not panic when w is nil.
	r := stats.NewReporter(c, time.Hour, nil)
	if r == nil {
		t.Fatal("expected non-nil reporter")
	}
}

func TestReporter_StartStop(t *testing.T) {
	c := stats.New()
	var buf bytes.Buffer
	r := stats.NewReporter(c, 20*time.Millisecond, &buf)
	r.Start()
	time.Sleep(55 * time.Millisecond)
	r.Stop()

	lines := strings.Count(buf.String(), "[stats]")
	// Expect at least 2 ticks in 55ms with a 20ms interval.
	if lines < 2 {
		t.Fatalf("expected at least 2 stat lines, got %d", lines)
	}
}
