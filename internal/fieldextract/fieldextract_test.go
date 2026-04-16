package fieldextract_test

import (
	"testing"

	"greplive/internal/fieldextract"
)

func TestNew_ValidPattern(t *testing.T) {
	e, err := fieldextract.New(`(?P<level>\w+)\s+(?P<msg>.+)`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	fields := e.Fields()
	if len(fields) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(fields))
	}
}

func TestNew_InvalidPattern(t *testing.T) {
	_, err := fieldextract.New(`(?P<bad`)
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestExtract_Match(t *testing.T) {
	e, _ := fieldextract.New(`level=(?P<level>\w+)\s+msg=(?P<msg>[^\s]+)`)
	result := e.Extract("level=ERROR msg=disk_full")
	if result["level"] != "ERROR" {
		t.Errorf("expected ERROR, got %q", result["level"])
	}
	if result["msg"] != "disk_full" {
		t.Errorf("expected disk_full, got %q", result["msg"])
	}
}

func TestExtract_NoMatch(t *testing.T) {
	e, _ := fieldextract.New(`(?P<level>\d+)`)
	result := e.Extract("no digits here")
	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}

func TestExtractKV_Basic(t *testing.T) {
	result := fieldextract.ExtractKV(`ts=2024-01-01 level=INFO msg="server started" port=8080`)
	if result["level"] != "INFO" {
		t.Errorf("expected INFO, got %q", result["level"])
	}
	if result["port"] != "8080" {
		t.Errorf("expected 8080, got %q", result["port"])
	}
	if result["msg"] != "server started" {
		t.Errorf("expected 'server started', got %q", result["msg"])
	}
}

func TestExtractKV_Empty(t *testing.T) {
	result := fieldextract.ExtractKV("no key value pairs here")
	if len(result) != 0 {
		t.Errorf("expected empty map, got %v", result)
	}
}

func TestFields_ReturnsCopy(t *testing.T) {
	e, _ := fieldextract.New(`(?P<a>\w+)`)
	f1 := e.Fields()
	f1[0] = "mutated"
	f2 := e.Fields()
	if f2[0] == "mutated" {
		t.Error("Fields() should return a copy")
	}
}
