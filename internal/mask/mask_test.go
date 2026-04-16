package mask

import (
	"testing"
)

func TestNew_ValidPattern(t *testing.T) {
	m, err := New("password", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m == nil {
		t.Fatal("expected non-nil Masker")
	}
}

func TestNew_EmptyPattern(t *testing.T) {
	_, err := New("", "")
	if err == nil {
		t.Fatal("expected error for empty pattern")
	}
}

func TestNew_InvalidPattern(t *testing.T) {
	_, err := New("[invalid", "")
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestApply_KeyEquals(t *testing.T) {
	m, _ := New("password", "***")
	got := m.Apply(`user=alice password=secret123 host=db`)
	want := `user=alice password=*** host=db`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestApply_JSONField(t *testing.T) {
	m, _ := New("token", "REDACTED")
	got := m.Apply(`{"user":"bob","token":"abc.def.ghi"}`)
	if got == `{"user":"bob","token":"abc.def.ghi"}` {
		t.Error("expected token value to be masked")
	}
}

func TestApply_NoMatch(t *testing.T) {
	m, _ := New("secret", "")
	line := `user=alice host=db`
	got := m.Apply(line)
	if got != line {
		t.Errorf("expected unchanged line, got %q", got)
	}
}

func TestApply_CaseInsensitive(t *testing.T) {
	m, _ := New("password", "")
	got := m.Apply(`PASSWORD=hunter2`)
	if got == `PASSWORD=hunter2` {
		t.Error("expected case-insensitive match to be masked")
	}
}

func TestApplyAll_MultipleFields(t *testing.T) {
	m1, _ := New("password", "")
	m2, _ := New("token", "")
	line := `password=abc token=xyz user=carol`
	got := ApplyAll([]*Masker{m1, m2}, line)
	if got == line {
		t.Error("expected both fields to be masked")
	}
	if contains(got, "abc") || contains(got, "xyz") {
		t.Errorf("sensitive values still present: %q", got)
	}
}

func TestNew_CustomPlaceholder(t *testing.T) {
	m, _ := New("secret", "[hidden]")
	got := m.Apply(`secret=mysecret`)
	if !contains(got, "[hidden]") {
		t.Errorf("expected custom placeholder in output, got %q", got)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		(func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		})())
}
