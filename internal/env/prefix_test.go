package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdd_EmptyPrefix_Noop(t *testing.T) {
	p := NewPrefixer("")
	input := []string{"FOO=bar", "BAZ=qux"}
	assert.Equal(t, input, p.Add(input))
}

func TestAdd_PrependsPrefix(t *testing.T) {
	p := NewPrefixer("APP_")
	got := p.Add([]string{"FOO=bar", "BAZ=qux"})
	assert.Equal(t, []string{"APP_FOO=bar", "APP_BAZ=qux"}, got)
}

func TestAdd_SkipsAlreadyPrefixed(t *testing.T) {
	p := NewPrefixer("APP_")
	got := p.Add([]string{"APP_FOO=bar"})
	assert.Equal(t, []string{"APP_FOO=bar"}, got)
}

func TestAdd_MalformedEntry_PassedThrough(t *testing.T) {
	p := NewPrefixer("APP_")
	got := p.Add([]string{"NOEQUALS"})
	assert.Equal(t, []string{"NOEQUALS"}, got)
}

func TestStrip_EmptyPrefix_Noop(t *testing.T) {
	p := NewPrefixer("")
	input := []string{"FOO=bar"}
	assert.Equal(t, input, p.Strip(input))
}

func TestStrip_RemovesPrefix(t *testing.T) {
	p := NewPrefixer("APP_")
	got := p.Strip([]string{"APP_FOO=bar", "APP_BAZ=qux"})
	assert.Equal(t, []string{"FOO=bar", "BAZ=qux"}, got)
}

func TestStrip_DropsUnprefixedEntries(t *testing.T) {
	p := NewPrefixer("APP_")
	got := p.Strip([]string{"APP_FOO=bar", "OTHER=val"})
	assert.Equal(t, []string{"FOO=bar"}, got)
}

func TestStrip_MalformedEntry_Dropped(t *testing.T) {
	p := NewPrefixer("APP_")
	got := p.Strip([]string{"NOEQUALS"})
	assert.Empty(t, got)
}

func TestAdd_ValueContainsEquals(t *testing.T) {
	p := NewPrefixer("X_")
	got := p.Add([]string{"KEY=a=b"})
	assert.Equal(t, []string{"X_KEY=a=b"}, got)
}
