package config

// PrefixConfig holds optional settings for adding or stripping an
// environment-variable prefix when secrets are injected into the child
// process environment.
type PrefixConfig struct {
	// Add is prepended to every secret-derived env key before the child
	// process is started.  Leave empty to disable.
	Add string `yaml:"add"`

	// Strip removes a leading prefix from secret-derived env keys before
	// injection.  Useful when Vault KV paths encode a service name that
	// should not appear in the process environment.
	Strip string `yaml:"strip"`
}

// ResolvePrefix returns the effective add/strip strings, falling back to
// empty strings when the config section is absent.
func ResolvePrefix(cfg *PrefixConfig) (add, strip string) {
	if cfg == nil {
		return "", ""
	}
	return cfg.Add, cfg.Strip
}
