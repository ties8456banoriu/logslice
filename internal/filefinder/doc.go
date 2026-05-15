// Package filefinder resolves log file paths from glob patterns,
// literal paths, and optional recursive directory traversal.
//
// Usage:
//
//	finder := filefinder.New(filefinder.WithRecursive(true))
//	files, err := finder.Expand([]string{"/var/log/app/*.log", "/data/logs"})
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, f := range files {
//		fmt.Println(f)
//	}
//
// Duplicate paths are automatically removed and results are returned
// in sorted order for deterministic processing.
package filefinder
