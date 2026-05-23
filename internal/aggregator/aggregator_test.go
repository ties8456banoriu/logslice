package aggregator

import (
	"sort"
	"testing"
	"time"
)

func TestNew_DefaultWindow(t *testing.T) {
	a := New(0)
	if a.window != time.Minute {
		t.Fatalf("expected default window 1m, got %v", a.window)
	}
}

func TestNew_NegativeWindow(t *testing.T) {
	a := New(-5 * time.Second)
	if a.window != time.Minute {
		t.Fatalf("expected default window 1m, got %v", a.window)
	}
}

func TestNew_CustomWindow(t *testing.T) {
	a := New(5 * time.Minute)
	if a.window != 5*time.Minute {
		t.Fatalf("expected 5m, got %v", a.window)
	}
}

func TestAdd_And_Counts(t *testing.T) {
	a := New(time.Minute)
	base := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	a.Add("error", base.Add(10*time.Second))
	a.Add("error", base.Add(20*time.Second))
	a.Add("error", base.Add(90*time.Second)) // next bucket
	a.Add("warn", base.Add(5*time.Second))

	counts := a.Counts("error")
	if counts == nil {
		t.Fatal("expected counts for 'error', got nil")
	}
	if got := counts[base]; got != 2 {
		t.Errorf("bucket %v: expected 2, got %d", base, got)
	}
	nextBucket := base.Add(time.Minute)
	if got := counts[nextBucket]; got != 1 {
		t.Errorf("bucket %v: expected 1, got %d", nextBucket, got)
	}

	warnCounts := a.Counts("warn")
	if warnCounts[base] != 1 {
		t.Errorf("warn bucket: expected 1, got %d", warnCounts[base])
	}
}

func TestCounts_UnknownKey(t *testing.T) {
	a := New(time.Minute)
	if a.Counts("missing") != nil {
		t.Error("expected nil for unknown key")
	}
}

func TestKeys(t *testing.T) {
	a := New(time.Minute)
	ts := time.Now()
	a.Add("alpha", ts)
	a.Add("beta", ts)
	a.Add("alpha", ts)

	keys := a.Keys()
	sort.Strings(keys)
	if len(keys) != 2 || keys[0] != "alpha" || keys[1] != "beta" {
		t.Errorf("unexpected keys: %v", keys)
	}
}

func TestReset(t *testing.T) {
	a := New(time.Minute)
	a.Add("x", time.Now())
	a.Reset()
	if len(a.Keys()) != 0 {
		t.Error("expected no keys after reset")
	}
}
