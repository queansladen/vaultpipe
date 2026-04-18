package vault

import (
	"fmt"
	"os"

	vaultapi "github.com/hashicorp/vault/api"
)

// Client wraps the Vault API client with helper methods for secret retrieval.
type Client struct {
	api *vaultapi.Client
}

// Config holds configuration for connecting to Vault.
type Config struct {
	Address string
	Token   string
}

// NewClient creates a new Vault client from the given config.
// If Address or Token are empty, it falls back to VAULT_ADDR and VAULT_TOKEN env vars.
func NewClient(cfg Config) (*Client, error) {
	vcfg := vaultapi.DefaultConfig()

	addr := cfg.Address
	if addr == "" {
		addr = os.Getenv("VAULT_ADDR")
	}
	if addr == "" {
		return nil, fmt.Errorf("vault address not set: provide --vault-addr or VAULT_ADDR")
	}
	vcfg.Address = addr

	client, err := vaultapi.NewClient(vcfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create vault client: %w", err)
	}

	token := cfg.Token
	if token == "" {
		token = os.Getenv("VAULT_TOKEN")
	}
	if token == "" {
		return nil, fmt.Errorf("vault token not set: provide --vault-token or VAULT_TOKEN")
	}
	client.SetToken(token)

	return &Client{api: client}, nil
}

// ReadSecretData reads a KV secret at the given path and returns its data map.
// Supports both KV v1 (secret/foo) and KV v2 (secret/data/foo) paths.
func (c *Client) ReadSecretData(path string) (map[string]interface{}, error) {
	secret, err := c.api.Logical().Read(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read secret at %q: %w", path, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("no secret found at path %q", path)
	}

	// KV v2 wraps data under a "data" key.
	if nested, ok := secret.Data["data"]; ok {
		if m, ok := nested.(map[string]interface{}); ok {
			return m, nil
		}
	}

	return secret.Data, nil
}
