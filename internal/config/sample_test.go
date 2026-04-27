package config

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/env"
)

func TestResolveSampler_NilConfig_Passthrough(t *testing.T) {
	s := ResolveSampler(nil)
	pairs := []string{"A=1", "B=2", "C=3"}
	out := s.Apply(pairs)
	if len(out) != len(pairs) {
		t.Fatalf("expected passthrough, got %d entries", len(out))
	}
}

func TestResolveSampler_EmptyMode_Passthrough(t *testing.T) {
	s := ResolveSampler(&SampleConfig{Mode: "", N: 2})
	pairs := []string{"A=1", "B=2", "C=3"}
	out := s.Apply(pairs)
	if len(out) != len(pairs) {
		t.Fatalf("expected passthrough for empty mode, got %d entries", len(out))
	}
}

func TestResolveSampler_InvalidMode_Passthrough(t *testing.T) {
	s := ResolveSampler(&SampleConfig{Mode: "bogus", N: 1})
	pairs := []string{"A=1", "B=2", "C=3"}
	out := s.Apply(pairs)
	if len(out) != len(pairs) {
		t.Fatalf("expected passthrough for invalid mode, got %d entries", len(out))
	}
}

func TestResolveSampler_FirstMode(t *testing.T) {
	s := ResolveSampler(&SampleConfig{Mode: "first", N: 2})
	pairs := []string{"A=1", "B=2", "C=3"}
	out := s.Apply(pairs)
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	if out[0] != "A=1" || out[1] != "B=2" {
		t.Errorf("unexpected entries: %v", out)
	}
}

func TestResolveSampler_LastMode(t *testing.T) {
	s := ResolveSampler(&SampleConfig{Mode: "last", N: 1})
	pairs := []string{"A=1", "B=2", "C=3"}
	out := s.Apply(pairs)
	if len(out) != 1 || out[0] != "C=3" {
		t.Errorf("expected [C=3], got %v", out)
	}
}

func TestResolveSampler_RandomMode_CorrectCount(t *testing.T) {
	s := ResolveSampler(&SampleConfig{Mode: "random", N: 2})
	pairs := []string{"A=1", "B=2", "C=3", "D=4"}
	out := s.Apply(pairs)
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
}

func TestResolveSampler_ReturnsCorrectType(t *testing.T) {
	s := ResolveSampler(&SampleConfig{Mode: "first", N: 1})
	if s == nil {
		t.Fatal("expected non-nil sampler")
	}
	_ = env.SampleModeFirst // ensure env package is used
}
