// Package filefinder provides utilities for discovering log files
// matching glob patterns or directory traversal.
package filefinder

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// Finder locates log files based on provided patterns or paths.
type Finder struct {
	recursive bool
}

// Option configures a Finder.
type Option func(*Finder)

// WithRecursive enables recursive directory traversal.
func WithRecursive(r bool) Option {
	return func(f *Finder) {
		f.recursive = r
	}
}

// New creates a new Finder with the given options.
func New(opts ...Option) *Finder {
	f := &Finder{}
	for _, o := range opts {
		o(f)
	}
	return f
}

// Expand resolves a list of patterns/paths into concrete file paths.
// Glob patterns are expanded; directories are walked if recursive is set.
func (f *Finder) Expand(patterns []string) ([]string, error) {
	seen := make(map[string]struct{})
	var results []string

	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return nil, fmt.Errorf("invalid pattern %q: %w", pattern, err)
		}
		if matches == nil {
			// treat as literal path
			matches = []string{pattern}
		}
		for _, m := range matches {
			paths, err := f.expand(m)
			if err != nil {
				return nil, err
			}
			for _, p := range paths {
				if _, ok := seen[p]; !ok {
					seen[p] = struct{}{}
					results = append(results, p)
				}
			}
		}
	}
	sort.Strings(results)
	return results, nil
}

func (f *Finder) expand(path string) ([]string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("cannot stat %q: %w", path, err)
	}
	if !info.IsDir() {
		return []string{path}, nil
	}
	if !f.recursive {
		return nil, fmt.Errorf("%q is a directory; use --recursive to traverse", path)
	}
	var files []string
	err = filepath.Walk(path, func(p string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !fi.IsDir() {
			files = append(files, p)
		}
		return nil
	})
	return files, err
}
