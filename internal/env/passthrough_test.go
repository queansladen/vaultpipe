package env

import (
	"testing"
)

func makeSnapshot(pairs []string) *Snapshot {
	return NewSnapshot(pairs)
}

func TestPassthrough_All_ForwardsHost(t *testing.T) {
	host := makeSnapshot([]string{"HOME=/root", "PATH=/usr/bin"})
	result := Passthrough(PassthroughAll, host, nil, nil)
	snap := NewSnapshot(result)
	if v, ok := snap.Get("HOME"); !ok || v != "/root" {
		t.Errorf("expected HOME=/root, got %q", v)
	}
}

func TestPassthrough_None_EmptyBase(t *testing.T) {
	host := makeSnapshot([]string{"HOME=/root", "SECRET=hunter2"})
	result := Passthrough(PassthroughNone, host, nil, nil)
	if len(result) != 0 {
		t.Errorf("expected empty env, got %v", result)
	}
}

func TestPassthrough_None_OverlayStillPresent(t *testing.T) {
	host := makeSnapshot([]string{"HOME=/root"})
	overlay := map[string]string{"TOKEN": "abc123"}
	result := Passthrough(PassthroughNone, host, nil, overlay)
	snap := NewSnapshot(result)
	if v, ok := snap.Get("TOKEN"); !ok || v != "abc123" {
		t.Errorf("expected TOKEN=abc123, got %q", v)
	}
	if _, ok := snap.Get("HOME"); ok {
		t.Error("expected HOME to be absent")
	}
}

func TestPassthrough_Filtered_NilFilter_AllowsAll(t *testing.T) {
	host := makeSnapshot([]string{"HOME=/root", "PATH=/bin"})
	result := Passthrough(PassthroughFiltered, host, nil, nil)
	snap := NewSnapshot(result)
	if _, ok := snap.Get("HOME"); !ok {
		t.Error("expected HOME to be present")
	}
}

func TestPassthrough_Filtered_WithRules(t *testing.T) {
	host := makeSnapshot([]string{"HOME=/root", "PATH=/bin", "MY_VAR=yes"})
	f := NewFilter([]string{"MY_VAR"})
	result := Passthrough(PassthroughFiltered, host, f, nil)
	snap := NewSnapshot(result)
	if _, ok := snap.Get("HOME"); ok {
		t.Error("HOME should be filtered out")
	}
	if v, ok := snap.Get("MY_VAR"); !ok || v != "yes" {
		t.Errorf("expected MY_VAR=yes, got %q", v)
	}
}

func TestPassthrough_OverlayWinsOverHost(t *testing.T) {
	host := makeSnapshot([]string{"TOKEN=old"})
	overlay := map[string]string{"TOKEN": "new"}
	result := Passthrough(PassthroughAll, host, nil, overlay)
	snap := NewSnapshot(result)
	if v, _ := snap.Get("TOKEN"); v != "new" {
		t.Errorf("expected overlay to win, got %q", v)
	}
}
