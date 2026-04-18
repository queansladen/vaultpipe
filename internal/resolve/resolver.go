package resolve

import (
	"fmt"

	"github.com/vaultpipe/vaultpipe/internal/audit"
	"github.com/vaultpipe/vaultpipe/internal/vault"
)

// SecretMapping maps environment variable names to vault paths and keys.
type SecretMapping struct {
	EnvVar string
	Path   string
	Key    string
}

// Resolver resolves secret mappings into environment variable key/value pairs.
type Resolver struct {
	client  *vault.Client
	auditor audit.Auditor
}

// NewResolver creates a new Resolver.
func NewResolver(client *vault.Client, auditor audit.Auditor) *Resolver {
	return &Resolver{client: client, auditor: auditor}
}

// Resolve takes a slice of SecretMappings and returns a map of env var name to secret value.
func (r *Resolver) Resolve(mappings []SecretMapping) (map[string]string, error) {
	result := make(map[string]string, len(mappings))
	for _, m := range mappings {
		data, err := r.client.ReadSecretData(m.Path)
		if err != nil {
			r.auditor.SecretFetchFailed(m.Path, err)
			return nil, fmt.Errorf("resolve %s: %w", m.EnvVar, err)
		}
		val, ok := data[m.Key]
		if !ok {
			err := fmt.Errorf("key %q not found at path %q", m.Key, m.Path)
			r.auditor.SecretFetchFailed(m.Path, err)
			return nil, err
		}
		r.auditor.SecretFetched(m.Path)
		result[m.EnvVar] = fmt.Sprintf("%v", val)
	}
	return result, nil
}
