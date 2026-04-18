package config

import "github.com/yourusername/vaultpipe/internal/env"

// InheritConfig describes how the parent environment should be inherited by
// the child process launched by vaultpipe.
type InheritConfig struct {
	// Mode is one of: "all", "none", "filtered". Defaults to "all".
	Mode string `yaml:"mode"`
	// Allow is a list of exact variable names or KEY_ prefixes used when
	// Mode is "filtered".
	Allow []string `yaml:"allow"`
}

// ResolveInheritMode converts the string mode from config into an InheritMode
// constant understood by the env package. Unknown values fall back to InheritAll.
func ResolveInheritMode(cfg InheritConfig) (env.InheritMode, *env.Filter) {
	switch cfg.Mode {
	case "none":
		return env.InheritNone, nil
	case "filtered":
		f := env.NewFilter(cfg.Allow)
		return env.InheritFiltered, f
	default:
		return env.InheritAll, nil
	}
}
