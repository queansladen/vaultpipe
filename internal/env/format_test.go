package env

import (
	"testing"
)

func TestFormat_RawMode_Passthrough(t *testing.T) {
	f := NewFormatter(FormatModeRaw)
	in := []string{"FOO=bar", "BAZ=qux"}
	out := f.Apply(in)
	if out[0] != "FOO=bar" || out[1] != "BAZ=qux" {
		t.Fatalf("expected passthrough, got %v", out)
	}
}

func TestFormat_EmptyMode_DefaultsToRaw(t *testing.T) {
	f := NewFormatter("")
	in := []string{"X=1"}
	out := f.Apply(in)
	if out[0] != "X=1" {
		t.Fatalf("expected raw passthrough, got %s", out[0])
	}
}

func TestFormat_ExportMode(t *testing.T) {
	f := NewFormatter(FormatModeExport)
	out := f.Apply([]string{"FOO=hello world"})
	if out[0] != "export FOO=hello world" {
		t.Fatalf("unexpected: %s", out[0])
	}
}

func TestFormat_DotenvMode_QuotesValue(t *testing.T) {
	f := NewFormatter(FormatModeDotenv)
	out := f.Apply([]string{"FOO=bar"})
	if out[0] != `FOO="bar"` {
		t.Fatalf("unexpected: %s", out[0])
	}
}

func TestFormat_DotenvMode_EscapesDoubleQuotes(t *testing.T) {
	f := NewFormatter(FormatModeDotenv)
	out := f.Apply([]string{`MSG=say "hi"`})
	expected := `MSG="say \"hi\""`
	if out[0] != expected {
		t.Fatalf("expected %s, got %s", expected, out[0])
	}
}

func TestFormat_MalformedEntry_PassedThrough(t *testing.T) {
	f := NewFormatter(FormatModeExport)
	out := f.Apply([]string{"NOEQUALS"})
	if out[0] != "NOEQUALS" {
		t.Fatalf("expected passthrough, got %s", out[0])
	}
}

func TestFormat_EmptySlice(t *testing.T) {
	f := NewFormatter(FormatModeExport)
	out := f.Apply([]string{})
	if len(out) != 0 {
		t.Fatalf("expected empty, got %v", out)
	}
}
