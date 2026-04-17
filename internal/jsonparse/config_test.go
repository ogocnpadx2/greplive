package jsonparse

import "testing"

func TestDefaultConfig_MessageKey(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.MessageKey != "message" {
		t.Fatalf("expected 'message', got %q", cfg.MessageKey)
	}
}

func TestDefaultConfig_Pretty(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Pretty {
		t.Fatal("expected Pretty to be false by default")
	}
}

func TestConfig_Build_CustomKey(t *testing.T) {
	cfg := Config{MessageKey: "msg"}
	f := cfg.Build()
	line := `{"msg":"hello"}`
	if got := f.ExtractMessage(line); got != "hello" {
		t.Fatalf("expected 'hello', got %q", got)
	}
}

func TestConfig_Build_EmptyKey_UsesDefault(t *testing.T) {
	cfg := Config{}
	f := cfg.Build()
	line := `{"message":"default"}`
	if got := f.ExtractMessage(line); got != "default" {
		t.Fatalf("expected 'default', got %q", got)
	}
}

func TestConfig_Build_Pretty(t *testing.T) {
	cfg := Config{Pretty: true}
	f := cfg.Build()
	line := `{"a":1}`
	got := f.Format(line)
	if got == line {
		t.Fatal("expected pretty-printed output to differ from compact input")
	}
}
