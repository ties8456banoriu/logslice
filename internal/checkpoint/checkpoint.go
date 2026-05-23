// Package checkpoint provides resumable processing state for log files.
// It tracks the last successfully processed byte offset per file so that
// logslice can resume from where it left off on a subsequent run.
package checkpoint

import (
	"encoding/json"
	"os"
	"sync"
)

// State holds per-file byte offsets.
type State map[string]int64

// Checkpoint persists and retrieves processing offsets for log files.
type Checkpoint struct {
	mu   sync.Mutex
	path string
	data State
}

// New loads an existing checkpoint file at path, or starts with an empty
// state if the file does not yet exist.
func New(path string) (*Checkpoint, error) {
	c := &Checkpoint{path: path, data: make(State)}
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return c, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()
	if err := json.NewDecoder(f).Decode(&c.data); err != nil {
		return nil, err
	}
	return c, nil
}

// Offset returns the last saved byte offset for the given file path.
// If no offset has been recorded, 0 is returned.
func (c *Checkpoint) Offset(file string) int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.data[file]
}

// SetOffset records the byte offset for the given file path in memory.
func (c *Checkpoint) SetOffset(file string, offset int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[file] = offset
}

// Save persists the current in-memory state to disk atomically by writing
// to a temporary file and renaming it over the target path.
func (c *Checkpoint) Save() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	tmp := c.path + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		return err
	}
	if err := json.NewEncoder(f).Encode(c.data); err != nil {
		f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return os.Rename(tmp, c.path)
}
