package env

import (
	"testing"
)

func TestQuote_NoShellSafe_Passthrough(t *testing.T) {
	q := NewQuoter(false)
	in := []string{"FOO=hello world", "BAR=plain"}
	out := q.Quote(in)
	if out[0] != "FOO=hello world" {
		t.Fatalf("expected unchanged, got %q", out[0])
	}
}

func TestQuote_ShellSafe_PlainValue_Unchanged(t *testing.T) {
	q := NewQuoter(true)
	out := q.Quote([]string{"FOO=plainvalue"})
	if out[0] != "FOO=plainvalue" {
		t.Fatalf("unexpected quoting: %q", out[0])
	}
}

func TestQuote_ShellSafe_SpaceInValue(t *testing.T) {
	q := NewQuoter(true)
	out := q.Quote([]string{"FOO=hello world"})
	if out[0] != "FOO='hello world'" {
		t.Fatalf("expected single-quoted, got %q", out[0])
	}
}

func TestQuote_ShellSafe_EmbeddedSingleQuote(t *testing.T) {
	q := NewQuoter(true)
	out := q.Quote([]string{"FOO=it's"})
	want := "FOO='it'\\''s'"
	if out[0] != want {
		t.Fatalf("want %q got %q", want, out[0])
	}
}

func TestQuote_MalformedEntry_PassedThrough(t *testing.T) {
	q := NewQuoter(true)
	out := q.Quote([]string{"NOEQUALS"})
	if out[0] != "NOEQUALS" {
		t.Fatalf("expected passthrough, got %q", out[0])
	}
}

func TestQuote_EmptyValue(t *testing.T) {
	q := NewQuoter(true)
	out := q.Quote([]string{"FOO="})
	if out[0] != "FOO=" {
		t.Fatalf("expected empty value unchanged, got %q", out[0])
	}
}

func TestQuote_ValueContainsEquals(t *testing.T) {
	q := NewQuoter(true)
	out := q.Quote([]string{"FOO=a=b"})
	// '=' is not a safe char so value should be quoted
	if out[0] != "FOO='a=b'" {
		t.Fatalf("expected quoted, got %q", out[0])
	}
}
