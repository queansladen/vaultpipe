package env

import (
	"sort"
	"testing"
)

func TestNewSnapshot_ParsesPairs(t *testing.T) {
	s := NewSnapshot([]string{"FOO=bar", "BAZ=qux"})
	if s.Len() != 2 {
		t.Fatalf("expected 2 entries, got %d", s.Len())
	}
}

func TestNewSnapshot_SkipsMalformed(t *testing.T) {
	s := NewSnapshot([]string{"NOEQUALS", "GOOD=val"})
	if s.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", s.Len())
	}
}

func TestGet_Present(t *testing.T) {
	s := NewSnapshot([]string{"KEY=secret"})
	v, ok := s.Get("KEY")
	if !ok {
		t.Fatal("expected key to be present")
	}
	if v != "secret" {
		t.Fatalf("expected 'secret', got %q", v)
	}
}

func TestGet_Missing(t *testing.T) {
	s := NewSnapshot([]string{})
	_, ok := s.Get("MISSING")
	if ok {
		t.Fatal("expected key to be absent")
	}
}

func TestNewSnapshot_ValueContainsEquals(t *testing.T) {
	s := NewSnapshot([]string{"URL=http://x.com?a=1&b=2"})
	v, ok := s.Get("URL")
	if !ok {
		t.Fatal("expected URL key")
	}
	if v != "http://x.com?a=1&b=2" {
		t.Fatalf("unexpected value: %q", v)
	}
}

func TestKeys_ReturnsAll(t *testing.T) {
	s := NewSnapshot([]string{"A=1", "B=2", "C=3"})
	keys := s.Keys()
	sort.Strings(keys)
	expected := []string{"A", "B", "C"}
	for i, k := range expected {
		if keys[i] != k {
			t.Fatalf("expected %q at index %d, got %q", k, i, keys[i])
		}
	}
}

func TestPairs_RoundTrip(t *testing.T) {
	input := []string{"FOO=bar", "BAZ=qux"}
	s := NewSnapshot(input)
	pairs := s.Pairs()
	if len(pairs) != len(input) {
		t.Fatalf("expected %d pairs, got %d", len(input), len(pairs))
	}
}
