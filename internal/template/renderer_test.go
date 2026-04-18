package template

import (
	"testing"
)

func testSecrets() map[string]map[string]string {
	return map[string]map[string]string{
		"secret/db": {
			"password": "s3cr3t",
			"user":     "admin",
		},
		"secret/api": {
			"key": "abc123",
		},
	}
}

func TestRender_PlainString(t *testing.T) {
	r := NewRenderer(testSecrets())
	out, err := r.Render("hello world")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "hello world" {
		t.Errorf("expected %q got %q", "hello world", out)
	}
}

func TestRender_SecretFunc(t *testing.T) {
	r := NewRenderer(testSecrets())
	out, err := r.Render(`{{ secret "secret/db" "password" }}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "s3cr3t" {
		t.Errorf("expected %q got %q", "s3cr3t", out)
	}
}

func TestRender_UnknownPath(t *testing.T) {
	r := NewRenderer(testSecrets())
	_, err := r.Render(`{{ secret "secret/missing" "key" }}`)
	if err == nil {
		t.Fatal("expected error for missing path")
	}
}

func TestRender_UnknownKey(t *testing.T) {
	r := NewRenderer(testSecrets())
	_, err := r.Render(`{{ secret "secret/db" "nokey" }}`)
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestRenderMap(t *testing.T) {
	r := NewRenderer(testSecrets())
	input := map[string]string{
		"DB_PASS": `{{ secret "secret/db" "password" }}`,
		"API_KEY": `{{ secret "secret/api" "key" }}`,
		"STATIC":  "unchanged",
	}
	out, err := r.RenderMap(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB_PASS"] != "s3cr3t" {
		t.Errorf("DB_PASS: expected %q got %q", "s3cr3t", out["DB_PASS"])
	}
	if out["API_KEY"] != "abc123" {
		t.Errorf("API_KEY: expected %q got %q", "abc123", out["API_KEY"])
	}
	if out["STATIC"] != "unchanged" {
		t.Errorf("STATIC: expected %q got %q", "unchanged", out["STATIC"])
	}
}

func TestRenderMap_Error(t *testing.T) {
	r := NewRenderer(testSecrets())
	input := map[string]string{
		"BAD": `{{ secret "nope" "nope" }}`,
	}
	_, err := r.RenderMap(input)
	if err == nil {
		t.Fatal("expected error")
	}
}
