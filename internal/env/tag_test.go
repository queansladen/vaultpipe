package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTag_EmptyPrefixSuffix_Passthrough(t *testing.T) {
	tagger := NewTagger("", "")
	pairs := []string{"FOO=bar", "BAZ=qux"}
	got := tagger.Tag(pairs)
	assert.Equal(t, pairs, got)
}

func TestTag_PrefixOnly(t *testing.T) {
	tagger := NewTagger("[vp]", "")
	got := tagger.Tag([]string{"SECRET=abc"})
	assert.Equal(t, []string{"SECRET=[vp]abc"}, got)
}

func TestTag_SuffixOnly(t *testing.T) {
	tagger := NewTagger("", "[end]")
	got := tagger.Tag([]string{"TOKEN=xyz"})
	assert.Equal(t, []string{"TOKEN=xyz[end]"}, got)
}

func TestTag_PrefixAndSuffix(t *testing.T) {
	tagger := NewTagger("<<", ">>")
	got := tagger.Tag([]string{"KEY=value"})
	assert.Equal(t, []string{"KEY=<<value>>"}, got)
}

func TestTag_MalformedEntry_PassedThrough(t *testing.T) {
	tagger := NewTagger("[", "]")
	got := tagger.Tag([]string{"NOEQUALS"})
	assert.Equal(t, []string{"NOEQUALS"}, got)
}

func TestTag_ValueContainsEquals(t *testing.T) {
	tagger := NewTagger("(", ")")
	got := tagger.Tag([]string{"URL=http://x.com?a=1"})
	assert.Equal(t, []string{"URL=(http://x.com?a=1)"}, got)
}

func TestStrip_EmptyPrefixSuffix_Passthrough(t *testing.T) {
	tagger := NewTagger("", "")
	pairs := []string{"FOO=bar"}
	got := tagger.Strip(pairs)
	assert.Equal(t, pairs, got)
}

func TestStrip_RemovesPrefixAndSuffix(t *testing.T) {
	tagger := NewTagger("<<", ">>")
	got := tagger.Strip([]string{"KEY=<<value>>"})
	assert.Equal(t, []string{"KEY=value"}, got)
}

func TestStrip_NoMarkers_ValueUnchanged(t *testing.T) {
	tagger := NewTagger("<<", ">>")
	got := tagger.Strip([]string{"KEY=plain"})
	assert.Equal(t, []string{"KEY=plain"}, got)
}

func TestStrip_MalformedEntry_PassedThrough(t *testing.T) {
	tagger := NewTagger("[", "]")
	got := tagger.Strip([]string{"NOEQUALS"})
	assert.Equal(t, []string{"NOEQUALS"}, got)
}
