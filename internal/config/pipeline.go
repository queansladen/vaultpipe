package config

import (
	"github.com/yourusername/vaultpipe/internal/env"
)

// BuildPipeline constructs an ordered env.Pipeline from the resolved
// configuration. Stages are appended in a deterministic order so that
// transformations compose predictably regardless of config file ordering.
//
// Order:
//  1. Coerce (key/value casing)
//  2. Rename
//  3. Replace (substring substitution)
//  4. Prefix  (add / strip prefix)
//  5. Tag     (wrap with prefix/suffix strings)
//  6. Scrub   (remove empty or prefixed keys)
//  7. Mask    (redact secret values)
//  8. Encode  (base64 etc.)
//  9. Truncate
// 10. Sort
func BuildPipeline(cfg *Config) (*env.Pipeline, error) {
	p := env.NewPipeline()

	if coercer := ResolveCaser(cfg.Env.Caser); coercer != nil {
		p.Add(stageOf(coercer))
	}

	if renamer := NewRenamerStage(cfg); renamer != nil {
		p.Add(renamer)
	}

	if replacer := ResolveReplacer(cfg.Env.Replace); replacer != nil {
		p.Add(stageOf(replacer))
	}

	if prefixer := ResolvePrefix(cfg.Env.Prefix); prefixer != nil {
		p.Add(stageOf(prefixer))
	}

	if tagger := ResolveTagger(cfg.Env.Tag); tagger != nil {
		p.Add(stageOf(tagger))
	}

	if scrubber := NewScrubberStage(cfg); scrubber != nil {
		p.Add(scrubber)
	}

	if masker := ResolveMask(cfg.Env.Mask); masker != nil {
		p.Add(stageOf(masker))
	}

	if sorter := ResolveSortConfig(cfg.Env.Sort); sorter != nil {
		p.Add(stageOf(sorter))
	}

	return p, nil
}

// stageOf wraps any type that exposes an Apply([]string)([]string,error)
// method as an env.Stage.
func stageOf(a interface {
	Apply([]string) ([]string, error)
}) env.Stage {
	return a.Apply
}

// NewRenamerStage returns a stage that applies rename mappings, or nil when
// no mappings are configured.
func NewRenamerStage(cfg *Config) env.Stage {
	if len(cfg.Env.Rename) == 0 {
		return nil
	}
	m := make(map[string]string, len(cfg.Env.Rename))
	for _, r := range cfg.Env.Rename {
		m[r.From] = r.To
	}
	r := env.NewRenamer(m)
	return r.Apply
}

// NewScrubberStage returns a stage that scrubs env vars, or nil when scrub
// is not configured.
func NewScrubberStage(cfg *Config) env.Stage {
	s := env.NewScrubber(cfg.Env.Scrub.Mode, cfg.Env.Scrub.Prefixes)
	if s == nil {
		return nil
	}
	return s.Apply
}
