package config

import "github.com/wakeward/vaultpipe/internal/env"

// LimitConfig holds env var count limit settings from YAML config.
type LimitConfig struct {
	Max    int    `yaml:"max"`
	Policy string `yaml:"policy"` // "error" | "truncate"
}

// ResolveLimit returns a configured Limiter from an optional LimitConfig.
// A nil cfg or max <= 0 produces a no-op limiter.
func ResolveLimit(cfg *LimitConfig) *env.Limiter {
	if cfg == nil || cfg.Max <= 0 {
		return env.NewLimiter(0, env.LimitPolicyError)
	}
	policy := env.LimitPolicyError
	if cfg.Policy == "truncate" {
		policy = env.LimitPolicyTruncate
	}
	return env.NewLimiter(cfg.Max, policy)
}
