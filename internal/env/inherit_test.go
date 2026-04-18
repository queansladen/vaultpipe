package env

import (
	"os"
	"strings"
	"testing"
)

func TestInherit_All(t *testing.T) {
	os.Setenv("_VP_TEST_VAR", "hello")
	t.Cleanup(func() { os.Unsetenv("_VP_TEST_VAR") })

	pairs := Inherit(InheritAll, nil)
	if len(pairs) == 0 {
		t.Fatal("expected non-empty environment")
	}
	found := false
	for _, p := range pairs {
		if strings.HasPrefix(p, "_VP_TEST_VAR=") {
			found = true
		}
	}
	if !found {
		t.Error("expected _VP_TEST_VAR in inherited environment")
	}
}

func TestInherit_None(t *testing.T) {
	pairs := Inherit(InheritNone, nil)
	if len(pairs) != 0 {
		t.Errorf("expected empty slice, got %d entries", len(pairs))
	}
}

func TestInherit_Filtered_NilFilter(t *testing.T) {
	pairs := Inherit(InheritFiltered, nil)
	if len(pairs) != 0 {
		t.Errorf("expected empty slice with nil filter, got %d", len(pairs))
	}
}

func TestInherit_Filtered_WithRules(t *testing.T) {
	os.Setenv("_VP_KEEP", "yes")
	os.Setenv("_VP_DROP", "no")
	t.Cleanup(func() {
		os.Unsetenv("_VP_KEEP")
		os.Unsetenv("_VP_DROP")
	})

	f := NewFilter([]string{"_VP_KEEP"})
	pairs := Inherit(InheritFiltered, f)

	for _, p := range pairs {
		if strings.HasPrefix(p, "_VP_DROP=") {
			t.Error("_VP_DROP should have been filtered out")
		}
	}
}

func TestMergeEnv_OverlayWins(t *testing.T) {
	base := []string{"FOO=base", "BAR=keep"}
	overlay := []string{"FOO=overlay", "BAZ=new"}

	result := MergeEnv(base, overlay)

	got := make(map[string]string)
	for _, p := range result {
		k, v, _ := strings.Cut(p, "=")
		got[k] = v
	}

	if got["FOO"] != "overlay" {
		t.Errorf("expected FOO=overlay, got %s", got["FOO"])
	}
	if got["BAR"] != "keep" {
		t.Errorf("expected BAR=keep, got %s", got["BAR"])
	}
	if got["BAZ"] != "new" {
		t.Errorf("expected BAZ=new, got %s", got["BAZ"])
	}
}

func TestMergeEnv_NoDuplicateKeys(t *testing.T) {
	base := []string{"X=1"}
	overlay := []string{"X=2"}
	result := MergeEnv(base, overlay)
	count := 0
	for _, p := range result {
		if strings.HasPrefix(p, "X=") {
			count++
		}
	}
	if count != 1 {
		t.Errorf("expected exactly 1 entry for X, got %d", count)
	}
}
