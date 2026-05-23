package ratelimiter_test

import (
	"context"
	"testing"
	"time"

	"logslice/internal/ratelimiter"
)

func TestNew_ZeroRate_NoLimit(t *testing.T) {
	rl := ratelimiter.New(0)
	defer rl.Stop()

	if rl.LinesPerSecond() != 0 {
		t.Fatalf("expected 0, got %d", rl.LinesPerSecond())
	}
	// Wait should return immediately with no blocking.
	ctx := context.Background()
	start := time.Now()
	for i := 0; i < 1000; i++ {
		if err := rl.Wait(ctx); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}
	if elapsed := time.Since(start); elapsed > 50*time.Millisecond {
		t.Errorf("expected near-instant return, took %v", elapsed)
	}
}

func TestNew_PositiveRate(t *testing.T) {
	rl := ratelimiter.New(100)
	defer rl.Stop()

	if rl.LinesPerSecond() != 100 {
		t.Fatalf("expected 100, got %d", rl.LinesPerSecond())
	}
}

func TestWait_ContextCancelled(t *testing.T) {
	rl := ratelimiter.New(1) // very slow: 1 token/sec
	defer rl.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	// Drain the first token that may already be buffered.
	_ = rl.Wait(ctx)

	// Second call should time out.
	err := rl.Wait(ctx)
	if err == nil {
		t.Fatal("expected error due to context cancellation, got nil")
	}
}

func TestStop_ReleasesWaiters(t *testing.T) {
	rl := ratelimiter.New(1)

	done := make(chan error, 1)
	go func() {
		// Drain the first available token.
		_ = rl.Wait(context.Background())
		// This second call should unblock when Stop is called.
		done <- rl.Wait(context.Background())
	}()

	time.Sleep(30 * time.Millisecond)
	rl.Stop()

	select {
	case err := <-done:
		if err == nil {
			t.Error("expected non-nil error after Stop")
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("goroutine did not unblock after Stop")
	}
}

func TestNegativeRate_TreatedAsUnlimited(t *testing.T) {
	rl := ratelimiter.New(-5)
	defer rl.Stop()

	ctx := context.Background()
	if err := rl.Wait(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
