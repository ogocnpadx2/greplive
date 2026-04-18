package tail

import "testing"

func TestDefaultConfig_Lines(t *testing.T) {
	c := DefaultConfig()
	if c.Lines != 10 {
		t.Fatalf("expected 10, got %d", c.Lines)
	}
}

func TestConfig_Build_UsesLines(t *testing.T) {
	c := Config{Lines: 5}
	b := c.Build()
	if b.cap != 5 {
		t.Fatalf("expected cap 5, got %d", b.cap)
	}
}

func TestConfig_Build_ZeroLines_UsesDefault(t *testing.T) {
	c := Config{Lines: 0}
	b := c.Build()
	if b.cap != 10 {
		t.Fatalf("expected cap 10, got %d", b.cap)
	}
}

func TestConfig_Build_NegativeLines_UsesDefault(t *testing.T) {
	c := Config{Lines: -3}
	b := c.Build()
	if b.cap != 10 {
		t.Fatalf("expected cap 10, got %d", b.cap)
	}
}
