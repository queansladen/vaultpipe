package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDedupe_NoDuplicates(t *testing.T) {
	input := []string{"FOO=1", "BAR=2", "BAZ=3"}
	out := Dedupe(input)
	assert.Equal(t, input, out)
}

func TestDedupe_KeepsLast(t *testing.T) {
	input := []string{"FOO=first", "BAR=1", "FOO=last"}
	out := Dedupe(input)
	// FOO should appear at its original position but with the last value.
	assert.Equal(t, []string{"FOO=last", "BAR=1"}, out)
}

func TestDedupe_Empty(t *testing.T) {
	out := Dedupe(nil)
	assert.Empty(t, out)
}

func TestDedupe_SingleEntry(t *testing.T) {
	out := Dedupe([]string{"KEY=val"})
	assert.Equal(t, []string{"KEY=val"}, out)
}

func TestDedupe_MultipleConsecutiveDuplicates(t *testing.T) {
	input := []string{"A=1", "A=2", "A=3", "B=x"}
	out := Dedupe(input)
	assert.Equal(t, []string{"A=3", "B=x"}, out)
}

func TestDedupeMap_ReturnsLastValue(t *testing.T) {
	input := []string{"FOO=first", "BAR=hello", "FOO=second"}
	m := DedupeMap(input)
	assert.Equal(t, "second", m["FOO"])
	assert.Equal(t, "hello", m["BAR"])
}

func TestDedupeMap_Empty(t *testing.T) {
	m := DedupeMap(nil)
	assert.NotNil(t, m)
	assert.Len(t, m, 0)
}

func TestDedupeMap_ValueContainsEquals(t *testing.T) {
	input := []string{"TOKEN=abc=def=ghi"}
	m := DedupeMap(input)
	assert.Equal(t, "abc=def=ghi", m["TOKEN"])
}
