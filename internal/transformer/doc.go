// Package transformer provides a composable, ordered pipeline for transforming
// structured log entries represented as map[string]any.
//
// A Transformer holds a sequence of TransformFuncs that are applied in
// registration order. Each function receives the output of the previous one,
// forming a chain. If any function returns an error the chain is aborted and
// the error is propagated to the caller.
//
// Usage:
//
//	addHost := func(e map[string]any) (map[string]any, error) {
//		e["host"] = "web-01"
//		return e, nil
//	}
//
//	tr, err := transformer.New(addHost)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	out, err := tr.Apply(entry)
package transformer
