package timewindow_test

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/timewindow"
)

func TestNew_ValidInputs(t *testing.T) {
	cases := []struct {
		start, end string
	}{
		{"2024-01-01T00:00:00Z", "2024-01-02T00:00:00Z"},
		{"2024-01-01T08:00:00", "2024-01-01T18:00:00"},
		{"2024-01-01 08:00:00", "2024-01-01 18:00:00"},
		{"2024-01-01", "2024-01-02"},
	}
	for _, c := range cases {
		tw, err := timewindow.New(c.start, c.end)
		if err != nil {
			t.Errorf("New(%q, %q) unexpected error: %v", c.start, c.end, err)
		}
		if tw == nil {
			t.Errorf("New(%q, %q) returned nil window", c.start, c.end)
		}
	}
}

func TestNew_InvalidInputs(t *testing.T) {
	cases := []struct {
		start, end string
	}{
		{"not-a-time", "2024-01-02"},
		{"2024-01-02", "not-a-time"},
		{"2024-01-02", "2024-01-01"}, // end before start
		{"2024-01-01", "2024-01-01"}, // equal times
	}
	for _, c := range cases {
		_, err := timewindow.New(c.start, c.end)
		if err == nil {
			t.Errorf("New(%q, %q) expected error but got none", c.start, c.end)
		}
	}
}

func TestContains(t *testing.T) {
	tw, _ := timewindow.New("2024-06-01T10:00:00Z", "2024-06-01T12:00:00Z")

	inside := time.Date(2024, 6, 1, 11, 0, 0, 0, time.UTC)
	if !tw.Contains(inside) {
		t.Error("expected inside time to be contained")
	}

	before := time.Date(2024, 6, 1, 9, 59, 59, 0, time.UTC)
	if tw.Contains(before) {
		t.Error("expected before time to NOT be contained")
	}

	after := time.Date(2024, 6, 1, 12, 0, 1, 0, time.UTC)
	if tw.Contains(after) {
		t.Error("expected after time to NOT be contained")
	}

	// boundary — inclusive
	if !tw.Contains(tw.Start) || !tw.Contains(tw.End) {
		t.Error("expected boundary times to be contained")
	}
}

func TestString(t *testing.T) {
	tw, _ := timewindow.New("2024-06-01T10:00:00Z", "2024-06-01T12:00:00Z")
	s := tw.String()
	if s == "" {
		t.Error("String() returned empty string")
	}
}
