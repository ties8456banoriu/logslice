// Package pipeline wires together file discovery, log parsing, time-window
// filtering, and output writing into a single reusable processing pipeline.
package pipeline

import (
	"fmt"
	"io"

	"github.com/yourorg/logslice/internal/filefinder"
	"github.com/yourorg/logslice/internal/logparser"
	"github.com/yourorg/logslice/internal/logreader"
	"github.com/yourorg/logslice/internal/output"
	"github.com/yourorg/logslice/internal/timewindow"
)

// Config holds all parameters needed to run the pipeline.
type Config struct {
	Patterns  []string
	From      string
	To        string
	Format    string
	TimeField  string
	TimeLayout string
	Recursive bool
	Writer    io.Writer
}

// Pipeline coordinates the end-to-end log-slice workflow.
type Pipeline struct {
	cfg Config
}

// New creates a Pipeline from the given Config.
func New(cfg Config) (*Pipeline, error) {
	if cfg.Writer == nil {
		return nil, fmt.Errorf("pipeline: writer must not be nil")
	}
	if len(cfg.Patterns) == 0 {
		return nil, fmt.Errorf("pipeline: at least one file pattern is required")
	}
	return &Pipeline{cfg: cfg}, nil
}

// Run executes the pipeline and returns the total number of log entries written.
func (p *Pipeline) Run() (int, error) {
	win, err := timewindow.New(p.cfg.From, p.cfg.To)
	if err != nil {
		return 0, fmt.Errorf("pipeline: invalid time window: %w", err)
	}

	var parserOpts []logparser.Option
	if p.cfg.TimeField != "" {
		parserOpts = append(parserOpts, logparser.WithTimeField(p.cfg.TimeField))
	}
	if p.cfg.TimeLayout != "" {
		parserOpts = append(parserOpts, logparser.WithTimeLayout(p.cfg.TimeLayout))
	}
	parser, err := logparser.NewParser(p.cfg.Format, parserOpts...)
	if err != nil {
		return 0, fmt.Errorf("pipeline: %w", err)
	}

	out, err := output.New(p.cfg.Format, p.cfg.Writer)
	if err != nil {
		return 0, fmt.Errorf("pipeline: %w", err)
	}

	finderOpts := []filefinder.Option{}
	if p.cfg.Recursive {
		finderOpts = append(finderOpts, filefinder.WithRecursive())
	}
	finder, err := filefinder.New(p.cfg.Patterns, finderOpts...)
	if err != nil {
		return 0, fmt.Errorf("pipeline: %w", err)
	}

	files, err := finder.Expand()
	if err != nil {
		return 0, fmt.Errorf("pipeline: file expansion: %w", err)
	}

	total := 0
	for _, f := range files {
		n, err := processFile(f, win, parser, out)
		if err != nil {
			return total, fmt.Errorf("pipeline: processing %q: %w", f, err)
		}
		total += n
	}
	return total, nil
}

func processFile(path string, win *timewindow.TimeWindow, parser logparser.Parser, out *output.Output) (int, error) {
	reader, err := logreader.New(path, win, parser)
	if err != nil {
		return 0, err
	}
	defer reader.Close()

	count := 0
	for {
		entry, ok, err := reader.Read()
		if err != nil {
			return count, err
		}
		if !ok {
			break
		}
		if err := out.Write(entry); err != nil {
			return count, err
		}
		count++
	}
	return count, nil
}
