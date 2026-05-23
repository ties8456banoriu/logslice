// Package aggregator provides time-bucket aggregation of log entry keys.
//
// It is useful for summarising high-volume log streams by grouping events
// into fixed-width time windows (e.g. per minute) and counting occurrences
// of each distinct key (such as a log level or error code).
//
// Basic usage:
//
//	a := aggregator.New(time.Minute)
//	a.Add("error", entry.Timestamp)
//	counts := a.Counts("error")
//
// Use NewFromOptions for a functional-option style:
//
//	a := aggregator.NewFromOptions(
//		aggregator.WithWindow(5 * time.Minute),
//	)
//
// The Aggregator is safe for concurrent use.
package aggregator
