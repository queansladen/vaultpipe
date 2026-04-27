package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompact_NoneMode_Passthrough(t *testing.T) {
	c := NewCompactor(CompactNone)
	input := []string{"A=1", "B=", "A=2"}
	got := c.Apply(input)
	assert.Equal(t, input, got)
}

func TestCompact_EmptyMode_DefaultsToNone(t *testing.T) {
	c := NewCompactor("")
	input := []string{"A=1", "B="}
	got := c.Apply(input)
	assert.Equal(t, input, got)
}

func TestCompact_DuplicatesMode_KeepsLast(t *testing.T) {
	c := NewCompactor(CompactDuplicates)
	input := []string{"A=first", "B=keep", "A=last"}
	got := c.Apply(input)
	assert.Equal(t, []string{"B=keep", "A=last"}, got)
}

func TestCompact_EmptyMode_RemovesEmptyValues(t *testing.T) {
	c := NewCompactor(CompactEmpty)
	input := []string{"A=hello", "B=", "C=world"}
	got := c.Apply(input)
	assert.Equal(t, []string{"A=hello", "C=world"}, got)
}

func TestCompact_BothMode_RemovesDupesAndEmpty(t *testing.T) {
	c := NewCompactor(CompactBoth)
	input := []string{"A=first", "B=", "A=last", "C=keep"}
	got := c.Apply(input)
	assert.Equal(t, []string{"A=last", "C=keep"}, got)
}

func TestCompact_EmptyInput_ReturnsNil(t *testing.T) {
	c := NewCompactor(CompactBoth)
	got := c.Apply([]string{})
	assert.Nil(t, got)
}

func TestCompact_AllRemoved_ReturnsNil(t *testing.T) {
	c := NewCompactor(CompactEmpty)
	input := []string{"A=", "B="}
	got := c.Apply(input)
	assert.Nil(t, got)
}

func TestCompact_MalformedEntry_PassedThrough(t *testing.T) {
	c := NewCompactor(CompactBoth)
	input := []string{"NOEQUALS", "A=ok"}
	got := c.Apply(input)
	assert.Equal(t, []string{"NOEQUALS", "A=ok"}, got)
}

func TestCompact_ValueContainsEquals(t *testing.T) {
	c := NewCompactor(CompactDuplicates)
	input := []string{"A=x=y", "B=z"}
	got := c.Apply(input)
	assert.Equal(t, []string{"A=x=y", "B=z"}, got)
}
