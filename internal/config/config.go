package config

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

// SecretMapping maps a Vault path to a set of env var overrides.
type SecretMapping struct {
	Path    string            `yaml:"path"`
	EnvVars map[string]string `yaml:"env_vars"`
}

// Config holds the top-level vaultpipe configuration.
type Config struct {
	VaultAddress string          `yaml:"vault_address"`
	VaultToken   string          `yaml:"vault_token"`
	Secrets      []SecretMapping `yaml:"secrets"`
	Command      []string        `yaml:"command"`
}

// Load reads and parses a YAML config file from the given path.
func Load(path string) (*Config, error) {
	if path == "" {
		return nil, errors.New("config path must not be empty")
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	decoder.KnownFields(true)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Validate checks that required fields are present.
func (c *Config) Validate() error {
	if c.VaultAddress == "" {
		return errors.New("vault_address is required")
	}
	if c.VaultToken == "" {
		return errors.New("vault_token is required")
	}
	if len(c.Command) == 0 {
		return errors.New("command is required")
	}
	return nil
}
