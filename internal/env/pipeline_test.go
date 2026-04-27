package env

import (
	"errors"
	"testing"
)

func appendSuffix(suffix string) Stage {
	return func(pairs []string) ([]string, error) {
		out := make([]string, len(pairs))
		for i, p := range pairs {
			out[i] = p + suffix
		}
		return out, nil
	}
}

func failStage(pairs []string) ([]string, error) {
	return nil, errors.New("stage failed")
}

func TestNewPipeline_Empty(t *testing.T) {
	p := NewPipeline()
	if p.Len() != 0 {
		t.Fatalf("expected 0 stages, got %d", p.Len())
	}
}

func TestRun_NoStages_Passthrough(t *testing.T) {
	input := []string{"A=1", "B=2"}
	p := NewPipeline()
	out, err := p.Run(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != len(input) || out[0] != input[0] {
		t.Fatalf("expected passthrough, got %v", out)
	}
}

func TestRun_SingleStage(t *testing.T) {
	p := NewPipeline(appendSuffix("_X"))
	out, err := p.Run([]string{"A=1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0] != "A=1_X" {
		t.Fatalf("expected A=1_X, got %s", out[0])
	}
}

func TestRun_MultipleStages_Chained(t *testing.T) {
	p := NewPipeline(appendSuffix("_A"), appendSuffix("_B"))
	out, err := p.Run([]string{"V=1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0] != "V=1_A_B" {
		t.Fatalf("expected V=1_A_B, got %s", out[0])
	}
}

func TestRun_StageError_HaltsPipeline(t *testing.T) {
	called := false
	guard := func(pairs []string) ([]string, error) {
		called = true
		return pairs, nil
	}
	p := NewPipeline(failStage, guard)
	_, err := p.Run([]string{"A=1"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if called {
		t.Fatal("subsequent stage should not have been called after error")
	}
}

func TestAdd_IncreasesLen(t *testing.T) {
	p := NewPipeline()
	p.Add(appendSuffix("_1"), appendSuffix("_2"))
	if p.Len() != 2 {
		t.Fatalf("expected 2 stages, got %d", p.Len())
	}
}

func TestRun_NilInput_ReturnsNil(t *testing.T) {
	p := NewPipeline(appendSuffix("_X"))
	out, err := p.Run(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Fatalf("expected empty output, got %v", out)
	}
}
