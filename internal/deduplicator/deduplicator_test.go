package deduplicator

import (
	"fmt"
	"testing"
)

func TestNew_DefaultCapacity(t *testing.T) {
	d := New(0)
	if d.capacity != 1024 {
		t.Errorf("expected default capacity 1024, got %d", d.capacity)
	}
}

func TestNew_NegativeCapacity(t *testing.T) {
	d := New(-5)
	if d.capacity != 1024 {
		t.Errorf("expected default capacity 1024, got %d", d.capacity)
	}
}

func TestIsDuplicate_FirstSeen(t *testing.T) {
	d := New(10)
	if d.IsDuplicate([]byte("hello")) {
		t.Error("expected first occurrence not to be a duplicate")
	}
}

func TestIsDuplicate_SecondSeen(t *testing.T) {
	d := New(10)
	d.IsDuplicate([]byte("hello"))
	if !d.IsDuplicate([]byte("hello")) {
		t.Error("expected second occurrence to be a duplicate")
	}
}

func TestIsDuplicate_DifferentData(t *testing.T) {
	d := New(10)
	d.IsDuplicate([]byte("hello"))
	if d.IsDuplicate([]byte("world")) {
		t.Error("expected different data not to be a duplicate")
	}
}

func TestIsDuplicate_EvictsOldest(t *testing.T) {
	d := New(3)
	d.IsDuplicate([]byte("a"))
	d.IsDuplicate([]byte("b"))
	d.IsDuplicate([]byte("c"))
	// "a" should be evicted now
	d.IsDuplicate([]byte("d"))
	if d.IsDuplicate([]byte("a")) {
		t.Error("expected evicted entry 'a' not to be a duplicate")
	}
}

func TestLen(t *testing.T) {
	d := New(10)
	for i := 0; i < 5; i++ {
		d.IsDuplicate([]byte(fmt.Sprintf("entry-%d", i)))
	}
	if d.Len() != 5 {
		t.Errorf("expected Len() == 5, got %d", d.Len())
	}
}

func TestReset_ClearsState(t *testing.T) {
	d := New(10)
	d.IsDuplicate([]byte("hello"))
	d.Reset()
	if d.Len() != 0 {
		t.Errorf("expected Len() == 0 after Reset, got %d", d.Len())
	}
	if d.IsDuplicate([]byte("hello")) {
		t.Error("expected entry to not be duplicate after Reset")
	}
}

func TestIsDuplicate_Concurrent(t *testing.T) {
	d := New(100)
	done := make(chan struct{})
	for i := 0; i < 10; i++ {
		go func(i int) {
			d.IsDuplicate([]byte(fmt.Sprintf("key-%d", i)))
			done <- struct{}{}
		}(i)
	}
	for i := 0; i < 10; i++ {
		<-done
	}
}
