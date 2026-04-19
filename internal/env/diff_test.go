package env

import (
	"testing"
)

func findChange(changes []Change, key string) (Change, bool) {
	for _, c := range changes {
		if c.Key == key {
			return c, true
		}
	}
	return Change{}, false
}

func TestDiff_NoChanges(t *testing.T) {
	env := []string{"FOO=bar", "BAZ=qux"}
	changes := Diff(env, env)
	if len(changes) != 0 {
		t.Fatalf("expected 0 changes, got %d", len(changes))
	}
}

func TestDiff_Added(t *testing.T) {
	before := []string{"FOO=bar"}
	after := []string{"FOO=bar", "NEW=value"}
	changes := Diff(before, after)
	c, ok := findChange(changes, "NEW")
	if !ok {
		t.Fatal("expected change for NEW")
	}
	if c.Change != ChangeAdded || c.New != "value" {
		t.Errorf("unexpected change: %+v", c)
	}
}

func TestDiff_Removed(t *testing.T) {
	before := []string{"FOO=bar", "OLD=gone"}
	after := []string{"FOO=bar"}
	changes := Diff(before, after)
	c, ok := findChange(changes, "OLD")
	if !ok {
		t.Fatal("expected change for OLD")
	}
	if c.Change != ChangeRemoved || c.Old != "gone" {
		t.Errorf("unexpected change: %+v", c)
	}
}

func TestDiff_Updated(t *testing.T) {
	before := []string{"FOO=old"}
	after := []string{"FOO=new"}
	changes := Diff(before, after)
	c, ok := findChange(changes, "FOO")
	if !ok {
		t.Fatal("expected change for FOO")
	}
	if c.Change != ChangeUpdated || c.Old != "old" || c.New != "new" {
		t.Errorf("unexpected change: %+v", c)
	}
}

func TestDiff_ValueContainsEquals(t *testing.T) {
	before := []string{"TOKEN=abc=def"}
	after := []string{"TOKEN=abc=xyz"}
	changes := Diff(before, after)
	c, ok := findChange(changes, "TOKEN")
	if !ok {
		t.Fatal("expected change for TOKEN")
	}
	if c.Old != "abc=def" || c.New != "abc=xyz" {
		t.Errorf("unexpected values: %+v", c)
	}
}

func TestDiff_Empty(t *testing.T) {
	changes := Diff(nil, nil)
	if len(changes) != 0 {
		t.Fatalf("expected 0 changes, got %d", len(changes))
	}
}
