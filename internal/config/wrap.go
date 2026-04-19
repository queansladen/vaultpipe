package config

import (
	"strings"

	"github.com/yourusername/vaultpipe/internal/env"
)

// WrapConfig holds optional key and value transformation modes.
type WrapConfig struct {
	KeyTransform   string `yaml:"key_transform"`   // "upper", "lower", or ""
	ValueTransform string `yaml:"value_transform"` // "upper", "lower", or ""
}

// ResolveWrapper builds an env.Wrapper from the provided WrapConfig.
// A nil config returns a no-op wrapper.
func ResolveWrapper(cfg *WrapConfig) *env.Wrapper {
	if cfg == nil {
		return env.NewWrapper(nil, nil)
	}
	return env.NewWrapper(
		resolveTransformFn(cfg.KeyTransform),
		resolveTransformFn(cfg.ValueTransform),
	)
}

func resolveTransformFn(mode string) func(string) string {
	switch strings.ToLower(mode) {
	case "upper":
		return strings.ToUpper
	case "lower":
		return strings.ToLower
	default:
		return nil
	}
}
