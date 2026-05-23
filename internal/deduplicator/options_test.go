package deduplicator

import "testing"

func TestDefaultConfig(t *testing.T) {
	cfg := defaultConfig()
	if cfg.capacity != 1024 {
		t.Errorf("expected default capacity 1024, got %d", cfg.capacity)
	}
}

func TestWithCapacity_Positive(t *testing.T) {
	d := NewFromOptions(WithCapacity(256))
	if d.capacity != 256 {
		t.Errorf("expected capacity 256, got %d", d.capacity)
	}
}

func TestWithCapacity_Zero_UsesDefault(t *testing.T) {
	d := NewFromOptions(WithCapacity(0))
	if d.capacity != 1024 {
		t.Errorf("expected default capacity 1024, got %d", d.capacity)
	}
}

func TestWithCapacity_Negative_UsesDefault(t *testing.T) {
	d := NewFromOptions(WithCapacity(-10))
	if d.capacity != 1024 {
		t.Errorf("expected default capacity 1024, got %d", d.capacity)
	}
}

func TestNewFromOptions_NoOptions(t *testing.T) {
	d := NewFromOptions()
	if d.capacity != 1024 {
		t.Errorf("expected default capacity 1024, got %d", d.capacity)
	}
	if d.Len() != 0 {
		t.Errorf("expected empty deduplicator, got Len=%d", d.Len())
	}
}

func TestNewFromOptions_MultipleOptions(t *testing.T) {
	d := NewFromOptions(WithCapacity(512), WithCapacity(64))
	// last write wins
	if d.capacity != 64 {
		t.Errorf("expected capacity 64, got %d", d.capacity)
	}
}
