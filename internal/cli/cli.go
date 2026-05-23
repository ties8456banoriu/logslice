// Package cli wires together all logslice sub-packages and exposes a single
// Run entry-point consumed by cmd/logslice/main.go.
package cli

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"logslice/internal/filefinder"
	"logslice/internal/logparser"
	"logslice/internal/logreader"
	"logslice/internal/output"
	"logslice/internal/ratelimiter"
	"logslice/internal/stats"
	"logslice/internal/timewindow"
)

// Run is the main entry-point for the CLI. It parses flags, resolves files,
// and writes matching log lines to w.
func Run(args []string, w io.Writer) error {
	cfg, err := parseFlags(args)
	if err != nil {
		return err
	}

	win, err := timewindow.New(cfg.from, cfg.to)
	if err != nil {
		return fmt.Errorf("invalid time window: %w", err)
	}

	parser, err := logparser.NewParser(cfg.format)
	if err != nil {
		return fmt.Errorf("log format: %w", err)
	}

	out, err := output.New(w, cfg.outputFmt)
	if err != nil {
		return fmt.Errorf("output format: %w", err)
	}

	ff := filefinder.New(filefinder.WithRecursive(cfg.recursive))
	files, err := ff.Expand(cfg.paths)
	if err != nil {
		return fmt.Errorf("resolving files: %w", err)
	}
	if len(files) == 0 {
		return fmt.Errorf("no files matched")
	}

	rl := ratelimiter.New(cfg.rateLimit)
	defer rl.Stop()

	st := stats.New()
	ctx := context.Background()

	for _, path := range files {
		if err := processFile(ctx, path, win, parser, out, rl, st); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "warning: %s: %v\n", path, err)
		}
	}

	st.Finish()
	if cfg.showStats {
		_, _ = fmt.Fprint(os.Stderr, st.Summary())
	}
	return nil
}

func processFile(
	ctx context.Context,
	path string,
	win *timewindow.TimeWindow,
	parser logparser.Parser,
	out *output.Output,
	rl *ratelimiter.RateLimiter,
	st *stats.Stats,
) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	reader := logreader.New(f, win, parser)
	for {
		entry, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			st.RecordError()
			continue
		}
		st.RecordLine()
		if entry == nil {
			continue
		}
		st.RecordMatch()
		if err := rl.Wait(ctx); err != nil {
			return err
		}
		if err := out.Write(entry); err != nil {
			return err
		}
	}
	return nil
}

type config struct {
	from      string
	to        string
	format    string
	outputFmt string
	paths     []string
	recursive bool
	showStats bool
	rateLimit int
}

func parseFlags(args []string) (*config, error) {
	fs := flag.NewFlagSet("logslice", flag.ContinueOnError)
	cfg := &config{}

	fs.StringVar(&cfg.from, "from", "", "start of time window (RFC3339)")
	fs.StringVar(&cfg.to, "to", "", "end of time window (RFC3339)")
	fs.StringVar(&cfg.format, "format", "json", "log format (json)")
	fs.StringVar(&cfg.outputFmt, "output", "json", "output format (json|compact)")
	fs.BoolVar(&cfg.recursive, "recursive", false, "recurse into directories")
	fs.BoolVar(&cfg.showStats, "stats", false, "print match statistics to stderr")
	fs.IntVar(&cfg.rateLimit, "rate", 0, "max lines per second (0 = unlimited)")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}
	if cfg.from == "" || cfg.to == "" {
		return nil, fmt.Errorf("--from and --to are required")
	}
	cfg.paths = fs.Args()
	if len(cfg.paths) == 0 {
		return nil, fmt.Errorf("at least one file or pattern is required")
	}
	if !strings.Contains(cfg.outputFmt, cfg.outputFmt) {
		return nil, fmt.Errorf("unknown output format: %s", cfg.outputFmt)
	}
	return cfg, nil
}
