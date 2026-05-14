package timewindow

import (
	"fmt"
	"time"
)

// TimeWindow represents a start and end time range for log filtering.
type TimeWindow struct {
	Start time.Time
	End   time.Time
}

// New creates a new TimeWindow from start and end time strings.
// Accepted formats: RFC3339, "2006-01-02T15:04:05", "2006-01-02 15:04:05", "2006-01-02".
func New(start, end string) (*TimeWindow, error) {
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	parsed := func(s string) (time.Time, error) {
		for _, f := range formats {
			if t, err := time.Parse(f, s); err == nil {
				return t, nil
			}
		}
		return time.Time{}, fmt.Errorf("unrecognised time format: %q", s)
	}

	s, err := parsed(start)
	if err != nil {
		return nil, fmt.Errorf("invalid start time: %w", err)
	}

	e, err := parsed(end)
	if err != nil {
		return nil, fmt.Errorf("invalid end time: %w", err)
	}

	if !e.After(s) {
		return nil, fmt.Errorf("end time must be after start time")
	}

	return &TimeWindow{Start: s, End: e}, nil
}

// Contains reports whether t falls within the window [Start, End] (inclusive).
func (tw *TimeWindow) Contains(t time.Time) bool {
	return !t.Before(tw.Start) && !t.After(tw.End)
}

// String returns a human-readable representation of the window.
func (tw *TimeWindow) String() string {
	return fmt.Sprintf("[%s, %s]", tw.Start.Format(time.RFC3339), tw.End.Format(time.RFC3339))
}
