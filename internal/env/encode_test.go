package env

import (
	"encoding/base64"
	"testing"
)

func b64(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func TestEncode_NoneMode_Passthrough(t *testing.T) {
	e := NewEncoder(EncodeModeNone)
	input := []string{"FOO=bar", "BAZ=qux"}
	out := e.Apply(input)
	if out[0] != "FOO=bar" || out[1] != "BAZ=qux" {
		t.Fatalf("expected passthrough, got %v", out)
	}
}

func TestEncode_EmptyMode_DefaultsToNone(t *testing.T) {
	e := NewEncoder("")
	if e.mode != EncodeModeNone {
		t.Fatalf("expected none, got %s", e.mode)
	}
}

func TestEncode_Base64_EncodesValues(t *testing.T) {
	e := NewEncoder(EncodeModeBase64)
	out := e.Apply([]string{"SECRET=hunter2", "TOKEN=abc123"})
	expect0 := "SECRET=" + b64("hunter2")
	expect1 := "TOKEN=" + b64("abc123")
	if out[0] != expect0 {
		t.Fatalf("got %s, want %s", out[0], expect0)
	}
	if out[1] != expect1 {
		t.Fatalf("got %s, want %s", out[1], expect1)
	}
}

func TestEncode_Base64_MalformedEntry_PassedThrough(t *testing.T) {
	e := NewEncoder(EncodeModeBase64)
	out := e.Apply([]string{"MALFORMED"})
	if out[0] != "MALFORMED" {
		t.Fatalf("expected passthrough, got %s", out[0])
	}
}

func TestEncode_Base64_EmptyValue(t *testing.T) {
	e := NewEncoder(EncodeModeBase64)
	out := e.Apply([]string{"EMPTY="})
	expected := "EMPTY=" + b64("")
	if out[0] != expected {
		t.Fatalf("got %s, want %s", out[0], expected)
	}
}

func TestEncode_DoesNotMutateInput(t *testing.T) {
	e := NewEncoder(EncodeModeBase64)
	input := []string{"KEY=val"}
	_ = e.Apply(input)
	if input[0] != "KEY=val" {
		t.Fatal("input was mutated")
	}
}
