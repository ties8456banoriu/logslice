// Package pipeline provides a high-level orchestration layer that composes
// the filefinder, logreader, and output packages into a single, reusable
// processing pipeline.
//
// # Overview
//
// A pipeline run consists of three stages:
//
//  1. File expansion – one or more path patterns (literals, globs, or
//     directories) are resolved to a deduplicated list of regular files
//     via [filefinder.FileFinder].
//
//  2. Log reading & filtering – each file is passed to [logreader.Reader],
//     which parses structured log lines and retains only those whose
//     timestamps fall within the configured [timewindow.TimeWindow].
//
//  3. Output writing – matching entries are serialised by [output.Writer]
//     in the requested format (JSON or compact).
//
// # Usage
//
//	res, err := pipeline.Run(ctx, pipeline.Config{
//	    Patterns:  []string{"./logs/*.log"},
//	    Window:    window,
//	    Reader:    reader,
//	    Writer:    writer,
//	    Recursive: true,
//	})
//
// The returned [Result] exposes a [stats.Stats] value and the number of
// files that were processed.
package pipeline
