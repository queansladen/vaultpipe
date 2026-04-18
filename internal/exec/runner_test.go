package exec

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestNewRunner(t *testing.T) {
	secrets := map[string]string{"DB_PASS": "secret"}
	r := NewRunner(secrets)
	if r == nil {
		t.Fatal("expected non-nil runner")
	}
	if r.Env["DB_PASS"] != "secret" {
		t.Errorf("expected DB_PASS=secret, got %s", r.Env["DB_PASS"])
	}
}

func TestRun_EmptyCommand(t *testing.T) {
	r := NewRunner(nil)
	err := r.Run("", nil)
	if err == nil {
		t.Fatal("expected error for empty command")
	}
	if !strings.Contains(err.Error(), "command must not be empty") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRun_InjectsSecrets(t *testing.T) {
	secrets := map[string]string{"VAULTPIPE_TEST_VAR": "injected_value"}
	r := NewRunner(secrets)

	// Redirect stdout to a buffer via a temp file.
	tmpFile, err := os.CreateTemp("", "vaultpipe-test-*")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	r.Stdout = tmpFile

	if err := r.Run("sh", []string{"-c", "echo $VAULTPIPE_TEST_VAR"}); err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	tmpFile.Seek(0, 0)
	buf := new(bytes.Buffer)
	buf.ReadFrom(tmpFile)

	output := strings.TrimSpace(buf.String())
	if output != "injected_value" {
		t.Errorf("expected 'injected_value', got %q", output)
	}
}

func TestRun_OverridesExistingEnv(t *testing.T) {
	os.Setenv("VAULTPIPE_OVERRIDE", "original")
	defer os.Unsetenv("VAULTPIPE_OVERRIDE")

	secrets := map[string]string{"VAULTPIPE_OVERRIDE": "overridden"}
	r := NewRunner(secrets)

	tmpFile, err := os.CreateTemp("", "vaultpipe-override-*")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	r.Stdout = tmpFile

	if err := r.Run("sh", []string{"-c", "echo $VAULTPIPE_OVERRIDE"}); err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	tmpFile.Seek(0, 0)
	buf := new(bytes.Buffer)
	buf.ReadFrom(tmpFile)

	output := strings.TrimSpace(buf.String())
	if output != "overridden" {
		t.Errorf("expected 'overridden', got %q", output)
	}
}
