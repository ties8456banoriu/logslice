package stats_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/logslice/internal/stats"
)

func TestNew_SetsStartTime(t *testing.T) {
	before := time.Now()
	c := stats.New()
	after := time.Now()

	if c.StartTime.Before(before) || c.StartTime.After(after) {
		t.Errorf("StartTime %v not in expected range [%v, %v]", c.StartTime, before, after)
	}
}

func TestFinish_SetsEndTime(t *testing.T) {
	c := stats.New()
	if !c.EndTime.IsZero() {
		t.Fatal("EndTime should be zero before Finish")
	}
	c.Finish()
	if c.EndTime.IsZero() {
		t.Fatal("EndTime should be set after Finish")
	}
}

func TestElapsed_BeforeFinish(t *testing.T) {
	c := stats.New()
	time.Sleep(2 * time.Millisecond)
	if c.Elapsed() < time.Millisecond {
		t.Error("expected elapsed > 1ms before Finish")
	}
}

func TestElapsed_AfterFinish(t *testing.T) {
	c := stats.New()
	time.Sleep(2 * time.Millisecond)
	c.Finish()
	e1 := c.Elapsed()
	time.Sleep(5 * time.Millisecond)
	e2 := c.Elapsed()
	if e1 != e2 {
		t.Errorf("Elapsed changed after Finish: %v vs %v", e1, e2)
	}
}

func TestMatchRate(t *testing.T) {
	cases := []struct {
		read, matched int
		want          float64
	}{
		{0, 0, 0},
		{10, 5, 0.5},
		{10, 10, 1.0},
		{10, 0, 0},
	}
	for _, tc := range cases {
		c := stats.New()
		c.LinesRead = tc.read
		c.LinesMatched = tc.matched
		if got := c.MatchRate(); got != tc.want {
			t.Errorf("MatchRate(%d/%d) = %v, want %v", tc.matched, tc.read, got, tc.want)
		}
	}
}

func TestPrint_ContainsExpectedFields(t *testing.T) {
	c := stats.New()
	c.FilesProcessed = 3
	c.LinesRead = 100
	c.LinesMatched = 42
	c.LinesSkipped = 58
	c.Finish()

	var buf bytes.Buffer
	c.Print(&buf)
	out := buf.String()

	for _, want := range []string{"3", "100", "42", "58", "42.0%"} {
		if !strings.Contains(out, want) {
			t.Errorf("Print output missing %q:\n%s", want, out)
		}
	}
}
