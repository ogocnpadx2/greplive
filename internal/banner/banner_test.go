package banner

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrint_ContainsPattern(t *testing.T) {
	var buf bytes.Buffer
	p := New(&buf)
	p.Print(Config{Pattern: `ERROR.*`})
	if !strings.Contains(buf.String(), `ERROR.*`) {
		t.Errorf("expected pattern in banner, got:\n%s", buf.String())
	}
}

func TestPrint_DefaultPattern(t *testing.T) {
	var buf bytes.Buffer
	p := New(&buf)
	p.Print(Config{})
	if !strings.Contains(buf.String(), "(none)") {
		t.Errorf("expected '(none)' when no pattern set")
	}
}

func TestPrint_DefaultLevel(t *testing.T) {
	var buf bytes.Buffer
	p := New(&buf)
	p.Print(Config{})
	if !strings.Contains(buf.String(), "all") {
		t.Errorf("expected level 'all' when no level set")
	}
}

func TestPrint_DefaultSourceStdin(t *testing.T) {
	var buf bytes.Buffer
	p := New(&buf)
	p.Print(Config{})
	if !strings.Contains(buf.String(), "stdin") {
		t.Errorf("expected source 'stdin' when no file set")
	}
}

func TestPrint_InputFile(t *testing.T) {
	var buf bytes.Buffer
	p := New(&buf)
	p.Print(Config{InputFile: "/var/log/app.log"})
	if !strings.Contains(buf.String(), "/var/log/app.log") {
		t.Errorf("expected file path in banner")
	}
}

func TestPrint_RateLimit(t *testing.T) {
	var buf bytes.Buffer
	p := New(&buf)
	p.Print(Config{RateLimit: 100})
	if !strings.Contains(buf.String(), "100") {
		t.Errorf("expected rate limit in banner")
	}
}

func TestPrint_NilWriterDefaultsToStderr(t *testing.T) {
	// Should not panic.
	p := New(nil)
	p.Print(Config{})
}

func TestPrint_DedupeAndTruncate(t *testing.T) {
	var buf bytes.Buffer
	p := New(&buf)
	p.Print(Config{Dedupe: true, Truncate: 120})
	out := buf.String()
	if !strings.Contains(out, "dedupe") {
		t.Errorf("expected dedupe in banner")
	}
	if !strings.Contains(out, "120") {
		t.Errorf("expected truncate width in banner")
	}
}
