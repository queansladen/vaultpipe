package env

import (
	"math/rand"
	"testing"
)

var samplePairs = []string{
	"ALPHA=1",
	"BETA=2",
	"GAMMA=3",
	"DELTA=4",
	"EPSILON=5",
}

func TestSample_NoneMode_Passthrough(t *testing.T) {
	s := NewSampler(SampleModeNone, 2, nil)
	out := s.Apply(samplePairs)
	if len(out) != len(samplePairs) {
		t.Fatalf("expected %d entries, got %d", len(samplePairs), len(out))
	}
}

func TestSample_ZeroN_Passthrough(t *testing.T) {
	s := NewSampler(SampleModeFirst, 0, nil)
	out := s.Apply(samplePairs)
	if len(out) != len(samplePairs) {
		t.Fatalf("expected passthrough when n=0, got %d entries", len(out))
	}
}

func TestSample_First_ReturnsHead(t *testing.T) {
	s := NewSampler(SampleModeFirst, 2, nil)
	out := s.Apply(samplePairs)
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	if out[0] != "ALPHA=1" || out[1] != "BETA=2" {
		t.Errorf("unexpected entries: %v", out)
	}
}

func TestSample_Last_ReturnsTail(t *testing.T) {
	s := NewSampler(SampleModeLast, 2, nil)
	out := s.Apply(samplePairs)
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	if out[0] != "DELTA=4" || out[1] != "EPSILON=5" {
		t.Errorf("unexpected entries: %v", out)
	}
}

func TestSample_Random_CorrectCount(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	s := NewSampler(SampleModeRandom, 3, rng)
	out := s.Apply(samplePairs)
	if len(out) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(out))
	}
}

func TestSample_NGreaterThanLen_CapsAtLen(t *testing.T) {
	s := NewSampler(SampleModeFirst, 100, nil)
	out := s.Apply(samplePairs)
	if len(out) != len(samplePairs) {
		t.Fatalf("expected %d entries, got %d", len(samplePairs), len(out))
	}
}

func TestSample_DoesNotMutateInput(t *testing.T) {
	rng := rand.New(rand.NewSource(7))
	s := NewSampler(SampleModeRandom, 3, rng)
	original := make([]string, len(samplePairs))
	copy(original, samplePairs)
	s.Apply(samplePairs)
	for i, v := range samplePairs {
		if v != original[i] {
			t.Errorf("input mutated at index %d: got %q want %q", i, v, original[i])
		}
	}
}

func TestSample_Keys_ReturnsKeySubset(t *testing.T) {
	s := NewSampler(SampleModeFirst, 2, nil)
	keys := s.Keys(samplePairs)
	if len(keys) != 2 || keys[0] != "ALPHA" || keys[1] != "BETA" {
		t.Errorf("unexpected keys: %v", keys)
	}
}
