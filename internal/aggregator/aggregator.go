package aggregator

import (
	"sync"
	"time"
)

// Aggregator groups log entries by a key field and counts occurrences
// within configurable time buckets.
type Aggregator struct {
	mu      sync.Mutex
	buckets map[string]map[time.Time]int64
	window  time.Duration
}

// New creates a new Aggregator with the given bucket window duration.
// If window is zero or negative, it defaults to one minute.
func New(window time.Duration) *Aggregator {
	if window <= 0 {
		window = time.Minute
	}
	return &Aggregator{
		buckets: make(map[string]map[time.Time]int64),
		window:  window,
	}
}

// Add records an occurrence of key at the given timestamp.
func (a *Aggregator) Add(key string, ts time.Time) {
	bucket := ts.Truncate(a.window)
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.buckets[key] == nil {
		a.buckets[key] = make(map[time.Time]int64)
	}
	a.buckets[key][bucket]++
}

// Counts returns a copy of the bucket counts for the given key.
// Returns nil if the key has not been seen.
func (a *Aggregator) Counts(key string) map[time.Time]int64 {
	a.mu.Lock()
	defer a.mu.Unlock()
	src, ok := a.buckets[key]
	if !ok {
		return nil
	}
	out := make(map[time.Time]int64, len(src))
	for t, c := range src {
		out[t] = c
	}
	return out
}

// Keys returns all keys that have been recorded.
func (a *Aggregator) Keys() []string {
	a.mu.Lock()
	defer a.mu.Unlock()
	keys := make([]string, 0, len(a.buckets))
	for k := range a.buckets {
		keys = append(keys, k)
	}
	return keys
}

// Reset clears all accumulated data.
func (a *Aggregator) Reset() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.buckets = make(map[string]map[time.Time]int64)
}
