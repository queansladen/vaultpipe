package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolveSortConfig_NilConfig(t *testing.T) {
	cfg := ResolveSortConfig(nil)
	assert.False(t, cfg.Enabled)
	assert.Equal(t, "asc", cfg.Order)
}

func TestResolveSortConfig_ExplicitAsc(t *testing.T) {
	cfg := ResolveSortConfig(&SortConfig{Enabled: true, Order: "asc"})
	assert.True(t, cfg.Enabled)
	assert.Equal(t, "asc", cfg.Order)
}

func TestResolveSortConfig_ExplicitDesc(t *testing.T) {
	cfg := ResolveSortConfig(&SortConfig{Enabled: true, Order: "desc"})
	assert.True(t, cfg.Enabled)
	assert.Equal(t, "desc", cfg.Order)
}

func TestResolveSortConfig_InvalidOrder_DefaultsToAsc(t *testing.T) {
	cfg := ResolveSortConfig(&SortConfig{Enabled: true, Order: "random"})
	assert.Equal(t, "asc", cfg.Order)
}

func TestResolveSortConfig_Disabled(t *testing.T) {
	cfg := ResolveSortConfig(&SortConfig{Enabled: false, Order: "desc"})
	assert.False(t, cfg.Enabled)
	assert.Equal(t, "desc", cfg.Order)
}
