package config

import "github.com/yourusername/vaultpipe/internal/env"

// MaskConfig holds masking options from the YAML config.
type MaskConfig struct {
	Mode      string `yaml:"mode"`
	Char      string `yaml:"char"`
	KeepChars int    `yaml:"keep_chars"`
}

// ResolveMask returns a Masker built from optional config.
// Defaults to MaskNone when cfg is nil.
func ResolveMask(cfg *MaskConfig) *env.Masker {
	if cfg == nil {
		return env.NewMasker(env.MaskNone, "*", 0)
	}
	mode := env.MaskMode(cfg.Mode)
	switch mode {
	case env.MaskFull, env.MaskPartial, env.MaskNone:
		// valid
	default:
		mode = env.MaskNone
	}
	return env.NewMasker(mode, cfg.Char, cfg.KeepChars)
}
