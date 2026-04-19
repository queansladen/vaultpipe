package env

import (
	"testing"
)

func TestJoin_EmptyBoth(t *testing.T) {
	j := NewJoiner(",")
	result := j.Join(nil, nil)
	if len(result) != 0 {
		t.Fatalf("expected empty, got %v", result)
	}
}

func TestJoin_BaseOnly(t *testing.T) {
	j := NewJoiner(",")
	result := j.Join([]string{"FOO=bar"}, nil)
	if len(result) != 1 || result[0] != "FOO=bar" {
		t.Fatalf("unexpected result: %v", result)
	}
}

func TestJoin_OverlayOnly(t *testing.T) {
	j := NewJoiner(",")
	result := j.Join(nil, []string{"FOO=bar"})
	if len(result) != 1 || result[0] != "FOO=bar" {
		t.Fatalf("unexpected result: %v", result)
	}
}

func TestJoin_NoDuplicates(t *testing.T) {
	j := NewJoiner(",")
	result := j.Join([]string{"A=1"}, []string{"B=2"})
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %v", result)
	}
}

func TestJoin_DuplicateKey_ConcatenatesValues(t *testing.T) {
	j := NewJoiner(",")
	result := j.Join([]string{"PATH=/usr/bin"}, []string{"PATH=/usr/local/bin"})
	if len(result) != 1 {
		t.Fatalf("expected 1 entry, got %v", result)
	}
	if result[0] != "PATH=/usr/bin,/usr/local/bin" {
		t.Errorf("unexpected value: %s", result[0])
	}
}

func TestJoin_CustomSeparator(t *testing.T) {
	j := NewJoiner(":")
	result := j.Join([]string{"PATH=/a"}, []string{"PATH=/b"})
	if result[0] != "PATH=/a:/b" {
		t.Errorf("unexpected value: %s", result[0])
	}
}

func TestJoin_DefaultSeparator_WhenEmpty(t *testing.T) {
	j := NewJoiner("")
	if j.separator != "," {
		t.Errorf("expected default separator ',', got %q", j.separator)
	}
}

func TestJoin_MalformedEntry_Skipped(t *testing.T) {
	j := NewJoiner(",")
	result := j.Join([]string{"MALFORMED", "GOOD=val"}, []string{"GOOD=extra"})
	if len(result) != 1 || result[0] != "GOOD=val,extra" {
		t.Errorf("unexpected result: %v", result)
	}
}
