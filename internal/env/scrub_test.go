package env

import (
	"testing"
)

func TestScrub_NoneMode_Passthrough(t *testing.T) {
	s := NewScrubber(ScrubNone)
	in := []string{"A=1", "B=", "C=3"}
	out := s.Apply(in)
	if len(out) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(out))
	}
}

func TestScrub_EmptyMode_RemovesEmptyValues(t *testing.T) {
	s := NewScrubber(ScrubEmpty)
	in := []string{"A=1", "B=", "C=", "D=4"}
	out := s.Apply(in)
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	if out[0] != "A=1" || out[1] != "D=4" {
		t.Errorf("unexpected entries: %v", out)
	}
}

func TestScrub_EmptyMode_KeepsNonEmpty(t *testing.T) {
	s := NewScrubber(ScrubEmpty)
	in := []string{"X=hello", "Y=world"}
	out := s.Apply(in)
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
}

func TestScrub_PrefixMode_RemovesMatchingKeys(t *testing.T) {
	s := NewScrubber(ScrubPrefix, "TMP_", "DEBUG_")
	in := []string{"TMP_FOO=1", "DEBUG_LEVEL=2", "APP_KEY=3", "DEBUG_=4"}
	out := s.Apply(in)
	if len(out) != 1 {
		t.Fatalf("expected 1 entry, got %d: %v", len(out), out)
	}
	if out[0] != "APP_KEY=3" {
		t.Errorf("unexpected entry: %s", out[0])
	}
}

func TestScrub_PrefixMode_NoPrefixes_KeepsAll(t *testing.T) {
	s := NewScrubber(ScrubPrefix)
	in := []string{"A=1", "B=2"}
	out := s.Apply(in)
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
}

func TestScrub_MalformedEntry_PassedThrough(t *testing.T) {
	s := NewScrubber(ScrubEmpty)
	in := []string{"NOEQUALS", "A=1"}
	out := s.Apply(in)
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
}

func TestScrub_EmptyMode_DefaultsWhenBlank(t *testing.T) {
	s := NewScrubber("")
	if s.mode != ScrubNone {
		t.Errorf("expected ScrubNone default, got %q", s.mode)
	}
}
