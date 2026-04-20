package env

import (
	"testing"
)

func TestShell_DefaultFormat_Export(t *testing.T) {
	s := NewSheller("", false)
	got := s.Apply([]string{"FOO=bar"})
	if len(got) != 1 || got[0] != "export FOO=bar" {
		t.Fatalf("expected 'export FOO=bar', got %v", got)
	}
}

func TestShell_ExportFormat_NoQuote(t *testing.T) {
	s := NewSheller(ShellFormatExport, false)
	got := s.Apply([]string{"KEY=hello world"})
	if got[0] != "export KEY=hello world" {
		t.Fatalf("unexpected: %q", got[0])
	}
}

func TestShell_ExportFormat_WithQuote(t *testing.T) {
	s := NewSheller(ShellFormatExport, true)
	got := s.Apply([]string{"KEY=hello world"})
	if got[0] != "export KEY='hello world'" {
		t.Fatalf("unexpected: %q", got[0])
	}
}

func TestShell_InlineFormat(t *testing.T) {
	s := NewSheller(ShellFormatInline, false)
	got := s.Apply([]string{"A=1", "B=2"})
	if got[0] != "A=1" || got[1] != "B=2" {
		t.Fatalf("unexpected: %v", got)
	}
}

func TestShell_UnsetFormat(t *testing.T) {
	s := NewSheller(ShellFormatUnset, false)
	got := s.Apply([]string{"SECRET=value"})
	if got[0] != "unset SECRET" {
		t.Fatalf("unexpected: %q", got[0])
	}
}

func TestShell_MalformedEntry_PassedThrough(t *testing.T) {
	s := NewSheller(ShellFormatExport, false)
	got := s.Apply([]string{"NOEQUALS"})
	if got[0] != "NOEQUALS" {
		t.Fatalf("expected passthrough, got %q", got[0])
	}
}

func TestShell_QuoteEscapesSingleQuote(t *testing.T) {
	s := NewSheller(ShellFormatExport, true)
	got := s.Apply([]string{"KEY=it's here"})
	// expected: export KEY='it'\''s here'
	want := "export KEY='it'\\''s here'"
	if got[0] != want {
		t.Fatalf("expected %q, got %q", want, got[0])
	}
}

func TestShell_EmptySlice(t *testing.T) {
	s := NewSheller(ShellFormatExport, false)
	got := s.Apply([]string{})
	if len(got) != 0 {
		t.Fatalf("expected empty, got %v", got)
	}
}

func TestShell_ValueContainsEquals(t *testing.T) {
	s := NewSheller(ShellFormatInline, false)
	got := s.Apply([]string{"URL=http://x.com?a=1"})
	if got[0] != "URL=http://x.com?a=1" {
		t.Fatalf("unexpected: %q", got[0])
	}
}
