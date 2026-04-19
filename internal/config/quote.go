package config

// QuoteConfig controls shell-quoting of secret values injected into the
// process environment when exporting to a shell-eval-able format.
type QuoteConfig struct {
	Enabled bool `yaml:"enabled"`
}

// ResolveQuote returns the effective QuoteConfig. A nil pointer is treated
// as disabled (the default behaviour — values are injected verbatim).
func ResolveQuote(q *QuoteConfig) QuoteConfig {
	if q == nil {
		return QuoteConfig{Enabled: false}
	}
	return *q
}
