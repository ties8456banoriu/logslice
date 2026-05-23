package deduplicator

import (
	"crypto/sha256"
	"encoding/hex"
	"sync"
)

// Deduplicator filters out duplicate log entries within a sliding window
// of recently seen content hashes.
type Deduplicator struct {
	mu      sync.Mutex
	seen    map[string]struct{}
	order   []string
	capacity int
}

// New creates a Deduplicator that remembers up to capacity recent entries.
// If capacity is zero or negative, a default of 1024 is used.
func New(capacity int) *Deduplicator {
	if capacity <= 0 {
		capacity = 1024
	}
	return &Deduplicator{
		seen:     make(map[string]struct{}, capacity),
		order:    make([]string, 0, capacity),
		capacity: capacity,
	}
}

// IsDuplicate returns true if the given data has been seen before.
// If not a duplicate, the hash is recorded for future calls.
func (d *Deduplicator) IsDuplicate(data []byte) bool {
	h := hash(data)

	d.mu.Lock()
	defer d.mu.Unlock()

	if _, exists := d.seen[h]; exists {
		return true
	}

	if len(d.order) >= d.capacity {
		oldest := d.order[0]
		d.order = d.order[1:]
		delete(d.seen, oldest)
	}

	d.seen[h] = struct{}{}
	d.order = append(d.order, h)
	return false
}

// Reset clears all recorded hashes.
func (d *Deduplicator) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.seen = make(map[string]struct{}, d.capacity)
	d.order = d.order[:0]
}

// Len returns the number of unique entries currently tracked.
func (d *Deduplicator) Len() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.seen)
}

func hash(data []byte) string {
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}
