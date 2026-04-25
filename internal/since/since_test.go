package since_test

import (
	"testing"
	"time"

	"github.com/user/greplive/internal/since"
)

func mustParse(layout, value string) time.Time {
	t, err := time.Parse(layout, value)
	if err != nil {
		panic(err)
	}
	return t
}

func TestNew_Disabled_WhenZeroCutoff(t *testing.T) {
	f := since.New(time.Time{})
	if f.Enabled() {
		t.Fatal("expected filter to be disabled for zero cutoff")
	}
}

func TestNew_Enabled_WhenCutoffSet(t *testing.T) {
	f := since.New(time.Now())
	if !f.Enabled() {
		t.Fatal("expected filter to be enabled")
	}
}

func TestAllow_Disabled_AlwaysTrue(t *testing.T) {
	f := since.New(time.Time{})
	for _, line := range []string{
		"2020-01-01T00:00:00Z some old log",
		"no timestamp here",
		"",
	} {
		if !f.Allow(line) {
			t.Errorf("disabled filter should allow %q", line)
		}
	}
}

func TestAllow_NoTimestamp_AlwaysTrue(t *testing.T) {
	cutoff := mustParse(time.RFC3339, "2024-01-01T00:00:00Z")
	f := since.New(cutoff)
	if !f.Allow("plain log line without any timestamp") {
		t.Fatal("line without timestamp should always be allowed")
	}
}

func TestAllow_TimestampAfterCutoff_Allowed(t *testing.T) {
	cutoff := mustParse(time.RFC3339, "2024-01-01T00:00:00Z")
	f := since.New(cutoff)
	line := "2024-06-15T12:00:00Z INFO server started"
	if !f.Allow(line) {
		t.Errorf("line after cutoff should be allowed: %q", line)
	}
}

func TestAllow_TimestampEqualCutoff_Allowed(t *testing.T) {
	cutoff := mustParse(time.RFC3339, "2024-01-01T00:00:00Z")
	f := since.New(cutoff)
	line := "2024-01-01T00:00:00Z INFO exactly at cutoff"
	if !f.Allow(line) {
		t.Errorf("line equal to cutoff should be allowed: %q", line)
	}
}

func TestAllow_TimestampBeforeCutoff_Dropped(t *testing.T) {
	cutoff := mustParse(time.RFC3339, "2024-01-01T00:00:00Z")
	f := since.New(cutoff)
	line := "2023-12-31T23:59:59Z WARN old event"
	if f.Allow(line) {
		t.Errorf("line before cutoff should be dropped: %q", line)
	}
}

func TestAllow_CustomLayout(t *testing.T) {
	layout := "2006-01-02 15:04:05"
	cutoff := mustParse(layout, "2024-03-01 00:00:00")
	f := since.New(cutoff, layout)

	allowed := "2024-03-15 08:30:00 INFO request received"
	if !f.Allow(allowed) {
		t.Errorf("expected line to be allowed: %q", allowed)
	}

	dropped := "2024-02-28 23:59:59 DEBUG old debug"
	if f.Allow(dropped) {
		t.Errorf("expected line to be dropped: %q", dropped)
	}
}
