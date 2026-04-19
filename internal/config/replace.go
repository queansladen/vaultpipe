package config

import "github.com/yourusername/vaultpipe/internal/env"

// ReplaceConfig holds find-and-replace rules from the YAML config.
type ReplaceConfig struct {
	Rules []ReplaceRuleConfig `yaml:"rules"`
}

// ReplaceRuleConfig is a single rule entry in the YAML config.
type ReplaceRuleConfig struct {
	Find    string `yaml:"find"`
	Replace string `yaml:"replace"`
}

// ResolveReplacer builds an env.Replacer from config, returning nil when
// no rules are defined.
func ResolveReplacer(cfg *ReplaceConfig) *env.Replacer {
	if cfg == nil || len(cfg.Rules) == 0 {
		return env.NewReplacer(nil)
	}
	rules := make([]env.ReplaceRule, 0, len(cfg.Rules))
	for _, r := range cfg.Rules {
		if r.Find == "" {
			continue
		}
		rules = append(rules, env.ReplaceRule{
			Find:    r.Find,
			Replace: r.Replace,
		})
	}
	return env.NewReplacer(rules)
}
