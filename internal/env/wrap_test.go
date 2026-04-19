package env

import (
	"strings"
	"testing"
)

func TestWrap_NilFunctions_Passthrough(t *testing.T) {
	w := NewWrapper(nil, nil)
	in := []string{"FOO=bar", "BAZ=qux"}
	out := w.Apply(in)
	if len(out) != 2 || out[0] != "FOO=bar" || out[1] != "BAZ=qux" {
		t.Fatalf("expected passthrough, got %v", out)
	}
}

func TestWrap_TransformsKey(t *testing.T) {
	w := NewWrapper(strings.ToLower, nil)
	out := w.Apply([]string{"FOO=bar"})
	if out[0] != "foo=bar" {
		t.Fatalf("expected foo=bar, got %s", out[0])
	}
}

func TestWrap_TransformsValue(t *testing.T) {
	w := NewWrapper(nil, strings.ToUpper)
	out := w.Apply([]string{"FOO=bar"})
	if out[0] != "FOO=BAR" {
		t.Fatalf("expected FOO=BAR, got %s", out[0])
	}
}

func TestWrap_TransformsBoth(t *testing.T) {
	w := NewWrapper(strings.ToLower, strings.ToUpper)
	out := w.Apply([]string{"Foo=bar"})
	if out[0] != "foo=BAR" {
		t.Fatalf("expected foo=BAR, got %s", out[0])
	}
}

func TestWrap_MalformedEntry_PassedThrough(t *testing.T) {
	w := NewWrapper(strings.ToLower, strings.ToUpper)
	out := w.Apply([]string{"NOEQUALS"})
	if out[0] != "NOEQUALS" {
		t.Fatalf("expected passthrough, got %s", out[0])
	}
}

func TestWrap_ValueContainsEquals(t *testing.T) {
	w := NewWrapper(nil, strings.ToUpper)
	out := w.Apply([]string{"KEY=val=ue"})
	if out[0] != "KEY=VAL=UE" {
		t.Fatalf("expected KEY=VAL=UE, got %s", out[0])
	}
}

func TestWrap_EmptySlice(t *testing.T) {
	w := NewWrapper(nil, nil)
	out := w.Apply([]string{})
	if len(out) != 0 {
		t.Fatalf("expected empty, got %v", out)
	}
}
