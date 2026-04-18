package scope

import (
	"testing"
)

func TestNew_BothEmpty_Disabled(t *testing.T) {
	sc, err := New("", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sc.Enabled() {
		t.Fatal("expected disabled scope")
	}
}

func TestNew_InvalidStart_ReturnsError(t *testing.T) {
	_, err := New("[", "end")
	if err == nil {
		t.Fatal("expected error for invalid start pattern")
	}
}

func TestNew_InvalidEnd_ReturnsError(t *testing.T) {
	_, err := New("start", "[")
	if err == nil {
		t.Fatal("expected error for invalid end pattern")
	}
}

func TestAllow_Disabled_PassesAll(t *testing.T) {
	sc, _ := New("", "")
	for _, line := range []string{"a", "b", "c"} {
		if !sc.Allow(line) {
			t.Fatalf("expected line %q to pass", line)
		}
	}
}

func TestAllow_StartOnly_EmitsFromMatch(t *testing.T) {
	sc, _ := New("BEGIN", "")
	lines := []string{"before", "BEGIN", "inside", "more"}
	want := []bool{false, true, true, true}
	for i, l := range lines {
		if got := sc.Allow(l); got != want[i] {
			t.Errorf("line %d %q: got %v want %v", i, l, got, want[i])
		}
	}
}

func TestAllow_StartAndEnd_EmitsRange(t *testing.T) {
	sc, _ := New("START", "STOP")
	cases := []struct {
		line string
		want bool
	}{
		{"noise", false},
		{"START", true},
		{"middle", true},
		{"STOP", true},
		{"after", false},
	}
	for _, c := range cases {
		if got := sc.Allow(c.line); got != c.want {
			t.Errorf("line %q: got %v want %v", c.line, got, c.want)
		}
	}
}

func TestReset_RestartsScope(t *testing.T) {
	sc, _ := New("START", "STOP")
	sc.Allow("START")
	sc.Allow("STOP")
	sc.Reset()
	if sc.Allow("middle") {
		t.Fatal("expected line to be blocked after reset")
	}
}
