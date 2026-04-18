package audit

import (
	"bytes"
	"testing"
)

func TestNoopLogger_NoOutput(t *testing.T) {
	n := &NoopLogger{}
	if err := n.SecretFetched("path", "key"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := n.SecretFetchFailed("path", "key", "reason"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := n.ProcessStarted("cmd"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewAuditor_Disabled(t *testing.T) {
	buf := &bytes.Buffer{}
	a := NewAuditor(false, buf)
	if _, ok := a.(*NoopLogger); !ok {
		t.Error("expected NoopLogger when disabled")
	}
	_ = a.SecretFetched("s", "k")
	if buf.Len() != 0 {
		t.Error("expected no output from NoopLogger")
	}
}

func TestNewAuditor_Enabled(t *testing.T) {
	buf := &bytes.Buffer{}
	a := NewAuditor(true, buf)
	if _, ok := a.(*Logger); !ok {
		t.Error("expected Logger when enabled")
	}
	_ = a.SecretFetched("secret/app", "db_pass")
	if buf.Len() == 0 {
		t.Error("expected output from Logger")
	}
}
