package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFlatten_NoneMode_Passthrough(t *testing.T) {
	f := NewFlattener(FlattenNone, "/")
	input := []string{"SECRET/DB/PASSWORD=hunter2", "PLAIN=value"}
	out, err := f.Apply(input)
	require.NoError(t, err)
	assert.Equal(t, input, out)
}

func TestFlatten_EmptyMode_Passthrough(t *testing.T) {
	f := NewFlattener("", "/")
	input := []string{"A/B=val"}
	out, err := f.Apply(input)
	require.NoError(t, err)
	assert.Equal(t, input, out)
}

func TestFlatten_EmptySplit_Passthrough(t *testing.T) {
	f := NewFlattener(FlattenUnderscore, "")
	input := []string{"A/B=val"}
	out, err := f.Apply(input)
	require.NoError(t, err)
	assert.Equal(t, input, out)
}

func TestFlatten_Underscore_ReplacesSlash(t *testing.T) {
	f := NewFlattener(FlattenUnderscore, "/")
	out, err := f.Apply([]string{"SECRET/DB/PASS=s3cr3t"})
	require.NoError(t, err)
	assert.Equal(t, []string{"SECRET_DB_PASS=s3cr3t"}, out)
}

func TestFlatten_Dot_ReplacesSeparator(t *testing.T) {
	f := NewFlattener(FlattenDot, "__")
	out, err := f.Apply([]string{"APP__DB__HOST=localhost"})
	require.NoError(t, err)
	assert.Equal(t, []string{"APP.DB.HOST=localhost"}, out)
}

func TestFlatten_Dash_ReplacesSeparator(t *testing.T) {
	f := NewFlattener(FlattenDash, ".")
	out, err := f.Apply([]string{"app.db.host=localhost"})
	require.NoError(t, err)
	assert.Equal(t, []string{"app-db-host=localhost"}, out)
}

func TestFlatten_ValueContainsEquals_OnlyKeyTransformed(t *testing.T) {
	f := NewFlattener(FlattenUnderscore, "/")
	out, err := f.Apply([]string{"A/B=x=y=z"})
	require.NoError(t, err)
	assert.Equal(t, []string{"A_B=x=y=z"}, out)
}

func TestFlatten_MalformedEntry_PassedThrough(t *testing.T) {
	f := NewFlattener(FlattenUnderscore, "/")
	input := []string{"NOEQUALS"}
	out, err := f.Apply(input)
	require.NoError(t, err)
	assert.Equal(t, input, out)
}

func TestFlatten_DoesNotMutateInput(t *testing.T) {
	f := NewFlattener(FlattenUnderscore, "/")
	input := []string{"A/B=val", "C/D=other"}
	orig := make([]string, len(input))
	copy(orig, input)
	_, err := f.Apply(input)
	require.NoError(t, err)
	assert.Equal(t, orig, input)
}
