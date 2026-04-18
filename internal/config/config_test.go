package config

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}
	return p
}

func TestLoad_EmptyPath(t *testing.T) {
	_, err := Load("")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := Load("/nonexistent/path/config.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoad_MissingVaultAddress(t *testing.T) {
	p := writeTemp(t, "vault_token: tok\ncommand: [echo]\n")
	_, err := Load(p)
	if err == nil {
		t.Fatal("expected validation error for missing vault_address")
	}
}

func TestLoad_MissingCommand(t *testing.T) {
	p := writeTemp(t, "vault_address: http://127.0.0.1:8200\nvault_token: tok\n")
	_, err := Load(p)
	if err == nil {
		t.Fatal("expected validation error for missing command")
	}
}

func TestLoad_Success(t *testing.T) {
	content := `
vault_address: http://127.0.0.1:8200
vault_token: mytoken
command: ["env"]
secrets:
  - path: secret/myapp
    env_vars:
      DB_PASSWORD: password
`
	p := writeTemp(t, content)
	cfg, err := Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.VaultAddress != "http://127.0.0.1:8200" {
		t.Errorf("unexpected vault_address: %s", cfg.VaultAddress)
	}
	if len(cfg.Secrets) != 1 {
		t.Fatalf("expected 1 secret mapping, got %d", len(cfg.Secrets))
	}
	if cfg.Secrets[0].EnvVars["DB_PASSWORD"] != "password" {
		t.Errorf("unexpected env var value")
	}
}
