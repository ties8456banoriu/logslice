// Package checkpoint provides resumable processing state for logslice.
//
// A Checkpoint records the last successfully read byte offset for each
// processed log file. Persisting this state allows logslice to skip
// already-processed content on subsequent invocations, which is useful
// when tailing large or frequently-rotated log directories.
//
// Usage:
//
//	cp, err := checkpoint.New("/var/run/logslice.ckpt")
//	if err != nil { ... }
//
//	// Before reading a file, seek to the saved offset:
//	offset := cp.Offset(filePath)
//
//	// After processing, record progress and persist:
//	cp.SetOffset(filePath, bytesRead)
//	if err := cp.Save(); err != nil { ... }
//
// The checkpoint file is written atomically (write-then-rename) to avoid
// corruption if the process is interrupted mid-save.
package checkpoint
