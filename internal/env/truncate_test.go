package env

import (
	"strings"
	"testing"
)

func TestTruncate_NoLimit(t *testing.T) {
	tr := NewTruncator(0, "")
	if got := tr.Truncate("hello"); got != "hello" {
		t.Fatalf("expected 'hello', got %q", got)
	}
}

func TestTruncate_WithinLimit(t *testing.T) {
	tr := NewTruncator(10, "")
	if got := tr.Truncate("short"); got != "short" {
		t.Fatalf("expected 'short', got %q", got)
	}
}

func TestTruncate_ExceedsLimit(t *testing.T) {
	tr := NewTruncator(5, "")
	got := tr.Truncate("hello world")
	if got != "hello" {
		t.Fatalf("expected 'hello', got %q", got)
	}
}

func TestTruncate_WithSuffix(t *testing.T) {
	tr := NewTruncator(8, "...")
	got := tr.Truncate("hello world")
	if got != "hello..." {
		t.Fatalf("expected 'hello...', got %q", got)
	}
	if len(got) > 8 {
		t.Fatalf("result exceeds maxBytes: %d", len(got))
	}
}

func TestTruncate_SuffixLargerThanLimit(t *testing.T) {
	tr := NewTruncator(2, "...")
	got := tr.Truncate("hello")
	// cut becomes 0, suffix still appended
	if !strings.HasSuffix(got, "...") {
		t.Fatalf("expected suffix, got %q", got)
	}
}

func TestApply_TruncatesValues(t *testing.T) {
	tr := NewTruncator(4, "")
	pairs := []string{"FOO=abcdef", "BAR=xy", "MALFORMED"}
	out := tr.Apply(pairs)
	if out[0] != "FOO=abcd" {
		t.Fatalf("expected FOO=abcd, got %q", out[0])
	}
	if out[1] != "BAR=xy" {
		t.Fatalf("expected BAR=xy, got %q", out[1])
	}
	if out[2] != "MALFORMED" {
		t.Fatalf("expected MALFORMED passthrough, got %q", out[2])
	}
}

func TestApply_PreservesKeys(t *testing.T) {
	tr := NewTruncator(3, "")
	pairs := []string{"MY_KEY=longvalue"}
	out := tr.Apply(pairs)
	if !strings.HasPrefix(out[0], "MY_KEY=") {
		t.Fatalf("key mangled: %q", out[0])
	}
}
