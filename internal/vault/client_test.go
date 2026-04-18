package vault

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient_MissingAddress(t *testing.T) {
	t.Setenv("VAULT_ADDR", "")
	t.Setenv("VAULT_TOKEN", "")

	_, err := NewClient(Config{})
	if err == nil {
		t.Fatal("expected error when address is missing, got nil")
	}
}

func TestNewClient_MissingToken(t *testing.T) {
	t.Setenv("VAULT_TOKEN", "")

	_, err := NewClient(Config{Address: "http://127.0.0.1:8200"})
	if err == nil {
		t.Fatal("expected error when token is missing, got nil")
	}
}

func TestNewClient_Success(t *testing.T) {
	client, err := NewClient(Config{
		Address: "http://127.0.0.1:8200",
		Token:   "test-token",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestReadSecretData_KVv1(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":{"API_KEY":"abc123","DB_PASS":"secret"}}`))
	}))
	defer server.Close()

	client, err := NewClient(Config{Address: server.URL, Token: "test"})
	if err != nil {
		t.Fatalf("client creation failed: %v", err)
	}

	data, err := client.ReadSecretData("secret/myapp")
	if err != nil {
		t.Fatalf("unexpected error reading secret: %v", err)
	}

	if data["API_KEY"] != "abc123" {
		t.Errorf("expected API_KEY=abc123, got %v", data["API_KEY"])
	}
}

func TestReadSecretData_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`null`))
	}))
	defer server.Close()

	client, err := NewClient(Config{Address: server.URL, Token: "test"})
	if err != nil {
		t.Fatalf("client creation failed: %v", err)
	}

	_, err = client.ReadSecretData("secret/missing")
	if err == nil {
		t.Fatal("expected error for missing secret, got nil")
	}
}
