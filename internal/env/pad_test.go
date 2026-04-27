package env

import (
	"testing"
)

func TestPad_NoneMode_Passthrough(t *testing.T) {
	p := NewPadder(PadModeNone, 10, ' ')
	in := []string{"FOO=bar", "BAZ=qux"}
	out, err := p.Apply(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i, v := range in {
		if out[i] != v {
			t.Errorf("[%d] got %q, want %q", i, out[i], v)
		}
	}
}

func TestPad_ZeroWidth_Passthrough(t *testing.T) {
	p := NewPadder(PadModeRight, 0, ' ')
	in := []string{"FOO=hi"}
	out, err := p.Apply(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0] != "FOO=hi" {
		t.Errorf("got %q, want %q", out[0], "FOO=hi")
	}
}

func TestPad_RightPad_PadsShortValue(t *testing.T) {
	p := NewPadder(PadModeRight, 8, '-')
	out, err := p.Apply([]string{"KEY=abc"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "KEY=abc-----"
	if out[0] != want {
		t.Errorf("got %q, want %q", out[0], want)
	}
}

func TestPad_LeftPad_PadsShortValue(t *testing.T) {
	p := NewPadder(PadModeLeft, 6, '0')
	out, err := p.Apply([]string{"NUM=42"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "NUM=000042"
	if out[0] != want {
		t.Errorf("got %q, want %q", out[0], want)
	}
}

func TestPad_ValueAlreadyWide_Passthrough(t *testing.T) {
	p := NewPadder(PadModeRight, 3, ' ')
	out, err := p.Apply([]string{"KEY=longvalue"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0] != "KEY=longvalue" {
		t.Errorf("got %q, want %q", out[0], "KEY=longvalue")
	}
}

func TestPad_MalformedEntry_PassedThrough(t *testing.T) {
	p := NewPadder(PadModeLeft, 10, ' ')
	out, err := p.Apply([]string{"NOEQUALS"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0] != "NOEQUALS" {
		t.Errorf("got %q, want %q", out[0], "NOEQUALS")
	}
}

func TestPad_DefaultChar_Space(t *testing.T) {
	p := NewPadder(PadModeRight, 5, 0)
	out, err := p.Apply([]string{"K=ab"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "K=ab   "
	if out[0] != want {
		t.Errorf("got %q, want %q", out[0], want)
	}
}
