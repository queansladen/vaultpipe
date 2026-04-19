package config

import (
	"github.com/yourusername/vaultpipe/internal/env"
)

// BuildTransformer constructs an env.Transformer from the resolved config,
// wiring together coerce, prefix, dedupe, sort, truncate, and quote steps
// in a consistent order.
func BuildTransformer(cfg *Config) *env.Transformer {
	var steps []func([]string) []string

	// Coerce case (upper/lower)
	coerceCfg := cfg.Env.Coerce
	if coerceCfg != "" {
		c := env.NewCoercer(coerceCfg)
		steps = append(steps, c.Apply)
	}

	// Prefix
	prefix := ResolvePrefix(cfg)
	if prefix != "" {
		p := env.NewPrefixer(prefix)
		steps = append(steps, p.Add)
	}

	// Dedupe — always applied to remove duplicates before further processing
	steps = append(steps, env.Dedupe)

	// Sort
	sortCfg := ResolveSortConfig(cfg)
	if sortCfg.Enabled {
		s := env.NewSorter(sortCfg.Order)
		steps = append(steps, s.Sort)
	}

	// Truncate
	truncCfg := ResolveTruncate(cfg)
	if truncCfg.Limit > 0 {
		tr := env.NewTruncator(truncCfg.Limit, truncCfg.Suffix)
		steps = append(steps, tr.Truncate)
	}

	// Quote
	if ResolveQuote(cfg) {
		q := env.NewQuoter()
		steps = append(steps, q.Quote)
	}

	return env.NewTransformer(steps...)
}
