package output_test

import (
	"bytes"
	"strings"
	"testing"

	"greplive/internal/output"
)

func TestWriteLine_Basic(t *testing.T) {
	var buf bytes.Buffer
	w := output.New(&buf)
	w.WriteLine("some info message")

	got := buf.String()
	if !strings.Contains(got, "some info message") {
		t.Errorf("expected output to contain original message, got: %q", got)
	}
	if !strings.HasSuffix(got, "\n") {
		t.Errorf("expected output to end with newline, got: %q", got)
	}
}

func TestWriteLine_WithLevel(t *testing.T) {
	var buf bytes.Buffer
	w := output.New(&buf, output.WithLevel(true))
	w.WriteLine("ERROR something went wrong")

	got := buf.String()
	if !strings.Contains(got, "ERROR") {
		t.Errorf("expected level label in output, got: %q", got)
	}
}

func TestWriteLine_WithTimestamp(t *testing.T) {
	var buf bytes.Buffer
	w := output.New(&buf, output.WithTimestamp(true))
	w.WriteLine("INFO starting up")

	got := buf.String()
	// Timestamp format [HH:MM:SS]
	if !strings.Contains(got, "[") || !strings.Contains(got, "]") {
		t.Errorf("expected timestamp brackets in output, got: %q", got)
	}
}

func TestWriteLine_NilWriterDefaultsToStdout(t *testing.T) {
	// Should not panic when nil is passed; falls back to os.Stdout
	w := output.New(nil)
	// We can't capture stdout easily here, just ensure no panic
	_ = w
}

func TestWriteLine_MultipleLines(t *testing.T) {
	var buf bytes.Buffer
	w := output.New(&buf, output.WithLevel(true), output.WithTimestamp(true))

	lines := []string{
		"DEBUG initializing",
		"WARN disk usage high",
		"ERROR disk full",
	}
	for _, l := range lines {
		w.WriteLine(l)
	}

	got := buf.String()
	parts := strings.Split(strings.TrimRight(got, "\n"), "\n")
	if len(parts) != 3 {
		t.Errorf("expected 3 output lines, got %d: %q", len(parts), got)
	}
}

func TestWriteLine_EmptyMessage(t *testing.T) {
	var buf bytes.Buffer
	w := output.New(&buf)
	w.WriteLine("")

	got := buf.String()
	// Even an empty message should produce a newline-terminated line
	if !strings.HasSuffix(got, "\n") {
		t.Errorf("expected output to end with newline for empty message, got: %q", got)
	}
}
