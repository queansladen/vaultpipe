package config

import "github.com/yourusername/vaultpipe/internal/env"

// SampleConfig holds sampling configuration from the YAML config.
type SampleConfig struct {
	Mode string `yaml:"mode"`
	N    int    `yaml:"n"`
}

// ResolveSampler constructs an env.Sampler from optional config.
// When cfg is nil or mode is empty, a passthrough sampler is returned.
func ResolveSampler(cfg *SampleConfig) *env.Sampler {
	if cfg == nil || cfg.Mode == "" {
		return env.NewSampler(env.SampleModeNone, 0, nil)
	}

	mode := env.SampleMode(cfg.Mode)
	switch mode {
	case env.SampleModeFirst, env.SampleModeLast, env.SampleModeRandom:
		// valid
	default:
		mode = env.SampleModeNone
	}

	return env.NewSampler(mode, cfg.N, nil)
}
