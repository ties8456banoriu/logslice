// Package cli wires together the logslice pipeline and exposes Run.
package cli

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/user/logslice/internal/filefinder"
	"github.com/user/logslice/internal/logparser"
	"github.com/user/logslice/internal/logreader"
	"github.com/user/logslice/internal/output"
	"github.com/user/logslice/internal/stats"
	"github.com/user/logslice/internal/timewindow"
)

// Run is the entry-point for the logslice CLI.
func Run(args []string, stdout, stderr io.Writer) int {
	flags, paths, err := parseFlags(args, stderr)
	if err != nil {
		return 2
	}

	win, err := timewindow.New(flags.from, flags.to)
	if err != nil {
		fmt.Fprintf(stderr, "error: invalid time window: %v\n", err)
		return 1
	}

	parser, err := logparser.NewParser(flags.format)
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}

	out, err := output.New(stdout, flags.outFormat)
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}

	finder := filefinder.New(filefinder.WithRecursive(flags.recursive))
	files, err := finder.Expand(paths)
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}
	if len(files) == 0 {
		fmt.Fprintln(stderr, "error: no input files found")
		return 1
	}

	ctr := stats.New()

	for _, f := range files {
		if err := processFile(f, win, parser, out, ctr); err != nil {
			fmt.Fprintf(stderr, "warning: %s: %v\n", f, err)
		}
		ctr.FilesProcessed++
	}

	ctr.Finish()
	if flags.printStats {
		ctr.Print(stderr)
	}
	return 0
}

func processFile(path string, win *timewindow.Window, p logparser.Parser, out *output.Output, ctr *stats.Counter) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	reader := logreader.New(f, win, p)
	for reader.Next() {
		ctr.LinesRead++
		entry := reader.Entry()
		if err := out.Write(entry); err != nil {
			return err
		}
		ctr.LinesMatched++
	}
	if err := reader.Err(); err != nil {
		return err
	}
	ctr.LinesRead += reader.Skipped()
	ctr.LinesSkipped += reader.Skipped()
	return nil
}

type cliFlags struct {
	from       string
	to         string
	format     string
	outFormat  string
	recursive  bool
	printStats bool
}

func parseFlags(args []string, stderr io.Writer) (cliFlags, []string, error) {
	fs := flag.NewFlagSet("logslice", flag.ContinueOnError)
	fs.SetOutput(stderr)

	var f cliFlags
	fs.StringVar(&f.from, "from", "", "start of time window (RFC3339)")
	fs.StringVar(&f.to, "to", "", "end of time window (RFC3339)")
	fs.StringVar(&f.format, "format", "json", "log format (json)")
	fs.StringVar(&f.outFormat, "out", "json", "output format (json|compact)")
	fs.BoolVar(&f.recursive, "r", false, "recurse into directories")
	fs.BoolVar(&f.printStats, "stats", false, "print processing statistics to stderr")

	if err := fs.Parse(args); err != nil {
		return f, nil, err
	}
	if f.from == "" || f.to == "" {
		fmt.Fprintln(stderr, "error: --from and --to are required")
		fs.Usage()
		return f, nil, fmt.Errorf("missing required flags")
	}
	return f, fs.Args(), nil
}
