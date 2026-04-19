package env

import "strings"

// Prefixer adds or strips a prefix from environment variable keys.
type Prefixer struct {
	prefix string
}

// NewPrefixer returns a Prefixer that operates on the given prefix.
func NewPrefixer(prefix string) *Prefixer {
	return &Prefixer{prefix: prefix}
}

// Add returns a new slice with the prefix prepended to every key.
// Entries that already carry the prefix are left unchanged.
func (p *Prefixer) Add(pairs []string) []string {
	if p.prefix == "" {
		return pairs
	}
	out := make([]string, 0, len(pairs))
	for _, pair := range pairs {
		key, val, ok := strings.Cut(pair, "=")
		if !ok {
			out = append(out, pair)
			continue
		}
		if !strings.HasPrefix(key, p.prefix) {
			key = p.prefix + key
		}
		out = append(out, key+"="+val)
	}
	return out
}

// Strip returns a new slice with the prefix removed from every key.
// Entries that do not carry the prefix are dropped.
func (p *Prefixer) Strip(pairs []string) []string {
	if p.prefix == "" {
		return pairs
	}
	out := make([]string, 0, len(pairs))
	for _, pair := range pairs {
		key, val, ok := strings.Cut(pair, "=")
		if !ok {
			continue
		}
		after, found := strings.CutPrefix(key, p.prefix)
		if !found {
			continue
		}
		out = append(out, after+"="+val)
	}
	return out
}
