package env

import (
	"strings"
	"testing"
)

func TestNewBuilderFromBase_EmptyBase(t *testing.T) {
	b := NewBuilderFromBase(nil)
	env := b.Build()
	if len(env) != 0 {
		t.Fatalf("expected empty env, got %d entries", len(env))
	}
}

func TestSet_InvalidKey_Empty(t *testing.T) {
	b := NewBuilderFromBase(nil)
	if err := b.Set("", "value"); err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestSet_InvalidKey_ContainsEquals(t *testing.T) {
	b := NewBuilderFromBase(nil)
	if err := b.Set("BAD=KEY", "value"); err == nil {
		t.Fatal("expected error for key containing '='")
	}
}

func TestSet_ValidKey(t *testing.T) {
	b := NewBuilderFromBase(nil)
	if err := b.Set("FOO", "bar"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	env := b.Build()
	if !containsEntry(env, "FOO=bar") {
		t.Fatalf("expected FOO=bar in env, got %v", env)
	}
}

func TestBuild_OverlayOverridesBase(t *testing.T) {
	base := []string{"HOME=/root", "PATH=/usr/bin"}
	b := NewBuilderFromBase(base)
	_ = b.Set("HOME", "/custom")

	env := b.Build()
	if containsEntry(env, "HOME=/root") {
		t.Fatal("base HOME should have been overridden")
	}
	if !containsEntry(env, "HOME=/custom") {
		t.Fatalf("expected HOME=/custom, got %v", env)
	}
}

func TestSetAll_SetsMultiple(t *testing.T) {
	b := NewBuilderFromBase(nil)
	secrets := map[string]string{"DB_PASS": "secret1", "API_KEY": "secret2"}
	if err := b.SetAll(secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	env := b.Build()
	if !containsEntry(env, "DB_PASS=secret1") {
		t.Errorf("missing DB_PASS in env")
	}
	if !containsEntry(env, "API_KEY=secret2") {
		t.Errorf("missing API_KEY in env")
	}
}

func TestKeys_ReturnsOverlayKeys(t *testing.T) {
	b := NewBuilderFromBase(nil)
	_ = b.Set("X", "1")
	_ = b.Set("Y", "2")
	keys := b.Keys()
	if len(keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(keys))
	}
}

func containsEntry(env []string, entry string) bool {
	for _, e := range env {
		if strings.EqualFold(e, entry) || e == entry {
			return true
		}
	}
	return false
}
