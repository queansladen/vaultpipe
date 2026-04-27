package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolveTagger_NilConfig(t *testing.T) {
	tagger := ResolveTagger(nil)
	assert.NotNil(t, tagger)
	// A nil config should produce a no-op tagger — values pass through unchanged.
	pairs := []string{"KEY=value"}
	assert.Equal(t, pairs, tagger.Tag(pairs))
}

func TestResolveTagger_EmptyStrings(t *testing.T) {
	tagger := ResolveTagger(&TagConfig{Prefix: "", Suffix: ""})
	pairs := []string{"KEY=value"}
	assert.Equal(t, pairs, tagger.Tag(pairs))
}

func TestResolveTagger_WithPrefix(t *testing.T) {
	tagger := ResolveTagger(&TagConfig{Prefix: "[vp]"})
	got := tagger.Tag([]string{"SECRET=abc"})
	assert.Equal(t, []string{"SECRET=[vp]abc"}, got)
}

func TestResolveTagger_WithSuffix(t *testing.T) {
	tagger := ResolveTagger(&TagConfig{Suffix: "[end]"})
	got := tagger.Tag([]string{"TOKEN=xyz"})
	assert.Equal(t, []string{"TOKEN=xyz[end]"}, got)
}

func TestResolveTagger_WithPrefixAndSuffix(t *testing.T) {
	tagger := ResolveTagger(&TagConfig{Prefix: "<<", Suffix: ">>"})
	got := tagger.Tag([]string{"KEY=value"})
	assert.Equal(t, []string{"KEY=<<value>>"}, got)
}
