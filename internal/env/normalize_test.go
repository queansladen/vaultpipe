package env

import (
	"testing"
)

func TestNormalize_NoneMode_Passthrough(t *testing.T) {
	n := NewNormalizer(NormalizeModeNone)
	input := []string{"MY-KEY=value", "another.key=123"}
	out, err := n.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i, p := range out {
		if p != input[i] {
			t.Errorf("pair %d: got %q, want %q", i, p, input[i])
		}
	}
}

func TestNormalize_EmptyMode_DefaultsToNone(t *testing.T) {
	n := NewNormalizer("")
	input := []string{"MY-KEY=val"}
	out, _ := n.Apply(input)
	if out[0] != "MY-KEY=val" {
		t.Errorf("got %q, want MY-KEY=val", out[0])
	}
}

func TestNormalize_SnakeMode_ReplacesHyphen(t *testing.T) {
	n := NewNormalizer(NormalizeModeSnake)
	out, err := n.Apply([]string{"MY-KEY=value"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0] != "MY_KEY=value" {
		t.Errorf("got %q, want MY_KEY=value", out[0])
	}
}

func TestNormalize_SnakeMode_CollapseRunOfSeparators(t *testing.T) {
	n := NewNormalizer(NormalizeModeSnake)
	out, _ := n.Apply([]string{"MY--KEY=val"})
	if out[0] != "MY_KEY=val" {
		t.Errorf("got %q, want MY_KEY=val", out[0])
	}
}

func TestNormalize_DotMode_ReplacesSeparator(t *testing.T) {
	n := NewNormalizer(NormalizeModeDot)
	out, _ := n.Apply([]string{"MY_KEY=val"})
	if out[0] != "MY.KEY=val" {
		t.Errorf("got %q, want MY.KEY=val", out[0])
	}
}

func TestNormalize_MalformedEntry_PassedThrough(t *testing.T) {
	n := NewNormalizer(NormalizeModeSnake)
	input := []string{"NOKEYVALUE"}
	out, _ := n.Apply(input)
	if out[0] != "NOKEYVALUE" {
		t.Errorf("got %q, want NOKEYVALUE", out[0])
	}
}

func TestNormalize_ValueContainsEquals_KeyOnly(t *testing.T) {
	n := NewNormalizer(NormalizeModeSnake)
	out, _ := n.Apply([]string{"MY-KEY=a=b"})
	if out[0] != "MY_KEY=a=b" {
		t.Errorf("got %q, want MY_KEY=a=b", out[0])
	}
}

func TestNormalize_LeadingTrailingSeparator_Trimmed(t *testing.T) {
	n := NewNormalizer(NormalizeModeSnake)
	out, _ := n.Apply([]string{"-KEY-=val"})
	if out[0] != "KEY=val" {
		t.Errorf("got %q, want KEY=val", out[0])
	}
}
