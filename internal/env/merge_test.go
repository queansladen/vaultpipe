package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMerge_Empty(t *testing.T) {
	m := NewMerger(StrategyLast)
	result := m.Merge()
	assert.Empty(t, result)
}

func TestMerge_SingleSlice(t *testing.T) {
	m := NewMerger(StrategyLast)
	result := m.Merge([]string{"FOO=bar", "BAZ=qux"})
	assert.Equal(t, []string{"FOO=bar", "BAZ=qux"}, result)
}

func TestMerge_StrategyLast_LaterWins(t *testing.T) {
	m := NewMerger(StrategyLast)
	base := []string{"FOO=original", "KEEP=yes"}
	overlay := []string{"FOO=override"}
	result := m.Merge(base, overlay)
	assert.Contains(t, result, "FOO=override")
	assert.Contains(t, result, "KEEP=yes")
}

func TestMerge_StrategyFirst_EarlierWins(t *testing.T) {
	m := NewMerger(StrategyFirst)
	base := []string{"FOO=original"}
	overlay := []string{"FOO=override", "NEW=value"}
	result := m.Merge(base, overlay)
	assert.Contains(t, result, "FOO=original")
	assert.Contains(t, result, "NEW=value")
}

func TestMerge_SkipsMalformed(t *testing.T) {
	m := NewMerger(StrategyLast)
	result := m.Merge([]string{"VALID=yes", "NOEQUALSSIGN", "ALSO=fine"})
	assert.Equal(t, []string{"VALID=yes", "ALSO=fine"}, result)
}

func TestMerge_PreservesInsertionOrder(t *testing.T) {
	m := NewMerger(StrategyLast)
	a := []string{"A=1", "B=2"}
	b := []string{"C=3", "A=99"}
	result := m.Merge(a, b)
	assert.Equal(t, []string{"A=99", "B=2", "C=3"}, result)
}

func TestMerge_ThreeSlices(t *testing.T) {
	m := NewMerger(StrategyLast)
	result := m.Merge(
		[]string{"X=1"},
		[]string{"X=2"},
		[]string{"X=3"},
	)
	assert.Equal(t, []string{"X=3"}, result)
}
