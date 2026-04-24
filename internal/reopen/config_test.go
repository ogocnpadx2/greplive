package reopen_test

import (
	"testing"

	"github.com/user/greplive/internal/reopen"
)

func TestDefaultConfig_AutoWatch(t *testing.T) {
	cfg := reopen.DefaultConfig()
	if cfg.AutoWatch {
		t.Fatal("default AutoWatch should be false")
	}
}

func TestConfig_Build_EmptyPath_ReturnsError(t *testing.T) {
	cfg := reopen.Config{Path: ""}
	_, err := cfg.Build()
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestConfig_Build_NonExistentPath_ReturnsError(t *testing.T) {
	cfg := reopen.Config{Path: "/no/such/file.log"}
	_, err := cfg.Build()
	if err == nil {
		t.Fatal("expected error for non-existent file")
	}
}

func TestConfig_Build_ValidPath_ReturnsReader(t *testing.T) {
	p := writeTmp(t, "hello")
	cfg := reopen.Config{Path: p}
	r, err := cfg.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r == nil {
		t.Fatal("expected non-nil Reader")
	}
	r.Close()
}
