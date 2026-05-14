// Package cli implements the command-line interface for logslice.
//
// It parses flags (--from, --to, --time-field, --format) and one or more
// log file paths, then wires together the internal pipeline:
//
//	logparser.JSONParser  →  logreader.Reader  →  output.Output
//
// Each file is processed sequentially; parse or read errors for a single
// file are reported as warnings on stderr so that remaining files are still
// processed.
//
// Usage:
//
//	logslice --from <RFC3339> --to <RFC3339> [--time-field <field>] [--format json|compact] <file>...
package cli
