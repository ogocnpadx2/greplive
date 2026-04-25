package strip

import "testing"

func TestDefaultConfig_ANSI_True(t *testing.T) {
	cfg := DefaultConfig()
	if !cfg.ANSI {
		t.Fatal("expected ANSI=true in default config")
	}
}

func TestDefaultConfig_Control_False(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Control {
		t.Fatal("expected Control=false in default config")
	}
}

func TestConfig_Build_ANSIOnly(t *testing.T) {
	s := Config{ANSI: true, Control: false}.Build()
	if !s.Enabled() {
		t.Fatal("expected stripper to be enabled")
	}
	got := s.Apply("\x1b[31mwarn\x1b[0m")
	if got != "warn" {
		t.Fatalf("expected %q, got %q", "warn", got)
	}
}

func TestConfig_Build_ControlOnly(t *testing.T) {
	s := Config{ANSI: false, Control: true}.Build()
	if !s.Enabled() {
		t.Fatal("expected stripper to be enabled")
	}
	// ANSI escape should survive; control char should be removed.
	got := s.Apply("\x1b[31m\x01ok")
	want := "\x1b[31mok"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestConfig_Build_BothDisabled(t *testing.T) {
	s := Config{ANSI: false, Control: false}.Build()
	if s.Enabled() {
		t.Fatal("expected stripper to be disabled when both false")
	}
}

func TestDefaultConfig_Build_RemovesANSI(t *testing.T) {
	s := DefaultConfig().Build()
	got := s.Apply("\x1b[32mINFO\x1b[0m message")
	want := "INFO message"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}
