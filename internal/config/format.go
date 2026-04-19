package config

import "github.com/yourusername/vaultpipe/internal/env"

// FormatConfig holds output format settings from the YAML config.
type FormatConfig struct {
	Mode string `yaml:"mode"` // raw | export | dotenv
}

// ResolveFormat returns an env.FormatMode from the config block.
// A nil config or empty mode defaults to FormatModeRaw.
func ResolveFormat(cfg *FormatConfig) env.FormatMode {
	if cfg == nil || cfg.Mode == "" {
		return env.FormatModeRaw
	}
	switch env.FormatMode(cfg.Mode) {
	case env.FormatModeExport, env.FormatModeDotenv:
		return env.FormatMode(cfg.Mode)
	default:
		return env.FormatModeRaw
	}
}
