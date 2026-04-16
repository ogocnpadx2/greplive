package cli

import (
	"testing"
	"time"

	"github.com/user/greplive/internal/severity"
)

func TestParseFlags_Defaults(t *testing.T) {
	cfg, err := ParseFlags([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Pattern != "" {
		t.Errorf("expected empty pattern, got %q", cfg.Pattern)
	}
	if cfg.MinLevel != severity.Debug {
		t.Errorf("expected debug level, got %v", cfg.MinLevel)
	}
	if cfg.ShowTimestamp {
		t.Error("expected ShowTimestamp false")
	}
	if cfg.StatsInterval != 5*time.Second {
		t.Errorf("expected 5s stats interval, got %v", cfg.StatsInterval)
	}
}

func TestParseFlags_AllFlags(t *testing.T) {
	cfg, err := ParseFlags([]string{
		"-pattern", "ERROR",
		"-level", "warn",
		"-file", "/var/log/app.log",
		"-timestamp",
		"-level-prefix",
		"-stats",
		"-stats-interval", "10s",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Pattern != "ERROR" {
		t.Errorf("expected pattern ERROR, got %q", cfg.Pattern)
	}
	if cfg.MinLevel != severity.Warn {
		t.Errorf("expected warn level, got %v", cfg.MinLevel)
	}
	if cfg.File != "/var/log/app.log" {
		t.Errorf("expected file /var/log/app.log, got %q", cfg.File)
	}
	if !cfg.ShowTimestamp {
		t.Error("expected ShowTimestamp true")
	}
	if !cfg.ShowLevel {
		t.Error("expected ShowLevel true")
	}
	if !cfg.ShowStats {
		t.Error("expected ShowStats true")
	}
	if cfg.StatsInterval != 10*time.Second {
		t.Errorf("expected 10s, got %v", cfg.StatsInterval)
	}
}

func TestParseFlags_InvalidLevel(t *testing.T) {
	_, err := ParseFlags([]string{"-level", "verbose"})
	if err == nil {
		t.Fatal("expected error for unknown level")
	}
}

func TestParseFlags_InvalidFlag(t *testing.T) {
	_, err := ParseFlags([]string{"-nonexistent"})
	if err == nil {
		t.Fatal("expected error for unknown flag")
	}
}

func TestParseFlags_InvalidStatsInterval(t *testing.T) {
	_, err := ParseFlags([]string{"-stats-interval", "notaduration"})
	if err == nil {
		t.Fatal("expected error for invalid stats interval")
	}
}
