package config

import "github.com/yourusername/vaultpipe/internal/env"

// HashConfig holds configuration for the env value hashing stage.
type HashConfig struct {
	Mode string `yaml:"mode"`
}

// ResolveHasher returns a configured *env.Hasher derived from cfg.
// If cfg is nil or mode is empty, the hasher defaults to HashModeNone
// (i.e. values are passed through unchanged).
func ResolveHasher(cfg *HashConfig) *env.Hasher {
	if cfg == nil {
		return env.NewHasher(env.HashModeNone)
	}
	return env.NewHasher(env.HashMode(cfg.Mode))
}
