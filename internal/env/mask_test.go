package env

import (
	"testing"
)

func TestMask_NoneMode_Passthrough(t *testing.T) {
	m := NewMasker(MaskNone, "*", 0)
	in := []string{"SECRET=hunter2", "OTHER=value"}
	out := m.Apply(in)
	for i, v := range in {
		if out[i] != v {
			t.Errorf("expected %q, got %q", v, out[i])
		}
	}
}

func TestMask_FullMode_MasksEntireValue(t *testing.T) {
	m := NewMasker(MaskFull, "*", 0)
	out := m.Apply([]string{"TOKEN=abc123"})
	if out[0] != "TOKEN=******" {
		t.Errorf("unexpected: %s", out[0])
	}
}

func TestMask_PartialMode_KeepsTrailing(t *testing.T) {
	m := NewMasker(MaskPartial, "*", 3)
	out := m.Apply([]string{"TOKEN=abcdef"})
	if out[0] != "TOKEN=***def" {
		t.Errorf("unexpected: %s", out[0])
	}
}

func TestMask_PartialMode_ShortValue_FullyMasked(t *testing.T) {
	m := NewMasker(MaskPartial, "*", 4)
	out := m.Apply([]string{"K=ab"})
	if out[0] != "K=**" {
		t.Errorf("unexpected: %s", out[0])
	}
}

func TestMask_MalformedEntry_PassedThrough(t *testing.T) {
	m := NewMasker(MaskFull, "*", 0)
	out := m.Apply([]string{"NOEQUALS"})
	if out[0] != "NOEQUALS" {
		t.Errorf("unexpected: %s", out[0])
	}
}

func TestMask_CustomChar(t *testing.T) {
	m := NewMasker(MaskFull, "#", 0)
	out := m.Apply([]string{"K=val"})
	if out[0] != "K=###" {
		t.Errorf("unexpected: %s", out[0])
	}
}

func TestMask_EmptyChar_DefaultsStar(t *testing.T) {
	m := NewMasker(MaskFull, "", 0)
	out := m.Apply([]string{"K=ab"})
	if out[0] != "K=**" {
		t.Errorf("unexpected: %s", out[0])
	}
}
