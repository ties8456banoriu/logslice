// Package cli wires together the command-line flags and the core pipeline
// (logparser → logreader → output) for logslice.
package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/example/logslice/internal/logparser"
	"github.com/example/logslice/internal/logreader"
	"github.com/example/logslice/internal/output"
	"github.com/example/logslice/internal/output/format_string"
	"github.com/example/logslice/internal/timewindow"
)

// Config holds all resolved CLI options.
type Config struct {
	From      string
	To        string
	TimeField string
	Format    string
	Files     []string
}

// Run parses args, builds the pipeline, and streams matching log lines to w.
func Run(args []string, w io.Writer, errW io.Writer) error {
	cfg, err := parseFlags(args, errW)
	if err != nil {
		return err
	}

	win, err := timewindow.New(cfg.From, cfg.To)
	if err != nil {
		return fmt.Errorf("time window: %w", err)
	}

	parser, err := logparser.NewJSONParser(cfg.TimeField)
	if err != nil {
		return fmt.Errorf("parser: %w", err)
	}

	fmt, err := format_string.ParseFormat(cfg.Format)
	if err != nil {
		return fmt.Errorf("output format: %w", err)
	}

	out, err := output.New(w, fmt)
	if err != nil {
		return fmt.Errorf("output: %w", err)
	}

	for _, path := range cfg.Files {
		if err := processFile(path, win, parser, out); err != nil {
			fmt.Fprintf(errW, "warn: %s: %v\n", path, err)
		}
	}
	return nil
}

func processFile(path string, win *timewindow.TimeWindow, p *logparser.JSONParser, out *output.Output) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	reader, err := logreader.New(f, win, p)
	if err != nil {
		return err
	}

	for {
		entry, ok, err := reader.Read()
		if err != nil {
			return err
		}
		if !ok {
			break
		}
		if err := out.Write(entry); err != nil {
			return err
		}
	}
	return nil
}

func parseFlags(args []string, errW io.Writer) (*Config, error) {
	fs := flag.NewFlagSet("logslice", flag.ContinueOnError)
	fs.SetOutput(errW)

	from := fs.String("from", "", "start of time window (RFC3339)")
	to := fs.String("to", "", "end of time window (RFC3339)")
	timeField := fs.String("time-field", "time", "JSON field containing the timestamp")
	fmt := fs.String("format", "json", "output format: json|compact")

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil, nil
		}
		return nil, err
	}

	if *from == "" || *to == "" {
		return nil, errors.New("--from and --to are required")
	}

	if fs.NArg() == 0 {
		return nil, errors.New("at least one log file must be specified")
	}

	return &Config{
		From:      *from,
		To:        *to,
		TimeField: *timeField,
		Format:    *fmt,
		Files:     fs.Args(),
	}, nil
}
