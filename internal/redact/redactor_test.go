package redact_test

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/redact"
)

func TestNew_DefaultMask(t *testing.T) {
	r := redact.New("")
	if r.Mask() != "[REDACTED]" {
		t.Errorf("expected default mask, got %q", r.Mask())
	}
}

func TestNew_CustomMask(t *testing.T) {
	r := redact.New("***")
	if r.Mask() != "***" {
		t.Errorf("expected '***', got %q", r.Mask())
	}
}

func TestAdd_IncreasesLen(t *testing.T) {
	r := redact.New("")
	r.Add("secret1", "secret2")
	if r.Len() != 2 {
		t.Errorf("expected 2 secrets, got %d", r.Len())
	}
}

func TestAdd_IgnoresEmpty(t *testing.T) {
	r := redact.New("")
	r.Add("", "real-secret", "")
	if r.Len() != 1 {
		t.Errorf("expected 1 secret, got %d", r.Len())
	}
}

func TestRedact_ReplacesSecret(t *testing.T) {
	r := redact.New("[REDACTED]")
	r.Add("s3cr3t")
	out := r.Redact("the password is s3cr3t ok")
	if out != "the password is [REDACTED] ok" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestRedact_MultipleSecrets(t *testing.T) {
	r := redact.New("[REDACTED]")
	r.Add("alpha", "beta")
	out := r.Redact("alpha and beta are secrets")
	if out != "[REDACTED] and [REDACTED] are secrets" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestRedact_NoSecrets_Passthrough(t *testing.T) {
	r := redact.New("")
	input := "nothing to hide"
	if got := r.Redact(input); got != input {
		t.Errorf("expected passthrough, got %q", got)
	}
}
