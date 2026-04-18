package env

import "strings"

// Filter holds rules for excluding environment variables from a snapshot
// or builder before passing them to a child process.
type Filter struct {
	prefixes []string
	exact    map[string]struct{}
}

// NewFilter returns a Filter that will drop variables matching any of the
// given exact names or prefixes.
func NewFilter(exact []string, prefixes []string) *Filter {
	em := make(map[string]struct{}, len(exact))
	for _, k := range exact {
		em[k] = struct{}{}
	}
	return &Filter{
		prefixes: prefixes,
		exact:    em,
	}
}

// Allow returns true when the key should be kept.
func (f *Filter) Allow(key string) bool {
	if _, blocked := f.exact[key]; blocked {
		return false
	}
	for _, p := range f.prefixes {
		if strings.HasPrefix(key, p) {
			return false
		}
	}
	return true
}

// Apply returns a new slice containing only the pairs whose key is allowed.
// Each pair is expected to be in KEY=VALUE form.
func (f *Filter) Apply(pairs []string) []string {
	out := make([]string, 0, len(pairs))
	for _, pair := range pairs {
		key := pair
		if idx := strings.IndexByte(pair, '='); idx >= 0 {
			key = pair[:idx]
		}
		if f.Allow(key) {
			out = append(out, pair)
		}
	}
	return out
}
