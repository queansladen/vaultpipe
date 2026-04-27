package env

import (
	"testing"
)

func TestTypecast_NoneMode_Passthrough(t *testing.T) {
	tc := NewTypecaster(CastNone)
	in := []string{"FOO=TRUE", "BAR=1.5", "BAZ=hello"}
	out, err := tc.Apply(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i, want := range in {
		if out[i] != want {
			t.Errorf("pair %d: got %q, want %q", i, out[i], want)
		}
	}
}

func TestTypecast_EmptyMode_DefaultsToNone(t *testing.T) {
	tc := NewTypecaster("")
	if tc.mode != CastNone {
		t.Fatalf("expected CastNone, got %q", tc.mode)
	}
}

func TestTypecast_BoolMode_NormalisesTrueVariants(t *testing.T) {
	tc := NewTypecaster(CastBool)
	cases := []struct{ in, want string }{
		{"ENABLED=TRUE", "ENABLED=true"},
		{"ENABLED=1", "ENABLED=true"},
		{"ENABLED=False", "ENABLED=false"},
		{"ENABLED=0", "ENABLED=false"},
	}
	for _, c := range cases {
		out, err := tc.Apply([]string{c.in})
		if err != nil {
			t.Fatalf("input %q: unexpected error: %v", c.in, err)
		}
		if out[0] != c.want {
			t.Errorf("input %q: got %q, want %q", c.in, out[0], c.want)
		}
	}
}

func TestTypecast_BoolMode_NonBool_Passthrough(t *testing.T) {
	tc := NewTypecaster(CastBool)
	in := []string{"FOO=hello"}
	out, err := tc.Apply(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0] != "FOO=hello" {
		t.Errorf("got %q, want %q", out[0], "FOO=hello")
	}
}

func TestTypecast_NumericMode_NormalisesIntegers(t *testing.T) {
	tc := NewTypecaster(CastNumeric)
	in := []string{"COUNT=007", "RATIO=3.14000"}
	out, err := tc.Apply(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0] != "COUNT=7" {
		t.Errorf("integer: got %q, want %q", out[0], "COUNT=7")
	}
	if out[1] != "RATIO=3.14" {
		t.Errorf("float: got %q, want %q", out[1], "RATIO=3.14")
	}
}

func TestTypecast_AutoMode_InfersType(t *testing.T) {
	tc := NewTypecaster(CastAuto)
	cases := []struct{ in, want string }{
		{"A=TRUE", "A=true"},
		{"B=042", "B=42"},
		{"C=2.500", "C=2.5"},
		{"D=hello", "D=hello"},
	}
	for _, c := range cases {
		out, err := tc.Apply([]string{c.in})
		if err != nil {
			t.Fatalf("input %q: unexpected error: %v", c.in, err)
		}
		if out[0] != c.want {
			t.Errorf("input %q: got %q, want %q", c.in, out[0], c.want)
		}
	}
}

func TestTypecast_MalformedEntry_PassedThrough(t *testing.T) {
	tc := NewTypecaster(CastAuto)
	in := []string{"NOEQUALS"}
	out, err := tc.Apply(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0] != "NOEQUALS" {
		t.Errorf("got %q, want %q", out[0], "NOEQUALS")
	}
}
