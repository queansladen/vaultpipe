package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReplace_NoRules_Passthrough(t *testing.T) {
	r := NewReplacer(nil)
	in := []string{"FOO=bar", "BAZ=qux"}
	assert.Equal(t, in, r.Apply(in))
}

func TestReplace_EmptyFind_Skipped(t *testing.T) {
	r := NewReplacer([]ReplaceRule{{Find: "", Replace: "x"}})
	in := []string{"FOO=bar"}
	assert.Equal(t, in, r.Apply(in))
}

func TestReplace_ReplacesMatchingSubstring(t *testing.T) {
	r := NewReplacer([]ReplaceRule{{Find: "localhost", Replace: "prod.internal"}})
	in := []string{"DB_HOST=localhost:5432"}
	out := r.Apply(in)
	assert.Equal(t, []string{"DB_HOST=prod.internal:5432"}, out)
}

func TestReplace_MultipleRules_AllApplied(t *testing.T) {
	r := NewReplacer([]ReplaceRule{
		{Find: "http", Replace: "https"},
		{Find: "localhost", Replace: "api.example.com"},
	})
	in := []string{"URL=http://localhost/path"}
	out := r.Apply(in)
	assert.Equal(t, []string{"URL=https://api.example.com/path"}, out)
}

func TestReplace_MalformedEntry_PassedThrough(t *testing.T) {
	r := NewReplacer([]ReplaceRule{{Find: "a", Replace: "b"}})
	in := []string{"MALFORMED"}
	assert.Equal(t, in, r.Apply(in))
}

func TestReplace_EmptySlice(t *testing.T) {
	r := NewReplacer([]ReplaceRule{{Find: "a", Replace: "b"}})
	assert.Empty(t, r.Apply([]string{}))
}

func TestReplace_DoesNotMutateInput(t *testing.T) {
	r := NewReplacer([]ReplaceRule{{Find: "old", Replace: "new"}})
	in := []string{"KEY=old_value"}
	orig := in[0]
	r.Apply(in)
	assert.Equal(t, orig, in[0])
}
