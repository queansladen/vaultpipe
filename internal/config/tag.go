package config

import "github.com/yourusername/vaultpipe/internal/env"

// TagConfig holds optional tagging settings for injected environment values.
type TagConfig struct {
	Prefix string `yaml:"prefix"`
	Suffix string `yaml:"suffix"`
}

// ResolveTagger returns an env.Tagger constructed from the provided TagConfig.
// If cfg is nil, a no-op tagger with empty prefix and suffix is returned.
func ResolveTagger(cfg *TagConfig) *env.Tagger {
	if cfg == nil {
		return env.NewTagger("", "")
	}
	return env.NewTagger(cfg.Prefix, cfg.Suffix)
}
