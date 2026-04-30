package config

import "github.com/yourusername/vaultpipe/internal/env"

// InterpConfig controls variable interpolation inside env-pair values.
type InterpConfig struct {
	// Enabled turns interpolation on. Defaults to false.
	Enabled bool `yaml:"enabled"`
	// Fallback, when true, allows unresolved references to fall through to the
	// host OS environment. Defaults to false.
	Fallback bool `yaml:"fallback"`
	// Overlay provides additional key→value pairs used during expansion before
	// the pair list and the OS environment are consulted.
	Overlay map[string]string `yaml:"overlay"`
}

// ResolveInterpolator returns a configured *env.Interpolator when interpolation
// is enabled, or nil when it is disabled or cfg is nil.
func ResolveInterpolator(cfg *InterpConfig) *env.Interpolator {
	if cfg == nil || !cfg.Enabled {
		return nil
	}
	return env.NewInterpolator(cfg.Overlay, cfg.Fallback)
}
