package ratelimit_test

import (
	"context"
	"testing"
	"time"

	"greplive/internal/ratelimit"
)

func TestNew_ZeroRate_NoLimiting(t *testing.T) {
	l := ratelimit.New(0)
	defer l.Stop()

	if l.Rate() != 0 {
		t.Fatalf("expected rate 0, got %d", l.Rate())
	}

	ctx := context.Background()
	// Should never block when rate is 0
	for i := 0; i < 100; i++ {
		if !l.Wait(ctx) {
			t.Fatal("Wait returned false with unlimited limiter")
		}
	}
}

func TestNew_PositiveRate_StoresRate(t *testing.T) {
	l := ratelimit.New(10)
	defer l.Stop()

	if l.Rate() != 10 {
		t.Fatalf("expected rate 10, got %d", l.Rate())
	}
}

func TestWait_ContextCancelled_ReturnsFalse(t *testing.T) {
	// Very low rate so tokens are scarce
	l := ratelimit.New(1)
	defer l.Stop()

	// Drain the initial bucket
	ctx := context.Background()
	l.Wait(ctx)

	// Now the bucket should be empty; cancel the context immediately
	cancelCtx, cancel := context.WithCancel(context.Background())
	cancel()

	result := l.Wait(cancelCtx)
	if result {
		t.Fatal("expected Wait to return false on cancelled context")
	}
}

func TestWait_TokensRefilled_ReturnsTrue(t *testing.T) {
	l := ratelimit.New(100)
	defer l.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	ok := l.Wait(ctx)
	if !ok {
		t.Fatal("expected Wait to succeed with 100 lines/sec limiter")
	}
}

func TestStop_DoesNotPanic(t *testing.T) {
	l := ratelimit.New(50)
	l.Stop()
	// calling Stop on a zero-rate limiter should also be safe
	l2 := ratelimit.New(0)
	l2.Stop()
}

func TestWait_NegativeRate_TreatedAsUnlimited(t *testing.T) {
	l := ratelimit.New(-5)
	defer l.Stop()

	ctx := context.Background()
	for i := 0; i < 50; i++ {
		if !l.Wait(ctx) {
			t.Fatal("Wait should never block for negative rate")
		}
	}
}
