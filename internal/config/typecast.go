package config

import "github.com/vaultpipe/vaultpipe/internal/env"

// TypecastConfig holds configuration for the typecast pipeline stage.
type TypecastConfig struct {
	// Mode controls how values are coerced.
	// Accepted values: "none", "bool", "numeric", "auto".
	Mode string `yaml:"mode"`
}

// ResolveTypecaster returns a Typecaster built from cfg, or a passthrough
// Typecaster when cfg is nil.
func ResolveTypecaster(cfg *TypecastConfig) *env.Typecaster {
	if cfg == nil {
		return env.NewTypecaster(env.CastNone)
	}
	return env.NewTypecaster(env.CastMode(cfg.Mode))
}
