package env

import "strings"

// Truncator limits the byte length of environment variable values.
type Truncator struct {
	maxBytes int
	suffix   string
}

// NewTruncator returns a Truncator that clips values to maxBytes.
// If suffix is non-empty it is appended after truncation (and counts toward the limit).
func NewTruncator(maxBytes int, suffix string) *Truncator {
	if maxBytes <= 0 {
		maxBytes = 0
	}
	return &Truncator{maxBytes: maxBytes, suffix: suffix}
}

// Truncate clips a single value string to the configured byte limit.
func (t *Truncator) Truncate(value string) string {
	if t.maxBytes == 0 || len(value) <= t.maxBytes {
		return value
	}
	cut := t.maxBytes - len(t.suffix)
	if cut < 0 {
		cut = 0
	}
	return value[:cut] + t.suffix
}

// Apply truncates the values in a KEY=VALUE slice, leaving keys intact.
func (t *Truncator) Apply(pairs []string) []string {
	out := make([]string, len(pairs))
	for i, p := range pairs {
		idx := strings.IndexByte(p, '=')
		if idx < 0 {
			out[i] = p
			continue
		}
		key := p[:idx]
		val := t.Truncate(p[idx+1:])
		out[i] = key + "=" + val
	}
	return out
}
