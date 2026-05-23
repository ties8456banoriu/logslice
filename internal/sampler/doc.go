// Package sampler implements deterministic log entry sampling for logslice.
//
// When processing large log files it is sometimes desirable to reduce output
// volume by keeping only a representative fraction of matching entries.
// The Sampler type provides a simple counter-based strategy: given a rate N,
// it retains every Nth entry that passes through the pipeline.
//
// Example usage:
//
//	s, err := sampler.New(10) // keep 1 in every 10 entries
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, entry := range entries {
//		if s.Keep() {
//			output.Write(entry)
//		}
//	}
//
// A rate of 0 or 1 disables sampling and every entry is kept.
// The Sampler is safe for use from a single goroutine; the internal counter
// uses atomic operations so concurrent reads are safe, though the sampling
// sequence is only meaningful when calls are serialised.
package sampler
