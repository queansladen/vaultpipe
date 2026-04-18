package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAllow_ExactMatch(t *testing.T) {
	f := NewFilter([]string{"SECRET_TOKEN", "AWS_SECRET_ACCESS_KEY"}, nil)
	assert.False(t, f.Allow("SECRET_TOKEN"))
	assert.False(t, f.Allow("AWS_SECRET_ACCESS_KEY"))
	assert.True(t, f.Allow("HOME"))
}

func TestAllow_PrefixMatch(t *testing.T) {
	f := NewFilter(nil, []string{"VAULT_", "TF_VAR_"})
	assert.False(t, f.Allow("VAULT_TOKEN"))
	assert.False(t, f.Allow("VAULT_ADDR"))
	assert.False(t, f.Allow("TF_VAR_region"))
	assert.True(t, f.Allow("PATH"))
}

func TestAllow_NoRules_AllowsAll(t *testing.T) {
	f := NewFilter(nil, nil)
	assert.True(t, f.Allow("ANY_KEY"))
}

func TestApply_FiltersSlice(t *testing.T) {
	f := NewFilter([]string{"SECRET"}, []string{"VAULT_"})
	input := []string{"HOME=/root", "SECRET=abc", "VAULT_TOKEN=xyz", "PATH=/usr/bin"}
	got := f.Apply(input)
	assert.Equal(t, []string{"HOME=/root", "PATH=/usr/bin"}, got)
}

func TestApply_ValueContainsEquals(t *testing.T) {
	f := NewFilter([]string{"DROP"}, nil)
	input := []string{"DROP=a=b=c", "KEEP=x=y"}
	got := f.Apply(input)
	assert.Equal(t, []string{"KEEP=x=y"}, got)
}

func TestApply_EmptyInput(t *testing.T) {
	f := NewFilter([]string{"X"}, []string{"Y_"})
	assert.Empty(t, f.Apply(nil))
	assert.Empty(t, f.Apply([]string{}))
}
