package sample

import (
	"testing"
)

func TestNew_ClampsRateBelow(t *testing.T) {
	s := New(-0.5, 42)
	if s.Rate() != 0 {
		t.Fatalf("expected 0, got %f", s.Rate())
	}
}

func TestNew_ClampsRateAbove(t *testing.T) {
	s := New(1.5, 42)
	if s.Rate() != 1.0 {
		t.Fatalf("expected 1.0, got %f", s.Rate())
	}
}

func TestEnabled_FullRate(t *testing.T) {
	s := New(1.0, 42)
	if s.Enabled() {
		t.Fatal("expected Enabled=false for rate 1.0")
	}
}

func TestEnabled_PartialRate(t *testing.T) {
	s := New(0.5, 42)
	if !s.Enabled() {
		t.Fatal("expected Enabled=true for rate 0.5")
	}
}

func TestKeep_ZeroRate_DropsAll(t *testing.T) {
	s := New(0.0, 42)
	for i := 0; i < 100; i++ {
		if s.Keep() {
			t.Fatal("expected Keep=false for rate 0.0")
		}
	}
}

func TestKeep_FullRate_KeepsAll(t *testing.T) {
	s := New(1.0, 42)
	for i := 0; i < 100; i++ {
		if !s.Keep() {
			t.Fatal("expected Keep=true for rate 1.0")
		}
	}
}

func TestKeep_HalfRate_Approximate(t *testing.T) {
	s := New(0.5, 99)
	kept := 0
	const n = 10_000
	for i := 0; i < n; i++ {
		if s.Keep() {
			kept++
		}
	}
	ratio := float64(kept) / n
	if ratio < 0.45 || ratio > 0.55 {
		t.Fatalf("expected ~0.5 keep ratio, got %f", ratio)
	}
}
