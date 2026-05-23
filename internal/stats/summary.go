package stats

import (
	"fmt"
	"io"
	"strings"
)

// WriteSummary writes a human-readable summary of s to w.
// It reports files scanned, lines read, matches found, match rate,
// and total elapsed time.
func (s *Stats) WriteSummary(w io.Writer) error {
	lines := []string{
		"--- logslice summary ---",
		fmt.Sprintf("files scanned : %d", s.Files()),
		fmt.Sprintf("lines read    : %d", s.Lines()),
		fmt.Sprintf("matches found : %d", s.Matches()),
		fmt.Sprintf("match rate    : %.2f%%", s.MatchRate()*100),
		fmt.Sprintf("elapsed       : %s", s.Elapsed().Round(1e6)),
	}
	_, err := fmt.Fprintln(w, strings.Join(lines, "\n"))
	return err
}

// Summary returns the summary as a string.
func (s *Stats) Summary() string {
	var sb strings.Builder
	_ = s.WriteSummary(&sb)
	return sb.String()
}
