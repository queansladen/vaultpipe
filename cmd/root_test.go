package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("write temp config: %v", err)
	}
	return p
}

func TestExecute_MissingConfigFlag(t *testing.T) {
	rootCmd.SetArgs([]string{})
	var buf bytes.Buffer
	rootCmd.SetErr(&buf)
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error when --config flag is missing")
	}
}

func TestExecute_MissingConfigFile(t *testing.T) {
	rootCmd.SetArgs([]string{"--config", "/nonexistent/path.yaml"})
	var buf bytes.Buffer
	rootCmd.SetErr(&buf)
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error for missing config file")
	}
}

func TestExecute_InvalidConfig(t *testing.T) {
	p := writeTempConfig(t, "vault:\n  address: \"\"\n")
	rootCmd.SetArgs([]string{"--config", p})
	var buf bytes.Buffer
	rootCmd.SetErr(&buf)
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error for invalid config")
	}
}
