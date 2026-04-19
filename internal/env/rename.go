package env

import "strings"

// Renamer renames environment variable keys according to a mapping.
type Renamer struct {
	mapping map[string]string
}

// NewRenamer creates a Renamer from a map of old key -> new key.
func NewRenamer(mapping map[string]string) *Renamer {
	m := make(map[string]string, len(mapping))
	for k, v := range mapping {
		if k != "" && v != "" {
			m[k] = v
		}
	}
	return &Renamer{mapping: m}
}

// Apply renames keys in the provided env pairs slice.
// Pairs that do not match any mapping are passed through unchanged.
// Malformed pairs (no '=') are passed through as-is.
func (r *Renamer) Apply(pairs []string) []string {
	if len(r.mapping) == 0 {
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
		if newKey, ok := r.mapping[key]; ok {
			out = append(out, newKey+"="+val)
		} else {
			out = append(out, p)
		}
	}
	return out
}
