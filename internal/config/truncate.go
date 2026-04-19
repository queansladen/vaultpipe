package config

// TruncateConfig holds optional value-length limiting settings.
type TruncateConfig struct {
	// MaxBytes is the maximum byte length of any injected secret value.
	// Zero means unlimited.
	MaxBytes int    `yaml:"max_bytes"`
	Suffix   string `yaml:"suffix"`
}

// ResolveTruncate returns a canonical TruncateConfig from the loaded Config.
// If the section is absent, a zero-value (unlimited) config is returned.
func ResolveTruncate(cfg *Config) TruncateConfig {
	if cfg == nil {
		return TruncateConfig{}
	}
	return cfg.Truncate
}
