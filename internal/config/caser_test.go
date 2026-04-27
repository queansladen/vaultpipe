package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolveCaser_NilConfig_Passthrough(t *testing.T) {
	fn := ResolveCaser(nil)
	in := []string{"hello=world"}
	out, err := fn(in)
	require.NoError(t, err)
	assert.Equal(t, in, out)
}

func TestResolveCaser_EmptyTarget_DefaultsToKey(t *testing.T) {
	fn := ResolveCaser(&CaserConfig{Mode: "upper"})
	out, err := fn([]string{"hello=world"})
	require.NoError(t, err)
	assert.Equal(t, []string{"HELLO=world"}, out)
}

func TestResolveCaser_ExplicitUpperKey(t *testing.T) {
	fn := ResolveCaser(&CaserConfig{Target: "key", Mode: "upper"})
	out, err := fn([]string{"foo=bar"})
	require.NoError(t, err)
	assert.Equal(t, []string{"FOO=bar"}, out)
}

func TestResolveCaser_ExplicitLowerValue(t *testing.T) {
	fn := ResolveCaser(&CaserConfig{Target: "value", Mode: "lower"})
	out, err := fn([]string{"KEY=UPPER"})
	require.NoError(t, err)
	assert.Equal(t, []string{"KEY=upper"}, out)
}

func TestResolveCaser_EmptyMode_Passthrough(t *testing.T) {
	fn := ResolveCaser(&CaserConfig{Target: "key", Mode: ""})
	in := []string{"hello=world"}
	out, err := fn(in)
	require.NoError(t, err)
	assert.Equal(t, in, out)
}
