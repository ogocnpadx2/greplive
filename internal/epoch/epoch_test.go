package epoch

import (
	"strings"
	"testing"
	"time"
)

func TestNew_Disabled_WhenEmptyLayout(t *testing.T) {
	c := New("", false)
	if c.Enabled() {
		t.Fatal("expected disabled")
	}
}

func TestNew_Enabled_WhenLayoutProvided(t *testing.T) {
	c := New(time.RFC3339, false)
	if !c.Enabled() {
		t.Fatal("expected enabled")
	}
}

func TestApply_Disabled_Unchanged(t *testing.T) {
	c := New("", false)
	line := "ts=1700000000 msg=hello"
	if got := c.Apply(line); got != line {
		t.Fatalf("want %q, got %q", line, got)
	}
}

func TestApply_ReplacesSecondEpoch(t *testing.T) {
	c := New("2006-01-02", true)
	ts := time.Unix(1700000000, 0).UTC().Format("2006-01-02")
	line := "ts=1700000000 msg=hello"
	got := c.Apply(line)
	if !strings.Contains(got, ts) {
		t.Fatalf("expected %q in %q", ts, got)
	}
	if strings.Contains(got, "1700000000") {
		t.Fatalf("epoch should have been replaced in %q", got)
	}
}

func TestApply_ReplacesMilliEpoch(t *testing.T) {
	c := New("2006-01-02", true)
	ts := time.UnixMilli(1700000000000).UTC().Format("2006-01-02")
	line := "ts=1700000000000 msg=hello"
	got := c.Apply(line)
	if !strings.Contains(got, ts) {
		t.Fatalf("expected %q in %q", ts, got)
	}
}

func TestApply_NoEpoch_Unchanged(t *testing.T) {
	c := New(time.RFC3339, false)
	line := "no timestamps here"
	if got := c.Apply(line); got != line {
		t.Fatalf("want %q, got %q", line, got)
	}
}

func TestApplyAll_MultipleConverters(t *testing.T) {
	c1 := New("2006-01-02", true)
	c2 := New("", false) // disabled – should be a no-op
	line := "ts=1700000000 msg=ok"
	got := ApplyAll(line, []*Converter{c1, c2})
	if strings.Contains(got, "1700000000") {
		t.Fatalf("epoch should have been replaced in %q", got)
	}
}

func TestParseLayout_KnownNames(t *testing.T) {
	cases := []struct {
		name string
		want string
	}{
		{"rfc3339", time.RFC3339},
		{"rfc3339ms", "2006-01-02T15:04:05.000Z07:00"},
		{"short", "2006-01-02 15:04:05"},
		{"date", "2006-01-02"},
		{"custom", "custom"},
	}
	for _, tc := range cases {
		if got := ParseLayout(tc.name); got != tc.want {
			t.Errorf("ParseLayout(%q) = %q, want %q", tc.name, got, tc.want)
		}
	}
}
