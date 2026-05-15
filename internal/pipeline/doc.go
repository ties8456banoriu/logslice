// Package pipeline provides a high-level orchestration layer for logslice.
//
// It connects the individual subsystems — file discovery (filefinder), log
// parsing (logparser), time-window filtering (logreader / timewindow), and
// formatted output (output) — into a single, easy-to-use Pipeline type.
//
// Typical usage:
//
//	p, err := pipeline.New(pipeline.Config{
//		Patterns:  []string{"/var/log/app/*.log"},
//		From:      "2024-01-01T00:00:00Z",
//		To:        "2024-01-02T00:00:00Z",
//		Format:    "json",
//		Writer:    os.Stdout,
//	})
//	if err != nil { ... }
//	n, err := p.Run()
package pipeline
