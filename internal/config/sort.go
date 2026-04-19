package config

// SortConfig controls how injected environment variables are ordered.
type SortConfig struct {
	Enabled    bool   `yaml:"enabled"`
	Order      string `yaml:"order"` // "asc" or "desc"
}

// ResolveSortConfig returns a normalised SortConfig.
// Defaults: enabled=false, order="asc".
func ResolveSortConfig(c *SortConfig) SortConfig {
	if c == nil {
		return SortConfig{Enabled: false, Order: "asc"}
	}
	order := c.Order
	if order != "asc" && order != "desc" {
		order = "asc"
	}
	return SortConfig{
		Enabled: c.Enabled,
		Order:   order,
	}
}
