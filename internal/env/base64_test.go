package env

import (
	"encoding/base64"
	"fmt"
	"testing"
)

func b64enc(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func TestDecode_NoneMode_Passthrough(t *testing.T) {
	d := NewDecoder(DecodeModeNone)
	input := []string{"KEY=value"}
	got, err := d.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got[0] != "KEY=value" {
		t.Fatalf("expected passthrough, got %q", got[0])
	}
}

func TestDecode_EmptyMode_DefaultsToNone(t *testing.T) {
	d := NewDecoder("")
	if d.mode != DecodeModeNone {
		t.Fatalf("expected none, got %q", d.mode)
	}
}

func TestDecode_ValuesMode_DecodesValue(t *testing.T) {
	d := NewDecoder(DecodeModeValues)
	encoded := b64enc("supersecret")
	input := []string{fmt.Sprintf("TOKEN=%s", encoded)}
	got, err := d.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got[0] != "TOKEN=supersecret" {
		t.Fatalf("expected decoded value, got %q", got[0])
	}
}

func TestDecode_KeysMode_DecodesKey(t *testing.T) {
	d := NewDecoder(DecodeModeKeys)
	encodedKey := b64enc("MY_KEY")
	input := []string{fmt.Sprintf("%s=rawvalue", encodedKey)}
	got, err := d.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got[0] != "MY_KEY=rawvalue" {
		t.Fatalf("expected decoded key, got %q", got[0])
	}
}

func TestDecode_BothMode_DecodesKeyAndValue(t *testing.T) {
	d := NewDecoder(DecodeModeBoth)
	input := []string{fmt.Sprintf("%s=%s", b64enc("API_KEY"), b64enc("abc123"))}
	got, err := d.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got[0] != "API_KEY=abc123" {
		t.Fatalf("expected both decoded, got %q", got[0])
	}
}

func TestDecode_MalformedEntry_PassedThrough(t *testing.T) {
	d := NewDecoder(DecodeModeValues)
	input := []string{"NOKEYEQUALS"}
	got, err := d.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got[0] != "NOKEYEQUALS" {
		t.Fatalf("expected passthrough, got %q", got[0])
	}
}

func TestDecode_InvalidBase64Value_PassedThrough(t *testing.T) {
	d := NewDecoder(DecodeModeValues)
	input := []string{"KEY=not-valid-base64!!!"}
	got, err := d.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got[0] != "KEY=not-valid-base64!!!" {
		t.Fatalf("expected passthrough on invalid base64, got %q", got[0])
	}
}
