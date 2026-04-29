package config

import (
	"github.com/your-org/vaultpipe/internal/env"
)

// WindowConfig holds configuration for the env windowing feature.
type WindowConfig struct {
	Mode   string `yaml:"mode"`
	Offset int    `yaml:"offset"`
	Size   int    `yaml:"size"`
}

// ResolveWindower constructs an env.Windower from the provided config.
// If cfg is nil or the mode is empty, a passthrough Windower (none mode) is returned.
func ResolveWindower(cfg *WindowConfig) (*env.Windower, error) {
	if cfg == nil || cfg.Mode == "" {
		w, err := env.NewWindower(env.WindowModeNone, 0, 0)
		return w, err
	}
	w, err := env.NewWindower(env.WindowMode(cfg.Mode), cfg.Offset, cfg.Size)
	if err != nil {
		return nil, err
	}
	return w, nil
}
