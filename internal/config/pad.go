package config

import "github.com/nicholasgasior/vaultpipe/internal/env"

// PadConfig holds the configuration for value padding.
type PadConfig struct {
	Mode  string `yaml:"mode"`
	Width int    `yaml:"width"`
	Char  string `yaml:"char"`
}

// ResolvePadder constructs an env.Padder from cfg.
// Returns a no-op padder when cfg is nil or mode is empty/none.
func ResolvePadder(cfg *PadConfig) *env.Padder {
	if cfg == nil {
		return env.NewPadder(env.PadModeNone, 0, 0)
	}

	mode := env.PadMode(cfg.Mode)
	switch mode {
	case env.PadModeLeft, env.PadModeRight:
		// valid
	default:
		return env.NewPadder(env.PadModeNone, 0, 0)
	}

	var char rune
	if len(cfg.Char) > 0 {
		char = rune(cfg.Char[0])
	}

	return env.NewPadder(mode, cfg.Width, char)
}
