package config

import "github.com/yourusername/vaultpipe/internal/env"

// CaserConfig holds configuration for the key/value casing transformation.
type CaserConfig struct {
	Target string `yaml:"target"` // "key", "value", or "both"
	Mode   string `yaml:"mode"`   // "upper", "lower", or "none"
}

// ResolveCaser returns a pipeline step based on the provided CaserConfig.
// If cfg is nil the step is a no-op (none mode).
func ResolveCaser(cfg *CaserConfig) func([]string) ([]string, error) {
	if cfg == nil {
		return env.NewCaser("key", env.CaseModeNone)
	}
	target := cfg.Target
	if target == "" {
		target = "key"
	}
	mode := env.CaseMode(cfg.Mode)
	if mode == "" {
		mode = env.CaseModeNone
	}
	return env.NewCaser(env.CaseMode(target), mode)
}
