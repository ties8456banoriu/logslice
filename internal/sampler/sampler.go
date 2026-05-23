// Package sampler provides log entry sampling strategies for reducing
// output volume when processing high-throughput log streams.
package sampler

import (
	"errors"
	"sync/atomic"
)

// Sampler decides whether a given log entry should be included in output.
type Sampler struct {
	rate    uint64 // keep 1 in every rate entries
	counter atomic.Uint64
}

// New creates a Sampler that retains every nth entry.
// A rate of 0 or 1 means keep all entries (no sampling).
// Returns an error if rate is negative (not applicable for uint, but validated
// via the int parameter to catch misuse).
func New(rate int) (*Sampler, error) {
	if rate < 0 {
		return nil, errors.New("sampler: rate must be non-negative")
	}
	r := uint64(rate)
	if r <= 1 {
		r = 1
	}
	return &Sampler{rate: r}, nil
}

// Keep reports whether the current entry should be kept.
// It increments an internal counter on every call and returns true
// when the counter is a multiple of the configured rate.
func (s *Sampler) Keep() bool {
	n := s.counter.Add(1)
	return n%s.rate == 0
}

// Reset resets the internal counter to zero.
func (s *Sampler) Reset() {
	s.counter.Store(0)
}

// Rate returns the configured sampling rate.
func (s *Sampler) Rate() uint64 {
	return s.rate
}
