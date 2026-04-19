package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSort_Ascending(t *testing.T) {
	s := NewSorter(false)
	input := []string{"ZEBRA=1", "APPLE=2", "MANGO=3"}
	out := s.Sort(input)
	assert.Equal(t, []string{"APPLE=2", "MANGO=3", "ZEBRA=1"}, out)
}

func TestSort_Descending(t *testing.T) {
	s := NewSorter(true)
	input := []string{"ZEBRA=1", "APPLE=2", "MANGO=3"}
	out := s.Sort(input)
	assert.Equal(t, []string{"ZEBRA=1", "MANGO=3", "APPLE=2"}, out)
}

func TestSort_DoesNotMutateInput(t *testing.T) {
	s := NewSorter(false)
	input := []string{"Z=1", "A=2"}
	orig := []string{"Z=1", "A=2"}
	s.Sort(input)
	assert.Equal(t, orig, input)
}

func TestSort_Empty(t *testing.T) {
	s := NewSorter(false)
	assert.Empty(t, s.Sort(nil))
	assert.Empty(t, s.Sort([]string{}))
}

func TestSort_MalformedEntry_UsedAsKey(t *testing.T) {
	s := NewSorter(false)
	input := []string{"ZZZ", "AAA", "MMM=val"}
	out := s.Sort(input)
	assert.Equal(t, []string{"AAA", "MMM=val", "ZZZ"}, out)
}

func TestSort_StableForEqualKeys(t *testing.T) {
	s := NewSorter(false)
	input := []string{"B=2", "A=1", "A=3"}
	out := s.Sort(input)
	assert.Equal(t, "A=1", out[0])
	assert.Equal(t, "A=3", out[1])
}
