package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLimit_NoLimit_Passthrough(t *testing.T) {
	l := NewLimiter(0, LimitPolicyError)
	pairs := []string{"A=1", "B=2", "C=3"}
	out, err := l.Apply(pairs)
	require.NoError(t, err)
	assert.Equal(t, pairs, out)
}

func TestLimit_WithinLimit_Passthrough(t *testing.T) {
	l := NewLimiter(5, LimitPolicyError)
	pairs := []string{"A=1", "B=2"}
	out, err := l.Apply(pairs)
	require.NoError(t, err)
	assert.Equal(t, pairs, out)
}

func TestLimit_ExceedsLimit_PolicyError(t *testing.T) {
	l := NewLimiter(2, LimitPolicyError)
	_, err := l.Apply([]string{"A=1", "B=2", "C=3"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "exceeds limit")
}

func TestLimit_ExceedsLimit_PolicyTruncate(t *testing.T) {
	l := NewLimiter(2, LimitPolicyTruncate)
	out, err := l.Apply([]string{"A=1", "B=2", "C=3"})
	require.NoError(t, err)
	assert.Equal(t, []string{"A=1", "B=2"}, out)
}

func TestLimit_ExactLimit_Passthrough(t *testing.T) {
	l := NewLimiter(3, LimitPolicyError)
	pairs := []string{"A=1", "B=2", "C=3"}
	out, err := l.Apply(pairs)
	require.NoError(t, err)
	assert.Equal(t, pairs, out)
}

func TestLimit_EmptySlice(t *testing.T) {
	l := NewLimiter(2, LimitPolicyError)
	out, err := l.Apply([]string{})
	require.NoError(t, err)
	assert.Empty(t, out)
}
