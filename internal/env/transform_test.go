package env

import (
	"strings"
	"testing"
)

func upperStep(pairs []string) []string {
	out := make([]string, len(pairs))
	for i, p := range pairs {
		out[i] = strings.ToUpper(p)
	}
	return out
}

func appendStep(pairs []string) []string {
	return append(pairs, "EXTRA=1")
}

func TestNewTransformer_Empty(t *testing.T) {
	tr := NewTransformer()
	if tr.Len() != 0 {
		t.Fatalf("expected 0 steps, got %d", tr.Len())
	}
}

func TestApply_NoSteps_Passthrough(t *testing.T) {
	tr := NewTransformer()
	input := []string{"FOO=bar", "BAZ=qux"}
	out := tr.Apply(input)
	if len(out) != len(input) {
		t.Fatalf("expected %d entries, got %d", len(input), len(out))
	}
	for i := range input {
		if out[i] != input[i] {
			t.Errorf("entry %d: expected %q got %q", i, input[i], out[i])
		}
	}
}

func TestApply_SingleStep(t *testing.T) {
	tr := NewTransformer(upperStep)
	out := tr.Apply([]string{"foo=bar"})
	if out[0] != "FOO=BAR" {
		t.Errorf("expected FOO=BAR, got %q", out[0])
	}
}

func TestApply_MultipleSteps_OrderPreserved(t *testing.T) {
	tr := NewTransformer(appendStep, upperStep)
	out := tr.Apply([]string{"foo=bar"})
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	if out[0] != "FOO=BAR" {
		t.Errorf("expected FOO=BAR, got %q", out[0])
	}
	if out[1] != "EXTRA=1" {
		t.Errorf("expected EXTRA=1, got %q", out[1])
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	tr := NewTransformer(appendStep)
	input := []string{"A=1"}
	_ = tr.Apply(input)
	if len(input) != 1 {
		t.Errorf("input was mutated, now has %d entries", len(input))
	}
}

func TestAdd_IncreasesLen(t *testing.T) {
	tr := NewTransformer()
	tr.Add(upperStep)
	if tr.Len() != 1 {
		t.Fatalf("expected 1 step, got %d", tr.Len())
	}
	tr.Add(appendStep, upperStep)
	if tr.Len() != 3 {
		t.Fatalf("expected 3 steps, got %d", tr.Len())
	}
}
