package sampler_test

import (
	"testing"

	"logslice/internal/sampler"
)

func TestNew_RateOne_KeepsAll(t *testing.T) {
	s, err := sampler.New(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 0; i < 10; i++ {
		if !s.Keep() {
			t.Errorf("expected Keep()=true for entry %d with rate=1", i)
		}
	}
}

func TestNew_RateZero_TreatedAsOne(t *testing.T) {
	s, err := sampler.New(0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Rate() != 1 {
		t.Errorf("expected rate=1, got %d", s.Rate())
	}
}

func TestNew_NegativeRate_ReturnsError(t *testing.T) {
	_, err := sampler.New(-1)
	if err == nil {
		t.Fatal("expected error for negative rate, got nil")
	}
}

func TestKeep_SamplesCorrectly(t *testing.T) {
	const rate = 3
	s, err := sampler.New(rate)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []bool{false, false, true, false, false, true, false, false, true}
	for i, want := range expected {
		got := s.Keep()
		if got != want {
			t.Errorf("entry %d: Keep()=%v, want %v", i+1, got, want)
		}
	}
}

func TestReset_ResetsCounter(t *testing.T) {
	s, _ := sampler.New(3)
	s.Keep() // 1
	s.Keep() // 2
	s.Reset()
	// After reset counter is 0; next Keep increments to 1 — not a multiple of 3.
	if s.Keep() {
		t.Error("expected Keep()=false immediately after Reset with rate=3")
	}
}

func TestRate_ReturnsConfiguredRate(t *testing.T) {
	s, _ := sampler.New(5)
	if s.Rate() != 5 {
		t.Errorf("expected rate=5, got %d", s.Rate())
	}
}

func TestKeep_HighRate_LowThroughput(t *testing.T) {
	s, _ := sampler.New(100)
	kept := 0
	for i := 0; i < 1000; i++ {
		if s.Keep() {
			kept++
		}
	}
	if kept != 10 {
		t.Errorf("expected 10 kept entries out of 1000 at rate=100, got %d", kept)
	}
}
