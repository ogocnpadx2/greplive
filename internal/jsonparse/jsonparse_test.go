package jsonparse

import (
	"strings"
	"testing"
)

func TestIsJSON_ValidObject(t *testing.T) {
	if !IsJSON(`{"level":"info","message":"hello"}`) {
		t.Fatal("expected true for valid JSON object")
	}
}

func TestIsJSON_PlainText(t *testing.T) {
	if IsJSON("plain log line") {
		t.Fatal("expected false for plain text")
	}
}

func TestIsJSON_EmptyString(t *testing.T) {
	if IsJSON("") {
		t.Fatal("expected false for empty string")
	}
}

func TestIsJSON_JSONArray(t *testing.T) {
	if IsJSON(`["a","b"]`) {
		t.Fatal("expected false for JSON array (not object)")
	}
}

func TestFormat_NonJSON_Unchanged(t *testing.T) {
	f := New()
	line := "not json"
	if got := f.Format(line); got != line {
		t.Fatalf("expected %q, got %q", line, got)
	}
}

func TestFormat_Pretty(t *testing.T) {
	f := New(WithPretty(true))
	line := `{"a":1,"b":2}`
	got := f.Format(line)
	if !strings.Contains(got, "\n") {
		t.Fatal("expected indented output to contain newlines")
	}
}

func TestFormat_NoPretty(t *testing.T) {
	f := New(WithPretty(false))
	line := `{"a":1}`
	if got := f.Format(line); got != line {
		t.Fatalf("expected unchanged line, got %q", got)
	}
}

func TestExtractMessage_Found(t *testing.T) {
	f := New(WithMessageKey("msg"))
	line := `{"level":"info","msg":"hello world"}`
	if got := f.ExtractMessage(line); got != "hello world" {
		t.Fatalf("expected 'hello world', got %q", got)
	}
}

func TestExtractMessage_Missing(t *testing.T) {
	f := New()
	line := `{"level":"info"}`
	if got := f.ExtractMessage(line); got != line {
		t.Fatalf("expected original line, got %q", got)
	}
}

func TestExtractMessage_NonJSON(t *testing.T) {
	f := New()
	line := "plain text"
	if got := f.ExtractMessage(line); got != line {
		t.Fatalf("expected original line, got %q", got)
	}
}

func TestNew_DefaultMessageKey(t *testing.T) {
	f := New()
	line := `{"message":"hi"}`
	if got := f.ExtractMessage(line); got != "hi" {
		t.Fatalf("expected 'hi', got %q", got)
	}
}
