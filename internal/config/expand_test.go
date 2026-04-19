package config

import (
	"testing"
)

func lookup(key string) string {
	m := map[string]string{
		"ENV":    "production",
		"REGION": "us-east-1",
	}
	return m[key]
}

func TestExpandSecretRefs_NoVariables(t *testing.T) {
	in := []SecretMapping{{Env: "DB_PASS", Path: "secret/db", Key: "password"}}
	out := ExpandSecretRefs(in, lookup)
	if out[0].Path != "secret/db" {
		t.Errorf("unexpected path: %q", out[0].Path)
	}
}

func TestExpandSecretRefs_ExpandsPath(t *testing.T) {
	in := []SecretMapping{{Env: "DB_PASS", Path: "secret/${ENV}/db", Key: "password"}}
	out := ExpandSecretRefs(in, lookup)
	if out[0].Path != "secret/production/db" {
		t.Errorf("unexpected path: %q", out[0].Path)
	}
}

func TestExpandSecretRefs_ExpandsKey(t *testing.T) {
	in := []SecretMapping{{Env: "X", Path: "secret/app", Key: "key_${REGION}"}}
	out := ExpandSecretRefs(in, lookup)
	if out[0].Key != "key_us-east-1" {
		t.Errorf("unexpected key: %q", out[0].Key)
	}
}

func TestExpandSecretRefs_NilLookup(t *testing.T) {
	in := []SecretMapping{{Env: "X", Path: "secret/${ENV}/db", Key: "pass"}}
	out := ExpandSecretRefs(in, nil)
	if out[0].Path != "secret/${ENV}/db" {
		t.Errorf("expected unchanged path, got %q", out[0].Path)
	}
}

func TestExpandSecretRefs_UnknownVar(t *testing.T) {
	in := []SecretMapping{{Env: "X", Path: "secret/${UNKNOWN}/db", Key: "pass"}}
	out := ExpandSecretRefs(in, lookup)
	if out[0].Path != "secret//db" {
		t.Errorf("expected empty substitution, got %q", out[0].Path)
	}
}

func TestExpandSecretRefs_PreservesEnvField(t *testing.T) {
	in := []SecretMapping{{Env: "MY_VAR", Path: "p", Key: "k"}}
	out := ExpandSecretRefs(in, lookup)
	if out[0].Env != "MY_VAR" {
		t.Errorf("env field mutated: %q", out[0].Env)
	}
}
