package cli

import (
	"flag"
	"testing"
	"time"
)

func newFS() *flag.FlagSet {
	return flag.NewFlagSet("test", flag.ContinueOnError)
}

func TestParseFlags_Defaults(t *testing.T) {
	f, err := ParseFlags(newFS(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if f.Pattern != "" || f.Level != "any" || f.Follow || f.Timestamp {
		t.Fatalf("unexpected defaults: %+v", f)
	}
	if f.StatsInterval != 0 {
		t.Fatalf("expected zero stats interval, got %v", f.StatsInterval)
	}
}

func TestParseFlags_AllFlags(t *testing.T) {
	args := []string{
		"-pattern", "ERROR",
		"-level", "warn",
		"-follow",
		"-timestamp",
		"-stats", "5s",
		"-json",
		"-max-rate", "100",
		"-dedupe",
		"-max-len", "200",
		"-before", "2",
		"-after", "3",
	}
	f, err := ParseFlags(newFS(), args)
	if err != nil {
		t.Fatal(err)
	}
	if f.Pattern != "ERROR" { t.Error("pattern") }
	if f.Level != "warn" { t.Error("level") }
	if !f.Follow { t.Error("follow") }
	if !f.Timestamp { t.Error("timestamp") }
	if f.StatsInterval != 5*time.Second { t.Error("stats") }
	if !f.JSON { t.Error("json") }
	if f.MaxRate != 100 { t.Error("max-rate") }
	if !f.Dedupe { t.Error("dedupe") }
	if f.MaxLen != 200 { t.Error("max-len") }
	if f.Before != 2 { t.Error("before") }
	if f.After != 3 { t.Error("after") }
}

func TestParseFlags_InvalidLevel(t *testing.T) {
	_, err := ParseFlags(newFS(), []string{"-level", "verbose"})
	if err == nil {
		t.Fatal("expected error for invalid level")
	}
}

func TestParseFlags_InvalidFlag(t *testing.T) {
	_, err := ParseFlags(newFS(), []string{"-unknown"})
	if err == nil {
		t.Fatal("expected error for unknown flag")
	}
}

func TestParseFlags_InvalidStatsInterval(t *testing.T) {
	_, err := ParseFlags(newFS(), []string{"-stats", "notaduration"})
	if err == nil {
		t.Fatal("expected error for invalid stats interval")
	}
}
