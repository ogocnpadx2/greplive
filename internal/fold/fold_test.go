package fold_test

import (
	"testing"

	"greplive/internal/fold"
)

func TestNew_Disabled_WhenNotCaseInsensitive(t *testing.T) {
	f, err := fold.New("error", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Enabled() {
		t.Fatal("expected folder to be disabled")
	}
}

func TestNew_Disabled_WhenEmptyPattern(t *testing.T) {
	f, err := fold.New("", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Enabled() {
		t.Fatal("expected folder to be disabled for empty pattern")
	}
}

func TestNew_InvalidPattern_ReturnsError(t *testing.T) {
	_, err := fold.New("[invalid", true)
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestNew_ValidPattern_Enabled(t *testing.T) {
	f, err := fold.New("error", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !f.Enabled() {
		t.Fatal("expected folder to be enabled")
	}
}

func TestMatch_CaseInsensitive(t *testing.T) {
	f, _ := fold.New("error", true)
	cases := []struct {
		line string
		want bool
	}{
		{"an ERROR occurred", true},
		{"Error: something", true},
		{"error", true},
		{"no match here", false},
	}
	for _, tc := range cases {
		got := f.Match(tc.line)
		if got != tc.want {
			t.Errorf("Match(%q) = %v, want %v", tc.line, got, tc.want)
		}
	}
}

func TestMatch_Disabled_AlwaysTrue(t *testing.T) {
	f, _ := fold.New("error", false)
	if !f.Match("completely unrelated line") {
		t.Fatal("disabled folder should always return true")
	}
}

func TestFold_LowersCase(t *testing.T) {
	if got := fold.Fold("HELLO World"); got != "hello world" {
		t.Errorf("Fold() = %q, want %q", got, "hello world")
	}
}
