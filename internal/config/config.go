package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// SecretRef describes a single Vault secret path to fetch.
type SecretRef struct {
	Path string `yaml:"path"`
}

// Config holds the full vaultpipe configuration.
type Config struct {
	VaultAddress string            `yaml:"vault_address"`
	VaultToken   string            `yaml:"vault_token"`
	Secrets      []SecretRef       `yaml:"secrets"`
	Env          map[string]string `yaml:"env"`
	Command      []string          `yaml:"command"`
}

// Load reads and validates a YAML config file from the given path.
// If path is empty, it returns an error.
func Load(path string) (*Config, error) {
	if path == "" {
		return nil, errors.New("config path must not be empty")
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open config: %w", err)
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	decoder.KnownFields(true)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("decode config: %w", err)
	}

	if err := validate(&cfg); err != nil {
		return nil, err
	}

	// Allow vault_token to be sourced from environment.
	if cfg.VaultToken == "" {
		cfg.VaultToken = os.Getenv("VAULT_TOKEN")
	}

	return &cfg, nil
}

func validate(cfg *Config) error {
	if cfg.VaultAddress == "" {
		return errors.New("vault_address is required")
	}
	if len(cfg.Command) == 0 {
		return errors.New("command is required")
	}
	return nil
}
