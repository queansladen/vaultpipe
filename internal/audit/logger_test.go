package audit

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func newTestLogger() (*Logger, *bytes.Buffer) {
	buf := &bytes.Buffer{}
	return NewLogger(buf), buf
}

func TestNewLogger_DefaultsToStderr(t *testing.T) {
	l := NewLogger(nil)
	if l.writer == nil {
		t.Fatal("expected non-nil writer")
	}
}

func TestLog_WritesValidJSON(t *testing.T) {
	l, buf := newTestLogger()
	if err := l.Log("test_action", "secret/data/foo", "password", true, "ok"); err != nil {
		t.Fatalf("Log() error: %v", err)
	}
	var e Event
	if err := json.Unmarshal(buf.Bytes(), &e); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if e.Action != "test_action" {
		t.Errorf("expected action test_action, got %s", e.Action)
	}
	if e.Path != "secret/data/foo" {
		t.Errorf("unexpected path: %s", e.Path)
	}
	if !e.Success {
		t.Error("expected success=true")
	}
}

func TestSecretFetched(t *testing.T) {
	l, buf := newTestLogger()
	_ = l.SecretFetched("secret/myapp", "api_key")
	if !strings.Contains(buf.String(), "secret_fetch") {
		t.Error("expected secret_fetch action in output")
	}
}

func TestSecretFetchFailed(t *testing.T) {
	l, buf := newTestLogger()
	_ = l.SecretFetchFailed("secret/myapp", "token", "permission denied")
	var e Event
	_ = json.Unmarshal(buf.Bytes(), &e)
	if e.Success {
		t.Error("expected success=false")
	}
	if e.Message != "permission denied" {
		t.Errorf("unexpected message: %s", e.Message)
	}
}

func TestProcessStarted(t *testing.T) {
	l, buf := newTestLogger()
	_ = l.ProcessStarted("myapp --serve")
	var e Event
	_ = json.Unmarshal(buf.Bytes(), &e)
	if e.Action != "process_start" {
		t.Errorf("expected process_start, got %s", e.Action)
	}
	if e.Message != "myapp --serve" {
		t.Errorf("unexpected message: %s", e.Message)
	}
}
