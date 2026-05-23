// Package stats provides a lightweight counter for tracking logslice
// processing metrics such as files processed, lines read, lines matched,
// and elapsed time.
//
// Usage:
//
//	c := stats.New()
//	// ... process files, incrementing c.FilesProcessed, c.LinesRead, etc.
//	c.Finish()
//	c.Print(os.Stderr)
package stats
