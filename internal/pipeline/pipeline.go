// Package pipeline wires together file finding, log reading, and output
// writing into a single reusable processing pipeline.
package pipeline

import (
	"context"
	"fmt"
	"io"

	"logslice/internal/filefinder"
	"logslice/internal/logreader"
	"logslice/internal/output"
	"logslice/internal/stats"
	"logslice/internal/timewindow"
)

// Config holds all parameters needed to run the pipeline.
type Config struct {
	// Patterns are file paths or glob patterns to process.
	Patterns []string
	// Window is the time range to filter log entries.
	Window *timewindow.TimeWindow
	// Reader is the configured log reader.
	Reader *logreader.Reader
	// Writer is the configured output writer.
	Writer *output.Writer
	// Recursive enables recursive directory expansion.
	Recursive bool
}

// Result summarises the outcome of a pipeline run.
type Result struct {
	Stats     *stats.Stats
	FileCount int
}

// Run executes the full pipeline: expand patterns → read & filter logs → write
// matching entries. It honours ctx cancellation at each stage.
func Run(ctx context.Context, cfg Config) (*Result, error) {
	if cfg.Window == nil {
		return nil, fmt.Errorf("pipeline: time window must not be nil")
	}
	if cfg.Reader == nil {
		return nil, fmt.Errorf("pipeline: reader must not be nil")
	}
	if cfg.Writer == nil {
		return nil, fmt.Errorf("pipeline: writer must not be nil")
	}

	opts := []filefinder.Option{}
	if cfg.Recursive {
		opts = append(opts, filefinder.WithRecursive())
	}
	finder, err := filefinder.New(cfg.Patterns, opts...)
	if err != nil {
		return nil, fmt.Errorf("pipeline: %w", err)
	}

	files, err := finder.Expand()
	if err != nil {
		return nil, fmt.Errorf("pipeline: expand files: %w", err)
	}

	st := stats.New()

	for _, f := range files {
		if ctx.Err() != nil {
			break
		}
		if err := processFile(ctx, f, cfg, st); err != nil {
			return nil, err
		}
	}

	st.Finish()

	return &Result{
		Stats:     st,
		FileCount: len(files),
	}, nil
}

func processFile(ctx context.Context, path string, cfg Config, st *stats.Stats) error {
	entries, err := cfg.Reader.Read(ctx, path)
	if err != nil {
		return fmt.Errorf("pipeline: read %q: %w", path, err)
	}
	for _, e := range entries {
		st.Record(true)
		if err := cfg.Writer.Write(io.Discard, e); err != nil {
			return fmt.Errorf("pipeline: write entry: %w", err)
		}
	}
	return nil
}
