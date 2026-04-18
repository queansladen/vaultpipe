package sanitize_test

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/sanitize"
)

func TestKey_Empty(t *testing.T) {
	_, err := sanitize.Key("")
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestKey_StartsWithDigit(t *testing.T) {
	_, err := sanitize.Key("1BAD")
	if err == nil {
		t.Fatal("expected error for key starting with digit")
	}
}

func TestKey_InvalidCharacter(t *testing.T) {
	_, err := sanitize.Key("bad.key")
	if err == nil {
		t.Fatal("expected error for key with dot")
	}
}

func TestKey_NormalisesHyphen(t *testing.T) {
	got, err := sanitize.Key("my-secret-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "MY_SECRET_KEY" {
		t.Fatalf("expected MY_SECRET_KEY, got %s", got)
	}
}

func TestKey_NormalisesSpace(t *testing.T) {
	got, err := sanitize.Key("my key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "MY_KEY" {
		t.Fatalf("expected MY_KEY, got %s", got)
	}
}

func TestKey_AlreadyValid(t *testing.T) {
	got, err := sanitize.Key("DB_PASSWORD")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "DB_PASSWORD" {
		t.Fatalf("expected DB_PASSWORD, got %s", got)
	}
}

func TestValue_TrimsWhitespace(t *testing.T) {
	got := sanitize.Value("  secret  ")
	if got != "secret" {
		t.Fatalf("expected 'secret', got %q", got)
	}
}

func TestValue_PreservesInterior(t *testing.T) {
	input := "  hello world  "
	got := sanitize.Value(input)
	if got != "hello world" {
		t.Fatalf("expected 'hello world', got %q", got)
	}
}

func TestValue_NoChange(t *testing.T) {
	got := sanitize.Value("already-clean")
	if got != "already-clean" {
		t.Fatalf("unexpected change: %q", got)
	}
}
