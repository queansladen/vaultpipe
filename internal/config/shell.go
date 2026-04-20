package config

import "internal/env"

// ShellConfig holds shell-output formatting options from the YAML config.
type ShellConfig struct {
	Enabled bool   `yaml:"enabled"`
	Format  string `yaml:"format"`
	Quote   bool   `yaml:"quote"`
}

// ResolveSheller returns a configured Sheller from optional ShellConfig.
// If cfg is nil or disabled, ResolveSheller returns nil.
func ResolveSheller(cfg *ShellConfig) *env.Sheller {
	if cfg == nil || !cfg.Enabled {
		return nil
	}
	format := env.ShellFormat(cfg.Format)
	switch format {
	case env.ShellFormatExport, env.ShellFormatInline, env.ShellFormatUnset:
		// valid
	default:
		format = env.ShellFormatExport
	}
	return env.NewSheller(format, cfg.Quote)
}
