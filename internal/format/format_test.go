package format_test

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"greplive/internal/format"
)

var fixedTime = time.Date(2024, 1, 2, 15, 4, 5, 0, time.UTC)

func TestNew_DefaultsToPlain(t *testing.T) {
	f := format.New("unknown")
	if f == nil {
		t.Fatal("expected non-nil formatter")
	}
}

func TestPlain_MessageOnly(t *testing.T) {
	f := format.New(format.ModePlain)
	out := f.Format(format.Entry{Message: "hello world"})
	if out != "hello world" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestPlain_WithLevelAndTimestamp(t *testing.T) {
	f := format.New(format.ModePlain)
	out := f.Format(format.Entry{
		Timestamp: fixedTime,
		Level:     "error",
		Message:   "boom",
	})
	if !strings.Contains(out, "[ERROR]") {
		t.Errorf("expected [ERROR] in %q", out)
	}
	if !strings.Contains(out, "2024-01-02T15:04:05Z") {
		t.Errorf("expected timestamp in %q", out)
	}
	if !strings.Contains(out, "boom") {
		t.Errorf("expected message in %q", out)
	}
}

func TestPlain_WithFields(t *testing.T) {
	f := format.New(format.ModePlain)
	out := f.Format(format.Entry{
		Message: "msg",
		Fields:  map[string]string{"host": "srv1"},
	})
	if !strings.Contains(out, "host=srv1") {
		t.Errorf("expected field in %q", out)
	}
}

func TestJSON_ValidOutput(t *testing.T) {
	f := format.New(format.ModeJSON)
	out := f.Format(format.Entry{
		Timestamp: fixedTime,
		Level:     "warn",
		Message:   "disk full",
		Fields:    map[string]string{"disk": "/dev/sda"},
	})
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(out), &m); err != nil {
		t.Fatalf("invalid JSON: %v — output: %s", err, out)
	}
	if m["msg"] != "disk full" {
		t.Errorf("unexpected msg: %v", m["msg"])
	}
	if m["level"] != "warn" {
		t.Errorf("unexpected level: %v", m["level"])
	}
	if m["disk"] != "/dev/sda" {
		t.Errorf("unexpected field: %v", m["disk"])
	}
}

func TestJSON_NoTimestamp(t *testing.T) {
	f := format.New(format.ModeJSON)
	out := f.Format(format.Entry{Message: "no time"})
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(out), &m); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if _, ok := m["ts"]; ok {
		t.Error("expected no ts field when timestamp is zero")
	}
}
