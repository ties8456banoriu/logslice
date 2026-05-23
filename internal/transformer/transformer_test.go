package transformer

import (
	"errors"
	"testing"
)

func TestNew_NoFunctions(t *testing.T) {
	tr, err := New()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr.Len() != 0 {
		t.Fatalf("expected 0 fns, got %d", tr.Len())
	}
}

func TestNew_NilFunction(t *testing.T) {
	_, err := New(nil)
	if err == nil {
		t.Fatal("expected error for nil function")
	}
}

func TestApply_NoOp(t *testing.T) {
	tr, _ := New()
	entry := map[string]any{"msg": "hello"}
	out, err := tr.Apply(entry)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["msg"] != "hello" {
		t.Fatalf("expected msg=hello, got %v", out["msg"])
	}
}

func TestApply_AddField(t *testing.T) {
	addEnv := func(e map[string]any) (map[string]any, error) {
		e["env"] = "test"
		return e, nil
	}
	tr, _ := New(addEnv)
	out, err := tr.Apply(map[string]any{"msg": "hi"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["env"] != "test" {
		t.Fatalf("expected env=test, got %v", out["env"])
	}
}

func TestApply_ChainOrder(t *testing.T) {
	var order []string
	makeStep := func(name string) TransformFunc {
		return func(e map[string]any) (map[string]any, error) {
			order = append(order, name)
			return e, nil
		}
	}
	tr, _ := New(makeStep("a"), makeStep("b"), makeStep("c"))
	_, _ = tr.Apply(map[string]any{})
	if len(order) != 3 || order[0] != "a" || order[1] != "b" || order[2] != "c" {
		t.Fatalf("unexpected order: %v", order)
	}
}

func TestApply_ErrorAbortsChain(t *testing.T) {
	called := false
	failFn := func(e map[string]any) (map[string]any, error) {
		return nil, errors.New("boom")
	}
	afterFn := func(e map[string]any) (map[string]any, error) {
		called = true
		return e, nil
	}
	tr, _ := New(failFn, afterFn)
	_, err := tr.Apply(map[string]any{})
	if err == nil {
		t.Fatal("expected error")
	}
	if called {
		t.Fatal("second function should not have been called")
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	original := map[string]any{"key": "original"}
	mutate := func(e map[string]any) (map[string]any, error) {
		e["key"] = "mutated"
		return e, nil
	}
	tr, _ := New(mutate)
	out, _ := tr.Apply(original)
	if original["key"] != "original" {
		t.Fatalf("original map was mutated")
	}
	if out["key"] != "mutated" {
		t.Fatalf("expected mutated key in output")
	}
}
