package env

import (
	"testing"
)

func TestCoerce_None_Passthrough(t *testing.T) {
	c := NewCoercer(CoerceNone)
	input := []string{"FOO=Hello", "BAR=World"}
	out := c.Coerce(input)
	if out[0] != "FOO=Hello" || out[1] != "BAR=World" {
		t.Fatalf("expected passthrough, got %v", out)
	}
}

func TestCoerce_Upper(t *testing.T) {
	c := NewCoercer(CoerceUpper)
	out := c.Coerce([]string{"KEY=hello world", "X=mixedCase"})
	if out[0] != "KEY=HELLO WORLD" {
		t.Errorf("expected KEY=HELLO WORLD, got %s", out[0])
	}
	if out[1] != "X=MIXEDCASE" {
		t.Errorf("expected X=MIXEDCASE, got %s", out[1])
	}
}

func TestCoerce_Lower(t *testing.T) {
	c := NewCoercer(CoerceLower)
	out := c.Coerce([]string{"KEY=HELLO", "X=MixedCase"})
	if out[0] != "KEY=hello" {
		t.Errorf("expected KEY=hello, got %s", out[0])
	}
	if out[1] != "X=mixedcase" {
		t.Errorf("expected X=mixedcase, got %s", out[1])
	}
}

func TestCoerce_MalformedEntry_PassedThrough(t *testing.T) {
	c := NewCoercer(CoerceUpper)
	out := c.Coerce([]string{"NOEQUALS"})
	if len(out) != 1 || out[0] != "NOEQUALS" {
		t.Errorf("expected malformed entry to pass through, got %v", out)
	}
}

func TestCoerce_EmptySlice(t *testing.T) {
	c := NewCoercer(CoerceLower)
	out := c.Coerce([]string{})
	if len(out) != 0 {
		t.Errorf("expected empty slice, got %v", out)
	}
}

func TestCoerce_ValueContainsEquals(t *testing.T) {
	c := NewCoercer(CoerceLower)
	out := c.Coerce([]string{"KEY=VAL=UE"})
	if out[0] != "KEY=val=ue" {
		t.Errorf("expected KEY=val=ue, got %s", out[0])
	}
}
