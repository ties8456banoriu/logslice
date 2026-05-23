package stats

import (
	"fmt"
	"io"
	"time"
)

// Counter tracks processing statistics for a logslice run.
type Counter struct {
	FilesProcessed int
	LinesRead      int
	LinesMatched   int
	LinesSkipped   int
	StartTime      time.Time
	EndTime        time.Time
}

// New returns a new Counter with the start time set to now.
func New() *Counter {
	return &Counter{StartTime: time.Now()}
}

// Finish records the end time of processing.
func (c *Counter) Finish() {
	c.EndTime = time.Now()
}

// Elapsed returns the duration between start and end.
// If Finish has not been called, it returns the duration since start.
func (c *Counter) Elapsed() time.Duration {
	if c.EndTime.IsZero() {
		return time.Since(c.StartTime)
	}
	return c.EndTime.Sub(c.StartTime)
}

// MatchRate returns the fraction of read lines that matched the window.
// Returns 0 if no lines were read.
func (c *Counter) MatchRate() float64 {
	if c.LinesRead == 0 {
		return 0
	}
	return float64(c.LinesMatched) / float64(c.LinesRead)
}

// Print writes a human-readable summary to w.
func (c *Counter) Print(w io.Writer) {
	fmt.Fprintf(w, "files processed : %d\n", c.FilesProcessed)
	fmt.Fprintf(w, "lines read      : %d\n", c.LinesRead)
	fmt.Fprintf(w, "lines matched   : %d\n", c.LinesMatched)
	fmt.Fprintf(w, "lines skipped   : %d\n", c.LinesSkipped)
	fmt.Fprintf(w, "match rate      : %.1f%%\n", c.MatchRate()*100)
	fmt.Fprintf(w, "elapsed         : %s\n", c.Elapsed().Round(time.Millisecond))
}
