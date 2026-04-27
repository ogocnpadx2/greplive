package until_test

import (
	"testing"

	"github.com/user/greplive/internal/until"
)

func TestNew_EmptyPattern_Disabled(t *testing.T) {
	u, err := until.New("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u.Enabled() {
		t.Fatal("expected disabled for empty pattern")
	}
}

func TestNew_InvalidPattern_ReturnsError(t *testing.T) {
	_, err := until.New("[invalid")
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestNew_ValidPattern_Enabled(t *testing.T) {
	u, err := until.New(`ERROR`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !u.Enabled() {
		t.Fatal("expected enabled for non-empty pattern")
	}
}

func TestAllow_Disabled_AlwaysTrue(t *testing.T) {
	u, _ := until.New("")
	for _, line := range []string{"foo", "bar", "ERROR: boom"} {
		if !u.Allow(line) {
			t.Errorf("disabled Until should always allow, got false for %q", line)
		}
	}
}

func TestAllow_LinesBeforeTrigger_Emitted(t *testing.T) {
	u, _ := until.New(`STOP`)
	lines := []string{"alpha", "beta", "gamma"}
	for _, l := range lines {
		if !u.Allow(l) {
			t.Errorf("expected line %q to be allowed before trigger", l)
		}
	}
}

func TestAllow_TriggerLine_NotEmitted(t *testing.T) {
	u, _ := until.New(`STOP`)
	if !u.Allow("before") {
		t.Fatal("expected before-line to be allowed")
	}
	if u.Allow("STOP here") {
		t.Fatal("trigger line should not be emitted")
	}
}

func TestAllow_LinesAfterTrigger_Dropped(t *testing.T) {
	u, _ := until.New(`STOP`)
	u.Allow("before")
	u.Allow("STOP now")
	for _, l := range []string{"after1", "after2", "after3"} {
		if u.Allow(l) {
			t.Errorf("expected line %q to be dropped after trigger", l)
		}
	}
}

func TestReset_AllowsLinesAgain(t *testing.T) {
	u, _ := until.New(`STOP`)
	u.Allow("line")
	u.Allow("STOP")
	if u.Allow("after") {
		t.Fatal("expected drop after trigger before reset")
	}
	u.Reset()
	if !u.Allow("after reset") {
		t.Fatal("expected allow after reset")
	}
}
