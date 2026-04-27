package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCaser_NoneMode_Passthrough(t *testing.T) {
	fn := NewCaser("key", CaseModeNone)
	in := []string{"Hello=World", "foo=bar"}
	out, err := fn(in)
	require.NoError(t, err)
	assert.Equal(t, in, out)
}

func TestCaser_EmptyMode_Passthrough(t *testing.T) {
	fn := NewCaser("key", "")
	in := []string{"Hello=World"}
	out, err := fn(in)
	require.NoError(t, err)
	assert.Equal(t, in, out)
}

func TestCaser_UpperKey(t *testing.T) {
	fn := NewCaser("key", CaseModeUpper)
	out, err := fn([]string{"hello=World", "foo_bar=baz"})
	require.NoError(t, err)
	assert.Equal(t, []string{"HELLO=World", "FOO_BAR=baz"}, out)
}

func TestCaser_LowerValue(t *testing.T) {
	fn := NewCaser("value", CaseModeLower)
	out, err := fn([]string{"KEY=HELLO_WORLD"})
	require.NoError(t, err)
	assert.Equal(t, []string{"KEY=hello_world"}, out)
}

func TestCaser_UpperBoth(t *testing.T) {
	fn := NewCaser("both", CaseModeUpper)
	out, err := fn([]string{"key=value"})
	require.NoError(t, err)
	assert.Equal(t, []string{"KEY=VALUE"}, out)
}

func TestCaser_MalformedEntry_PassedThrough(t *testing.T) {
	fn := NewCaser("key", CaseModeUpper)
	out, err := fn([]string{"NOEQUALS"})
	require.NoError(t, err)
	assert.Equal(t, []string{"NOEQUALS"}, out)
}

func TestCaser_EmptySlice(t *testing.T) {
	fn := NewCaser("key", CaseModeUpper)
	out, err := fn([]string{})
	require.NoError(t, err)
	assert.Empty(t, out)
}
