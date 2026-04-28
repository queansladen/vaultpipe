package env

import (
	"strings"
	"testing"
)

func TestHash_NoneMode_Passthrough(t *testing.T) {
	h := NewHasher(HashModeNone)
	input := []string{"KEY=secret", "OTHER=value"}
	out, err := h.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i, p := range input {
		if out[i] != p {
			t.Errorf("pair %d: got %q, want %q", i, out[i], p)
		}
	}
}

func TestHash_EmptyMode_DefaultsToNone(t *testing.T) {
	h := NewHasher("")
	if h.mode != HashModeNone {
		t.Errorf("expected HashModeNone, got %q", h.mode)
	}
}

func TestHash_UnknownMode_DefaultsToNone(t *testing.T) {
	h := NewHasher("blake2")
	if h.mode != HashModeNone {
		t.Errorf("expected HashModeNone, got %q", h.mode)
	}
}

func TestHash_MD5_ReplacesValue(t *testing.T) {
	h := NewHasher(HashModeMD5)
	out, err := h.Apply([]string{"TOKEN=hello"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// md5("hello") = 5d41402abc4b2a76b9719d911017c592
	want := "TOKEN=5d41402abc4b2a76b9719d911017c592"
	if out[0] != want {
		t.Errorf("got %q, want %q", out[0], want)
	}
}

func TestHash_SHA256_ReplacesValue(t *testing.T) {
	h := NewHasher(HashModeSHA256)
	out, err := h.Apply([]string{"SECRET=world"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// sha256("world") = 486ea46224d1bb4fb680f34f7c9ad96a8f24ec88be73ea8e5a6c65260e9cb8a7
	want := "SECRET=486ea46224d1bb4fb680f34f7c9ad96a8f24ec88be73ea8e5a6c65260e9cb8a7"
	if out[0] != want {
		t.Errorf("got %q, want %q", out[0], want)
	}
}

func TestHash_MalformedEntry_PassedThrough(t *testing.T) {
	h := NewHasher(HashModeSHA256)
	input := []string{"NOEQUALS"}
	out, err := h.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0] != "NOEQUALS" {
		t.Errorf("got %q, want %q", out[0], "NOEQUALS")
	}
}

func TestHash_MD5_EmptyValue(t *testing.T) {
	h := NewHasher(HashModeMD5)
	out, err := h.Apply([]string{"KEY="})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(out[0], "KEY=") || len(out[0]) == len("KEY=") {
		t.Errorf("expected non-empty hash for empty value, got %q", out[0])
	}
}
