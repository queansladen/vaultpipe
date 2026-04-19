package env

import (
	"sort"
	"testing"
)

func TestNewCheckpoint_ParsesPairs(t *testing.T) {
	cp := NewCheckpoint("before", []string{"FOO=bar", "BAZ=qux"})
	if cp.Len() != 2 {
		t.Fatalf("expected 2, got %d", cp.Len())
	}
	if cp.Name() != "before" {
		t.Fatalf("unexpected name %s", cp.Name())
	}
}

func TestNewCheckpoint_SkipsMalformed(t *testing.T) {
	cp := NewCheckpoint("x", []string{"NOEQUALSSIGN", "GOOD=val"})
	if cp.Len() != 1 {
		t.Fatalf("expected 1, got %d", cp.Len())
	}
}

func TestDiff_NoChanges(t *testing.T) {
	a := NewCheckpoint("a", []string{"K=v"})
	b := NewCheckpoint("b", []string{"K=v"})
	if d := a.Diff(b); len(d) != 0 {
		t.Fatalf("expected no diff, got %v", d)
	}
}

func TestDiff_ValueChanged(t *testing.T) {
	a := NewCheckpoint("a", []string{"K=old"})
	b := NewCheckpoint("b", []string{"K=new"})
	d := a.Diff(b)
	if len(d) != 1 {
		t.Fatalf("expected 1 diff entry, got %v", d)
	}
	if d[0] != "K=old->new" {
		t.Fatalf("unexpected diff: %s", d[0])
	}
}

func TestDiff_KeyAdded(t *testing.T) {
	a := NewCheckpoint("a", nil)
	b := NewCheckpoint("b", []string{"NEW=val"})
	d := a.Diff(b)
	if len(d) != 1 || d[0] != "NEW=<unset>->val" {
		t.Fatalf("unexpected diff: %v", d)
	}
}

func TestDiff_KeyRemoved(t *testing.T) {
	a := NewCheckpoint("a", []string{"GONE=x"})
	b := NewCheckpoint("b", nil)
	d := a.Diff(b)
	if len(d) != 1 || d[0] != "GONE=x-><unset>" {
		t.Fatalf("unexpected diff: %v", d)
	}
}

func TestDiff_MultipleChanges(t *testing.T) {
	a := NewCheckpoint("a", []string{"A=1", "B=2"})
	b := NewCheckpoint("b", []string{"A=9", "C=3"})
	d := a.Diff(b)
	sort.Strings(d)
	if len(d) != 3 {
		t.Fatalf("expected 3 diffs, got %v", d)
	}
}
