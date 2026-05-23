// Package deduplicator provides a content-hash-based deduplication filter
// for log entries processed by logslice.
//
// It maintains a bounded sliding window of SHA-256 hashes of previously seen
// log lines. When the window is full, the oldest entry is evicted to make room
// for new ones, preventing unbounded memory growth during long-running scans.
//
// Basic usage:
//
//	d := deduplicator.New(512)
//	for _, line := range lines {
//		if !d.IsDuplicate(line) {
//			// process unique line
//		}
//	}
//
// Functional options are available via NewFromOptions:
//
//	d := deduplicator.NewFromOptions(
//		deduplicator.WithCapacity(2048),
//	)
//
// Deduplicator is safe for concurrent use.
package deduplicator
