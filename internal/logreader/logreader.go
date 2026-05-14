// Package logreader provides functionality for reading log entries
// from files and filtering them by a time window.
package logreader

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/user/logslice/internal/logparser"
	"github.com/user/logslice/internal/timewindow"
)

// Reader reads log entries from a source and filters by time window.
type Reader struct {
	parser *logparser.JSONParser
	window *timewindow.TimeWindow
}

// New creates a new Reader with the given parser and time window.
func New(parser *logparser.JSONParser, window *timewindow.TimeWindow) *Reader {
	return &Reader{
		parser: parser,
		window: window,
	}
}

// ReadFile opens the named file and returns all log entries within the time window.
func (r *Reader) ReadFile(path string) ([]logparser.Entry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("logreader: open %q: %w", path, err)
	}
	defer f.Close()
	return r.Read(f)
}

// Read reads log entries from the given reader and returns those within the time window.
func (r *Reader) Read(src io.Reader) ([]logparser.Entry, error) {
	var results []logparser.Entry
	scanner := bufio.NewScanner(src)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		if line == "" {
			continue
		}

		entry, err := r.parser.Parse(line)
		if err != nil {
			// Skip unparseable lines; callers can enable strict mode if needed.
			continue
		}

		if r.window.Contains(entry.Time) {
			results = append(results, entry)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("logreader: scan error: %w", err)
	}

	return results, nil
}
