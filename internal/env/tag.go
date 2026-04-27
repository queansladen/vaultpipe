package env

import (
	"fmt"
	"strings"
)

// Tagger annotates environment variable values with a prefix or suffix tag,
// useful for marking secrets injected by vaultpipe at runtime.
type Tagger struct {
	prefix string
	suffix string
}

// NewTagger returns a Tagger that wraps values with the given prefix and suffix.
// Either may be empty, in which case that side is left unchanged.
func NewTagger(prefix, suffix string) *Tagger {
	return &Tagger{prefix: prefix, suffix: suffix}
}

// Tag annotates a slice of KEY=VALUE pairs, returning new pairs with the value
// wrapped by the configured prefix and suffix.
func (t *Tagger) Tag(pairs []string) []string {
	if t.prefix == "" && t.suffix == "" {
		return pairs
	}
	out := make([]string, 0, len(pairs))
	for _, p := range pairs {
		idx := strings.IndexByte(p, '=')
		if idx < 0 {
			out = append(out, p)
			continue
		}
		key := p[:idx]
		val := p[idx+1:]
		out = append(out, fmt.Sprintf("%s=%s%s%s", key, t.prefix, val, t.suffix))
	}
	return out
}

// Strip removes the configured prefix and suffix tags from values in a slice
// of KEY=VALUE pairs. Values that do not carry both markers are left unchanged.
func (t *Tagger) Strip(pairs []string) []string {
	if t.prefix == "" && t.suffix == "" {
		return pairs
	}
	out := make([]string, 0, len(pairs))
	for _, p := range pairs {
		idx := strings.IndexByte(p, '=')
		if idx < 0 {
			out = append(out, p)
			continue
		}
		key := p[:idx]
		val := p[idx+1:]
		if t.prefix != "" && strings.HasPrefix(val, t.prefix) {
			val = val[len(t.prefix):]
		}
		if t.suffix != "" && strings.HasSuffix(val, t.suffix) {
			val = val[:len(val)-len(t.suffix)]
		}
		out = append(out, fmt.Sprintf("%s=%s", key, val))
	}
	return out
}
