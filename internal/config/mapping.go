package config

import (
	"fmt"

	"github.com/vaultpipe/vaultpipe/internal/resolve"
)

// SecretEntry represents a single secret mapping in the config file.
type SecretEntry struct {
	EnvVar string `yaml:"env"`
	Path   string `yaml:"path"`
	Key    string `yaml:"key"`
}

// ToMappings converts config SecretEntries to resolve.SecretMappings.
func ToMappings(entries []SecretEntry) ([]resolve.SecretMapping, error) {
	if len(entries) == 0 {
		return nil, fmt.Errorf("no secret mappings defined")
	}
	mappings := make([]resolve.SecretMapping, 0, len(entries))
	for i, e := range entries {
		if e.EnvVar == "" {
			return nil, fmt.Errorf("entry %d: missing env", i)
		}
		if e.Path == "" {
			return nil, fmt.Errorf("entry %d: missing path", i)
		}
		if e.Key == "" {
			return nil, fmt.Errorf("entry %d: missing key", i)
		}
		mappings = append(mappings, resolve.SecretMapping{
			EnvVar: e.EnvVar,
			Path:   e.Path,
			Key:    e.Key,
		})
	}
	return mappings, nil
}
