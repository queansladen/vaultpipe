package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRename_EmptyMapping_Passthrough(t *testing.T) {
	r := NewRenamer(nil)
	input := []string{"FOO=bar", "BAZ=qux"}
	got := r.Apply(input)
	assert.Equal(t, input, got)
}

func TestRename_RenamesMatchingKey(t *testing.T) {
	r := NewRenamer(map[string]string{"OLD_KEY": "NEW_KEY"})
	got := r.Apply([]string{"OLD_KEY=hello"})
	assert.Equal(t, []string{"NEW_KEY=hello"}, got)
}

func TestRename_NonMatchingKey_Passthrough(t *testing.T) {
	r := NewRenamer(map[string]string{"OLD_KEY": "NEW_KEY"})
	got := r.Apply([]string{"OTHER=value"})
	assert.Equal(t, []string{"OTHER=value"}, got)
}

func TestRename_MalformedEntry_PassedThrough(t *testing.T) {
	r := NewRenamer(map[string]string{"FOO": "BAR"})
	got := r.Apply([]string{"MALFORMED"})
	assert.Equal(t, []string{"MALFORMED"}, got)
}

func TestRename_ValueContainsEquals(t *testing.T) {
	r := NewRenamer(map[string]string{"K": "K2"})
	got := r.Apply([]string{"K=a=b=c"})
	assert.Equal(t, []string{"K2=a=b=c"}, got)
}

func TestRename_IgnoresEmptyMappingEntries(t *testing.T) {
	r := NewRenamer(map[string]string{"": "NEW", "OLD": ""})
	assert.Empty(t, r.mapping)
}

func TestRename_MultipleKeys(t *testing.T) {
	r := NewRenamer(map[string]string{
		"A": "X",
		"B": "Y",
	})
	got := r.Apply([]string{"A=1", "B=2", "C=3"})
	assert.Equal(t, []string{"X=1", "Y=2", "C=3"}, got)
}
