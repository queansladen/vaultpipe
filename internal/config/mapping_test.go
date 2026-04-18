package config_test

import (
	"testing"

	"github.com/vaultpipe/vaultpipe/internal/config"
)

func TestToMappings_Empty(t *testing.T) {
	_, err := config.ToMappings(nil)
	if err == nil {
		t.Fatal("expected error for empty entries")
	}
}

func TestToMappings_MissingEnv(t *testing.T) {
	entries := []config.SecretEntry{{Path: "secret/db", Key: "pass"}}
	_, err := config.ToMappings(entries)
	if err == nil {
		t.Fatal("expected error for missing env")
	}
}

func TestToMappings_MissingPath(t *testing.T) {
	entries := []config.SecretEntry{{EnvVar: "DB_PASS", Key: "pass"}}
	_, err := config.ToMappings(entries)
	if err == nil {
		t.Fatal("expected error for missing path")
	}
}

func TestToMappings_MissingKey(t *testing.T) {
	entries := []config.SecretEntry{{EnvVar: "DB_PASS", Path: "secret/db"}}
	_, err := config.ToMappings(entries)
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestToMappings_Success(t *testing.T) {
	entries := []config.SecretEntry{
		{EnvVar: "DB_PASS", Path: "secret/db", Key: "password"},
		{EnvVar: "API_KEY", Path: "secret/api", Key: "key"},
	}
	mappings, err := config.ToMappings(entries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(mappings) != 2 {
		t.Fatalf("expected 2 mappings, got %d", len(mappings))
	}
	if mappings[0].EnvVar != "DB_PASS" {
		t.Errorf("expected DB_PASS, got %s", mappings[0].EnvVar)
	}
	if mappings[1].Key != "key" {
		t.Errorf("expected key, got %s", mappings[1].Key)
	}
}
