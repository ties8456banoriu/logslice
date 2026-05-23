package transformer

import "fmt"

// TransformFunc is a function that transforms a log entry represented as a
// map of string keys to arbitrary values. It returns the transformed map and
// any error encountered.
type TransformFunc func(map[string]any) (map[string]any, error)

// Transformer applies an ordered chain of TransformFuncs to each log entry.
type Transformer struct {
	fns []TransformFunc
}

// New returns a new Transformer that will apply the given functions in order.
// If no functions are provided the transformer acts as a no-op pass-through.
func New(fns ...TransformFunc) (*Transformer, error) {
	for i, fn := range fns {
		if fn == nil {
			return nil, fmt.Errorf("transformer: function at index %d is nil", i)
		}
	}
	return &Transformer{fns: fns}, nil
}

// Apply runs every registered TransformFunc against entry in registration
// order. Each function receives the output of the previous one. The original
// map is never mutated; a shallow copy is made before the first transform.
// An error from any function aborts the chain and is returned immediately.
func (t *Transformer) Apply(entry map[string]any) (map[string]any, error) {
	if len(t.fns) == 0 {
		return entry, nil
	}

	// shallow copy so callers retain their original map
	current := make(map[string]any, len(entry))
	for k, v := range entry {
		current[k] = v
	}

	for _, fn := range t.fns {
		var err error
		current, err = fn(current)
		if err != nil {
			return nil, fmt.Errorf("transformer: %w", err)
		}
	}
	return current, nil
}

// Len returns the number of transform functions registered.
func (t *Transformer) Len() int { return len(t.fns) }
