// Package stats tracks processing metrics for a logslice run,
// including file counts, line counts, match counts, and elapsed time.
//
// Usage:
//
//	s := stats.New()
//	s.AddFile()
//	s.AddLine()
//	s.AddMatch()
//	s.Finish()
//	fmt.Println(s.Summary())
package stats
