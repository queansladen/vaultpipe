package env

import (
	"os"
	"testing"
)

func TestInterp_NoVariables_Passthrough(t *testing.T) {
	i := NewInterpolator(nil, false)
	in := []string{"FOO=bar", "BAZ=qux"}
	out, err := i.Apply(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 || out[0] != "FOO=bar" || out[1] != "BAZ=qux" {
		t.Fatalf("expected passthrough, got %v", out)
	}
}

func TestInterp_ExpandsFromOverlay(t *testing.T) {
	i := NewInterpolator(map[string]string{"HOST": "localhost"}, false)
	out, err := i.Apply([]string{"ADDR=${HOST}:5432"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0] != "ADDR=localhost:5432" {
		t.Fatalf("got %q", out[0])
	}
}

func TestInterp_ExpandsFromPreviousPair(t *testing.T) {
	i := NewInterpolator(nil, false)
	out, err := i.Apply([]string{"BASE=/opt/app", "LOG=${BASE}/logs"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[1] != "LOG=/opt/app/logs" {
		t.Fatalf("got %q", out[1])
	}
}

func TestInterp_UnknownVar_NoFallback_EmptyString(t *testing.T) {
	i := NewInterpolator(nil, false)
	out, err := i.Apply([]string{"X=${UNDEFINED}"}	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0] != "X=" {
		t.Fatalf("got %q", out[0])
	}
}

func TestInterp_FallbackToOS(t *testing.T) {
	t.Setenv("OS_VAR", "from-os")
	_ = os.Getenv("OS_VAR") // ensure set
	i := NewInterpolator(nil, true)
	out, err := i.Apply([]string{"Y=$OS_VAR"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0] != "Y=from-os" {
		t.Fatalf("got %q", out[0])
	}
}

func TestInterp_OverlayTakesPrecedenceOverOS(t *testing.T) {
	t.Setenv("PRIO", "from-os")
	i := NewInterpolator(map[string]string{"PRIO": "from-overlay"}, true)
	out, err := i.Apply([]string{"Z=${PRIO}"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0] != "Z=from-overlay" {
		t.Fatalf("got %q", out[0])
	}
}

func TestInterp_MalformedEntry_PassedThrough(t *testing.T) {
	i := NewInterpolator(nil, false)
	out, err := i.Apply([]string{"NOEQUALS"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0] != "NOEQUALS" {
		t.Fatalf("got %q", out[0])
	}
}
